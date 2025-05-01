package constants

const (
	LOG_DIR                  string = "./logs"
	LOG_NAME                 string = "atoms_of_code.log"
	HTTP_GET_METHOD          string = "GET"
	HTTP_POST_METHOD         string = "POST"
	MODEL                    string = "deepseek-r1-distill-qwen-7b"
	ERROR_DESCRIPTION        string = "error_description"
	ELAPSED_TIME			 string = "elapsed_time(s)"
	QUESTION				 string = "question"
	DATASET_ID				 string = "dataset_id"
	APPL_JSON				 string = "application/json"
	CT_TYPE					 string = "Content-Type"
	AUTH					 string = "Authorization"
	URL_AGENT_ONE			 string = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="
	URL_AGENT_TWO			 string = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="
	AGENT_ONE				 string = "QUESTION_ONE"
	AGENT_TWO				 string = "QUESTION_TWO"
	URL						 string = "URL"
	STATUS_CODE_STR			 string = "STATUS_CODE"
	QUESTION_HEADER 		 string = "Analyze the Java code below, as I will ask you about it later:"
	REQUEST_TIMEOUT_INTERVAL int    = 120  // 2min 
	WORKERS					 int 	= 3 // max goroutines, because i'm using gemini 3rpm.
	QUESTION_ONE_NUMBER		 int    = 1
	TEMPERATURE				 float32 = 0.7
	MAX_TOKENS				 int     = -1
	STREAM                   bool    = false
)
