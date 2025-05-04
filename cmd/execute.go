package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
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
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	globalConfig *config.Config
	retryInterval = time.Second

    globalPause = struct{
	mutex sync.Mutex
	conditional *sync.Cond
	paused bool
	}{
		paused: false,
	}
)

func triggerGlobalPause(duration time.Duration){
	color.Cyan("Global pause activated!")
	globalPause.mutex.Lock()
	globalPause.paused = true 
	globalPause.mutex.Unlock()

	go func() {
		time.Sleep(duration)

		globalPause.mutex.Lock()
		globalPause.paused = false
		globalPause.mutex.Unlock()
		globalPause.conditional.Broadcast() 
	}()
}

func init(){
	globalPause.conditional = sync.NewCond(&globalPause.mutex)
}

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
		globalPause.mutex.Lock()
		for globalPause.paused {
			globalPause.conditional.Wait()
		}
		globalPause.mutex.Unlock()


		task.Question.Question = constants.QUESTION_HEADER + task.Dataset.Class + task.Question.Question
		insertExecutor(wg, task.Dataset, loggerUsecase, connDB, datasetUsecase, task.Question)
	}
}

// insertExecutor handles the processing of a single dataset-question pair.
// It sends concurrent requests to two agents, handles responses or errors,
// logs errors, and persists valid results as Atom entities in the database.
func insertExecutor(wg *sync.WaitGroup, datasetRow *domain.Datasets, loggerUsecase *usecase.LoggerUsecase, connDB *gorm.DB, datasetUsecase *usecase.DatasetUsecase, question *domain.Question) {

	activeChannels := 0

	loggerUsecase.Info(fmt.Sprintf("Request: %d\n", datasetRow.ID))
	color.Green(fmt.Sprintf("Request: %d\n", datasetRow.ID))

	errorRepository := repository.NewErrorRepository(connDB)
	errorUsecase := usecase.NewErrorUsecase(errorRepository)

	atomRepository := repository.NewAtomRepository(connDB)
	atomUsecase := usecase.NewAtomUsecase(atomRepository)

	defer wg.Done()

	questionOne := make(chan dto.ClientResponseDTO, 1)
	questionTwo := make(chan dto.ClientResponseDTO, 1)

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

	if int(question.ID) == constants.QUESTION_ONE_NUMBER && !datasetRow.MarkedByAgentOne {
		activeChannels++
		go questionOneFn(question,loggerUsecase,  &atom, clientUsecase, errorUsecase,questionOne)
		
	} 
	if int(question.ID) != constants.QUESTION_ONE_NUMBER && !datasetRow.MarkedByAgentTwo {
		activeChannels++
		go questionTwoFn(question,loggerUsecase,  &atom, clientUsecase, errorUsecase,questionTwo)
	}

	timeout := time.After(time.Duration(constants.REQUEST_TIMEOUT_INTERVAL) * time.Second)
	color.Blue(fmt.Sprintf("QuestionID: %d", atom.QuestionID))
	// select answers from mult channels
	for i := 0; i < activeChannels; i++ {
		select {
		case responseQuestionOne := <-questionOne:
			handleResponse(responseQuestionOne, &atom, datasetRow, atomUsecase, datasetUsecase, 1, loggerUsecase)

		case responseQuestionTwo := <-questionTwo:
			handleResponse(responseQuestionTwo, &atom, datasetRow, atomUsecase, datasetUsecase, 2, loggerUsecase)

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

// removeHeadersOpenAI used to limited open ai tokens
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

func questionOneFn(
	question *domain.Question,
	loggerUsecase *usecase.LoggerUsecase,
	atom *domain.Atom,
	clientUsecase *usecase.APIRequestUsecase,
	errorUsecase *usecase.ErrorUsecase,
	questionOne chan dto.ClientResponseDTO){
	retryInterval = 0

	urlAgentOne := constants.URL_AGENT_TWO + globalConfig.GEMINI_KEY

	body, err := jsonBody(parse(question.Question))

	headers := map[string]string{"Content-Type": "application/json"}

	if err != nil {
		loggerUsecase.Error("[-] Body from question one contains error" + err.Error())
		return
	}

	for i := 0; i < constants.MAX_RETRIES; i++{
		response, _ := clientUsecase.Post(urlAgentOne, headers, body)


		defer response.Body.Close()

		loggerUsecase.Info("[+] Status From question one: " + response.Status)

		if response.StatusCode == http.StatusOK {

			status, responseDto := usecase.ResponseParser(response, loggerUsecase, 1)

			if status {
				responseAgentOne := responseDto
				color.Blue("OK -> Question One")
				questionOne <- responseAgentOne
			} else {
				color.Red("Not OK -> Question One: " + response.Status)
				log.Println("Error in question one status " + response.Status)
				errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Status Not ok :%s", responseDto.Candidates[0].Content.Parts[0].Text), constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, questionOne)
				atom.ErrorID = errID
				return
			}
		}else if response.StatusCode == http.StatusTooManyRequests{
			color.Red("[-] Too many requests! Waiting...")
			bodyBytes, _ := io.ReadAll(response.Body)
			waitDelay(bodyBytes)

		}else if response.StatusCode == http.StatusServiceUnavailable {
			color.Red("Not OK -> Question One: Server overloaded, retrying")

			log.Printf("Server overloaded, retrying... attempt %d", i+1)
			time.Sleep(retryInterval)
			retryInterval *= 4 

		} else {
			color.Red("Not OK! Status: " + response.Status)
			errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Fail to get success in request of question one"), constants.AGENT_ONE, constants.URL_AGENT_ONE, response.Status, questionOne)
			atom.ErrorID = errID
			return
		}
	}
}

func waitDelay(bodyBytes []byte){

		var data map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &data); err == nil {
			details := data["error"].(map[string]interface{})["details"].([]interface{})
			for _, d := range details {
				detailMap := d.(map[string]interface{})
				if detailMap["@type"] == "type.googleapis.com/google.rpc.RetryInfo" {
					retryDelay := detailMap["retryDelay"].(string)
					color.Yellow("Waiting: " + retryDelay)

					dur, err := time.ParseDuration(retryDelay)
					if err == nil {
						color.Yellow("Waiting:...")
						time.Sleep(dur)
						triggerGlobalPause(dur)
					}else{
						time.Sleep(20 * time.Second)
					}
				}
			}
		}
}

func questionTwoFn(
	question *domain.Question,
	loggerUsecase *usecase.LoggerUsecase,
	atom *domain.Atom,
	clientUsecase *usecase.APIRequestUsecase,
	errorUsecase *usecase.ErrorUsecase,
	questionTwo chan dto.ClientResponseDTO){
	retryInterval = 0

	body, err := jsonBody(parse(question.Question))

	if err != nil {
		loggerUsecase.Error("[-] Body from question two contains error" + err.Error())
		return
	}
	headers := map[string]string{"Content-Type": "application/json"}

	urlStr := constants.URL_AGENT_TWO + globalConfig.GEMINI_KEY
	
	for i := 0; i < constants.MAX_RETRIES; i++{
		response, _ := clientUsecase.Post(urlStr, headers, body)

	
		defer response.Body.Close()

		loggerUsecase.Info("[+] Status From question two: " + response.Status)

		if response.StatusCode == http.StatusOK {

			status, responseDto := usecase.ResponseParser(response, loggerUsecase, 2)
	
			if status {
				responseAgentTwo := responseDto
				color.Blue("OK -> Question Two")
				questionTwo <- responseAgentTwo
			} else {
				errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("erro: %s", responseDto.Candidates[0].Content.Parts[0].Text), constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, questionTwo)
				atom.ErrorID = errID
				return
			}
		}else if response.StatusCode == http.StatusTooManyRequests{
			color.Red("[-] Too many requests! Waiting...")
			bodyBytes, _ := io.ReadAll(response.Body)
			waitDelay(bodyBytes)

		}else if response.StatusCode == http.StatusServiceUnavailable {
			log.Printf("Server overloaded, retrying... attempt %d", i+1)
			time.Sleep(retryInterval)
			retryInterval *= 4 
		}else {
			log.Printf("[-] Error status %s", response.Status)
			errID := usecase.HandleAgentError(loggerUsecase, errorUsecase, fmt.Errorf("[-] Fail to get success in request of question two"), constants.AGENT_TWO, constants.URL_AGENT_TWO, response.Status, questionTwo)
			atom.ErrorID = errID
			return
		}
	}
}

func handleResponse(response dto.ClientResponseDTO, atom *domain.Atom, datasetRow *domain.Datasets, atomUsecase *usecase.AtomUsecase, datasetUsecase *usecase.DatasetUsecase, agent int, loggerUsecase *usecase.LoggerUsecase) {
	atom.Answer = response.Candidates[0].Content.Parts[0].Text
	atom.IsCorrect = usecase.CheckIfAnswerContainsAtomOfConfusion(atom.Answer, atom.AtomSearched)
	atom.AtomFinded = usecase.CheckWhatAtomOfConfusion(atom.Answer)
	if atom.Answer != "" && strings.Contains(strings.ToLower(atom.Answer), "yes"){
		_, err := atomUsecase.Create(atom)
		if err != nil {
			color.Red("[-] Fail to insert ATOM " + err.Error())
			loggerUsecase.Error("[-] Fail to insert ATOM ", zap.String(constants.ERROR_DESCRIPTION, err.Error()))
			return
		}
		if agent == 1 {
			datasetRow.MarkedByAgentOne = true
		} else {
			datasetRow.MarkedByAgentTwo = true
		}
		datasetUsecase.UpdateMarkedByAgent(agent, int(datasetRow.ID))
	}else{
		loggerUsecase.Error("EMPTY ANSWER? " + response.Candidates[0].Content.Parts[0].Text)
		color.Yellow(fmt.Sprintf("Answer for question: %d is empty", agent))
	}
}