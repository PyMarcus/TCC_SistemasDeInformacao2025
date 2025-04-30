package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/constants"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/config"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/db"
	adapters "github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http/dto"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/repository"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	globalCfg *config.Config
)


// execute is the main entry point of the application.
// It initializes the logger, loads configuration, establishes a database connection,
// sets up the repositories and use cases, and starts the task processing pool.
func execute(){

	startTime := time.Now()

	loggerUsecase := usecase.NewLoggerUsecase(usecase.LoggerConfig)

	loggerUsecase.Error("[+] Starting...")

	cfg, err := config.LoadConfig(".env")
	if err != nil{
		loggerUsecase.Error("[-] Fail to load .env", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		os.Exit(1)
	}
	globalCfg = cfg

	// database conn
	connDB, err := db.NewPostgresConn(cfg.DatabaseUrl)
	if err != nil{
		loggerUsecase.Error("[-] Fail to connect with database ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		os.Exit(1)
	}
	
	// repositories
	datasetRepository := repository.NewDatasetRepository(connDB)
	questionRepository := repository.NewQuestionRepository(connDB)
	
	// usecases
	datasetUsecase := usecase.NewDatasetUsecase(datasetRepository)
	questionUsecase := usecase.NewQuestionUsecase(questionRepository)
	

	poolExecutor(datasetUsecase, questionUsecase, loggerUsecase, connDB)

	elapsedTime := time.Since(startTime).Seconds()
	loggerUsecase.Error("[+] Complete!", zap.String(constants.ELAPSED_TIME, fmt.Sprintf("%f", elapsedTime)))
}

// poolExecutor initializes the worker pool and task queue,
// retrieves all datasets and questions, and dispatches tasks
// for concurrent processing using a worker pool.
func poolExecutor(
		datasetUsecase  *usecase.DatasetUsecase,
		questionUsecase *usecase.QuestionUsecase,
		loggerUsecase 	*usecase.LoggerUsecase,
		connDB 			*gorm.DB){

	var wg sync.WaitGroup


	questions, err := questionUsecase.FindAll()
	if err != nil{
		loggerUsecase.Error("[-] Fail to list questions ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		os.Exit(1)
	}

	dataset, err := datasetUsecase.FindAll()
	if err != nil{
		loggerUsecase.Error("[-] Fail to list dataset ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		os.Exit(1)
	}

	// channel creation
	tasksCh := make(chan *domain.Task, len(dataset) * len(questions))

	// workerpool
	for i := 0; i < constants.WORKERS; i++{
		go workerPool(tasksCh, &wg, loggerUsecase, connDB, datasetUsecase, questions)
	}

	// input to channel
	for _, ds := range dataset {
		stringDatasetID := fmt.Sprintf("%d", ds.ID)
		loggerUsecase.Info("[+] Reading dataset ", zap.String(constants.DATASET_ID, stringDatasetID))
		for _, q := range questions {
			wg.Add(1)
			tasksCh <- &domain.Task{Dataset: ds, Question: q}
		}
	}

	wg.Wait()

	close(tasksCh)
	 
}

// workerPool listens to the task channel and processes each task
// concurrently by invoking the insertExecutor function.
// It is meant to be run as a goroutine.
func workerPool(tasksCh <-chan *domain.Task, wg *sync.WaitGroup, loggerUsecase *usecase.LoggerUsecase, connDB *gorm.DB, datasetUsecase *usecase.DatasetUsecase, questions []*domain.Question){
	for task := range tasksCh {
		// atention + code + question.
		task.Question.Question = constants.QUESTION_HEADER + task.Dataset.Class + task.Question.Question
		insertExecutor(wg, task.Dataset, loggerUsecase, connDB, datasetUsecase, task.Question)
	}
}

// insertExecutor handles the processing of a single dataset-question pair.
// It sends concurrent requests to two agents, handles responses or errors,
// logs errors, and persists valid results as Atom entities in the database.
func insertExecutor(wg *sync.WaitGroup, datasetRow *domain.Datasets, loggerUsecase 	*usecase.LoggerUsecase, connDB *gorm.DB, datasetUsecase *usecase.DatasetUsecase, question *domain.Question){
	
	errorRepository := repository.NewErrorRepository(connDB)
	errorUsecase := usecase.NewErrorUsecase(errorRepository)

	atomRepository := repository.NewAtomRepository(connDB)
	atomUsecase := usecase.NewAtomUsecase(atomRepository)

	defer wg.Done()

	agentOne := make(chan dto.ClientResponseDTO, 1)
	agentTwo := make(chan dto.ClientResponseDTO, 1)
	totalChannels := 2

	clientService := adapters.NewApiRequestService()
	clientUsecase := usecase.NewAPIRequestUsecase(clientService)


	QUESTION_ID := 1

	atom := domain.Atom{
		QuestionID:            QUESTION_ID,
		Question:              question.Question,
		AgentOneAnswer:        "",
		AgentTwoAnswer:        "",
		DatasetID:             int(datasetRow.ID),
		AtomSearched:          datasetRow.Atom,
		AtomFindedByAgentOne:  "",
		AtomFindedByAgentTwo:  "",
		AgentOneIsCorrect:     false,
		AgentTwoIsCorrect:     false,
		Failed:                false,
		ErrorID:               0,
		UpdatedAt:             time.Now(),
		CreatedAt:             time.Now(),
	}

	// agent one request
	if !datasetRow.MarkedByAgentOne{
		go func(){
			// gemini
			urlAgentOne := constants.URL_AGENT_ONE + globalCfg.AGENT_ONE_KEY
			
			body, err := agentOneJsonBody(parse(question.Question))

			headers := map[string]string{"Content-Type": "application/json"}

			if err != nil{
				loggerUsecase.Error("[-] Body from agent one contains error" + err.Error())
				return
			}
			
			response, err := clientUsecase.Post(urlAgentOne, headers, body)
			
			if err != nil{
				errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, err, constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, agentOne)
				atom.ErrorID = errID
				return
			}

			loggerUsecase.Info("[+] Status From agent one: " + response.Status)

			if response.StatusCode == http.StatusOK{
				datasetRow.MarkedByAgentOne = true
				datasetUsecase.UpdateMarkedByAgent(1, int(datasetRow.ID))

				status, responseStr := usecase.ResponseParser(response)

				if status{
					responseAgentOne := dto.ClientResponseDTO{
						Message: responseStr,
						Api: constants.AGENT_ONE,
					}

					agentOne <- responseAgentOne
				}else{
					errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Status Not ok :%s",responseStr), constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, agentOne)
					atom.ErrorID = errID
					return
				}
			}else{
				errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Fail to get success in request to agent one"), constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, agentOne)
				atom.ErrorID = errID
				return
			}
		}()
	}

	// agent two request
	if !datasetRow.MarkedByAgentTwo{
		go func(){
			body, err := agentTwoJsonBody(parse(question.Question))

			if err != nil{
				loggerUsecase.Error("[-] Body from agent one contains error" + err.Error())
				return
			}
			headers := map[string]string{"Content-Type": "application/json"}
			headers["Authorization"] = "Bearer " + globalCfg.AGENT_TWO_KEY

			response, err := clientUsecase.Post(constants.URL_AGENT_TWO, headers, body)

			if err != nil{
				errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, err, constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, agentTwo)
				atom.ErrorID = errID
				return
			}
			loggerUsecase.Info("[+] Status From agent two: " + response.Status)
			if response.StatusCode == http.StatusOK{
				datasetRow.MarkedByAgentTwo = true
				datasetUsecase.UpdateMarkedByAgent(2, int(datasetRow.ID))

				status, responseStr := usecase.ResponseParser(response)
				
				if status{
					responseAgentTwo := dto.ClientResponseDTO{
						Message: responseStr,
						Api: constants.AGENT_TWO,
					}

					agentTwo <- responseAgentTwo
				}else{
					errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("erro: %s", responseStr), constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, agentTwo)
					atom.ErrorID = errID
					return
				}
			}else{
				errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Fail to get success in request to agent two"), constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, agentTwo)
				atom.ErrorID = errID
				return
			}
		}()
	}

	var responseAgentOneDTO dto.ClientResponseDTO
	var responseAgentTwoDTO dto.ClientResponseDTO

	timeout := time.After(60 * time.Second)

	// select answers from mult channels
	for i := 0; i < totalChannels; i++{
		select{
		case responseAgentOne := <-agentOne:
			responseAgentOneDTO = responseAgentOne
		case responseAgentTwo := <-agentTwo:
			responseAgentTwoDTO = responseAgentTwo
		case <-timeout:
			loggerUsecase.Error("Timeout API")
			return
		}
	}

	atom.AgentOneAnswer = responseAgentOneDTO.Message
	atom.AgentTwoAnswer = responseAgentTwoDTO.Message

	atom.AgentOneIsCorrect = usecase.CheckIfAnswerContainsAtomOfConfusion(atom.AgentOneAnswer)
	atom.AgentTwoIsCorrect = usecase.CheckIfAnswerContainsAtomOfConfusion(atom.AgentOneAnswer)

	atom.AtomFindedByAgentOne = usecase.CheckWhatAtomOfConfusion(atom.AgentOneAnswer)
	atom.AtomFindedByAgentTwo = usecase.CheckWhatAtomOfConfusion(atom.AgentTwoAnswer)

	_, err := atomUsecase.Create(&atom)
	if err != nil{
		loggerUsecase.Error("[-] Fail to insert ATOM ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		return
	}
}

func agentOneJsonBody(content string) (string, error){
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": content,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(body)
    if err != nil {
        fmt.Println("Erro ao converter para JSON:", err)
        return "", err
    }

	return string(jsonData), nil
}

func agentTwoJsonBody(content string) (string, error){
	body := dto.RequestBody{
		Messages: []dto.Message{
			{
				Role:    "user",
				Content: content,
			},
		},
		Model: constants.MODEL,
	}

	jsonData, err := json.Marshal(body)
    if err != nil {
        fmt.Println("Erro ao converter para JSON:", err)
        return "", err
    }

	return string(jsonData), nil
}

func parse(message string) string{
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}