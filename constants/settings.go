package constants

const (
	LOG_DIR                  string = "./logs"
	LOG_NAME                 string = "atoms_of_code.log"
	HTTP_GET_METHOD          string = "GET"
	HTTP_POST_METHOD         string = "POST"
	ERROR_DESCRIPTION        string = "error_description"
	ELAPSED_TIME			 string = "elapsed_time(s)"
	QUESTION				 string = "question"
	DATASET_ID				 string = "dataset_id"
	APPL_JSON				 string = "application/json"
	CT_TYPE					 string = "Content-Type"
	AUTH					 string = "Authorization"
	URL_AGENT_ONE			 string = ""
	URL_AGENT_TWO			 string = ""
	AGENT_ONE				 string = "AGENT_ONE"
	AGENT_TWO				 string = "AGENT_TWO"
	URL						 string = "URL"
	STATUS_CODE_STR			 string = "STATUS_CODE"
	QUESTION_HEADER 		 string = "Analyze the Java code below, as I will ask you about it later:"
	REQUEST_TIMEOUT_INTERVAL int    = 2  // min
	WORKERS					 int 	= 10 // max goroutines, because GEMINI allows 10 requests per min
)
