package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
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
	globalConfig *config.Config
)


// execute is the main entry point of the application.
// It initializes the logger, loads configuration, establishes a database connection,
// sets up the repositories and use cases, and starts the task processing pool.
func execute() {

	startTime := time.Now()

	loggerUsecase := usecase.NewLoggerUsecase(usecase.LoggerConfig)

	loggerUsecase.Error("[+] Starting...")

	cfg, err := config.LoadConfig(".env")
	globalConfig = cfg

	if err != nil {
		loggerUsecase.Error("[-] Fail to load .env", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		os.Exit(1)
	}

	// database conn
	connDB, err := db.NewPostgresConn(cfg.DatabaseUrl)
	if err != nil {
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
	datasetUsecase *usecase.DatasetUsecase,
	questionUsecase *usecase.QuestionUsecase,
	loggerUsecase *usecase.LoggerUsecase,
	connDB *gorm.DB) {

	var wg sync.WaitGroup

	questions, err := questionUsecase.FindAll()
	if err != nil {
		loggerUsecase.Error("[-] Fail to list questions ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		os.Exit(1)
	}

	dataset, err := datasetUsecase.FindAll()
	if err != nil {
		loggerUsecase.Error("[-] Fail to list dataset ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
		os.Exit(1)
	}

	// channel creation
	tasksCh := make(chan *domain.Task, len(dataset)*len(questions))

	// workerpool
	for i := 0; i < constants.WORKERS; i++ {
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
func workerPool(tasksCh <-chan *domain.Task, wg *sync.WaitGroup, loggerUsecase *usecase.LoggerUsecase, connDB *gorm.DB, datasetUsecase *usecase.DatasetUsecase, questions []*domain.Question) {
	for task := range tasksCh {
		// atention + code + question.
		//task.Question.Question = constants.QUESTION_HEADER + removeHeadersOpenAI(task.Dataset.Class) + task.Question.Question

		task.Question.Question = constants.QUESTION_HEADER + task.Dataset.Class + task.Question.Question
		insertExecutor(wg, task.Dataset, loggerUsecase, connDB, datasetUsecase, task.Question)
	}
}

// insertExecutor handles the processing of a single dataset-question pair.
// It sends concurrent requests to two agents, handles responses or errors,
// logs errors, and persists valid results as Atom entities in the database.
func insertExecutor(wg *sync.WaitGroup, datasetRow *domain.Datasets, loggerUsecase *usecase.LoggerUsecase, connDB *gorm.DB, datasetUsecase *usecase.DatasetUsecase, question *domain.Question) {
	
	log.Printf("Request: %d\n", datasetRow.ID)
	errorRepository := repository.NewErrorRepository(connDB)
	errorUsecase := usecase.NewErrorUsecase(errorRepository)

	atomRepository := repository.NewAtomRepository(connDB)
	atomUsecase := usecase.NewAtomUsecase(atomRepository)

	defer wg.Done()

	questionOne := make(chan dto.ClientResponseDTO, 1)
	questionTwo := make(chan dto.ClientResponseDTO, 1)
	totalChannels := 2

	clientService := adapters.NewApiRequestService()
	clientUsecase := usecase.NewAPIRequestUsecase(clientService)

	atom := domain.Atom{
		QuestionID:   int(question.ID),
		Question:     question.Question,
		Answer:       "",
		DatasetID:    int(datasetRow.ID),
		AtomSearched: datasetRow.Atom,
		AtomFinded:   "",
		IsCorrect:    false,
		Failed:       false,
		ErrorID:      0,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}

	if int(question.ID) == constants.QUESTION_ONE_NUMBER {
		// Question One
		if !datasetRow.MarkedByAgentOne {
			go func() {
				urlAgentOne := constants.URL_AGENT_TWO + globalConfig.GEMINI_KEY
				
				body, err := jsonBody(parse(question.Question))

				headers := map[string]string{"Content-Type": "application/json"}

				if err != nil {
					loggerUsecase.Error("[-] Body from question one contains error" + err.Error())
					return
				}
				
				response, err := clientUsecase.Post(urlAgentOne, headers, body)

				if err != nil {
					log.Println("Error in question one " + err.Error())
					errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, err, constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, questionOne)
					atom.ErrorID = errID
					return
				}
				
				defer response.Body.Close()

				loggerUsecase.Info("[+] Status From question one: " + response.Status)

				if response.StatusCode == http.StatusOK {
					
					status, responseDto := usecase.ResponseParser(response)

					if status {
						responseAgentOne := responseDto

						questionOne <- responseAgentOne
					} else {
						log.Println("Error in question one status " + response.Status)
						errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Status Not ok :%s", responseDto.Candidates[0].Content.Parts[0].Text), constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, questionOne)
						atom.ErrorID = errID
						return
					}
				} else {
					errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Fail to get success in request of question one"), constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, questionOne)
					atom.ErrorID = errID
					return
				}
			}()
		}
	} else {
		if !datasetRow.MarkedByAgentTwo {
			go func() {
				body, err := jsonBody(parse(question.Question))

				if err != nil {
					loggerUsecase.Error("[-] Body from question two contains error" + err.Error())
					return
				}
				headers := map[string]string{"Content-Type": "application/json"}

				urlStr := constants.URL_AGENT_TWO + globalConfig.GEMINI_KEY
				response, err := clientUsecase.Post(urlStr, headers, body)

				if err != nil {
					log.Println("Error in question two " + err.Error())
					errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, err, constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, questionTwo)
					atom.ErrorID = errID
					return
				}

				defer response.Body.Close()

				loggerUsecase.Info("[+] Status From question two: " + response.Status)

				if response.StatusCode == http.StatusOK {
					
					status, responseDto := usecase.ResponseParser(response)

					if status {
						responseAgentTwo := responseDto
						questionTwo <- responseAgentTwo
					} else {
						errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("erro: %s", responseDto.Candidates[0].Content.Parts[0].Text), constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, questionTwo)
						atom.ErrorID = errID
						return
					}
				} else {
					errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Fail to get success in request of question two"), constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, questionTwo)
					atom.ErrorID = errID
					return
				}
			}()
		}
	}


	var responseQuestionOneDTO dto.ClientResponseDTO
	var responseQuestionTwoDTO dto.ClientResponseDTO

	timeout := time.After(time.Duration(constants.REQUEST_TIMEOUT_INTERVAL) * time.Second)

	// select answers from mult channels
	for i := 0; i < totalChannels; i++ {
		select {
		case responseQuestionOne := <-questionOne:
			responseQuestionOneDTO = responseQuestionOne
			atom.Answer = responseQuestionOneDTO.Candidates[0].Content.Parts[0].Text
			atom.IsCorrect = usecase.CheckIfAnswerContainsAtomOfConfusion(atom.Answer)
			atom.AtomFinded = usecase.CheckWhatAtomOfConfusion(atom.Answer)

			_, err := atomUsecase.Create(&atom)
			if err != nil {
				loggerUsecase.Error("[-] Fail to insert ATOM ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
				return
			}

			datasetRow.MarkedByAgentTwo = true
			datasetUsecase.UpdateMarkedByAgent(1, int(datasetRow.ID))

		case responseQuestionTwo := <-questionTwo:
			responseQuestionTwoDTO = responseQuestionTwo
			atom.Answer = responseQuestionTwoDTO.Candidates[0].Content.Parts[0].Text
			atom.IsCorrect = usecase.CheckIfAnswerContainsAtomOfConfusion(atom.Answer)
			atom.AtomFinded = usecase.CheckWhatAtomOfConfusion(atom.Answer)

			_, err := atomUsecase.Create(&atom)
			if err != nil {
				loggerUsecase.Error("[-] Fail to insert ATOM ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
				return
			}
			datasetRow.MarkedByAgentTwo = true
			datasetUsecase.UpdateMarkedByAgent(2, int(datasetRow.ID))

		case <-timeout:
			log.Println("Timeout")
			loggerUsecase.Error("API Timeout")
			return
		}
	}

}

func jsonBody(content string) (string, error) {
	body := dto.RequestBody{
		Contents: []dto.Content{
			{
				Parts: []dto.Part{
					{
						Text: content,
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

func parse(message string) string {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return err.Error()
	}
	return string(jsonData)
}

func removeHeadersOpenAI(code string) string {
	// Remove comments
	code = regexp.MustCompile(constants.SINGLE_COMMENTS).ReplaceAllString(code, "")
	// Remove m comments
	code = regexp.MustCompile(constants.MULTI_COMMENTS).ReplaceAllString(code, "")
	// Remove imports
	code = regexp.MustCompile(constants.IMPORTS).ReplaceAllString(code, "")
	// Remove (package)
	code = regexp.MustCompile(constants.HEADERS).ReplaceAllString(code, "")
	return code
}