package usecase

import (
	"io"
	"net/http"
)

func ResponseParser(response *http.Response) (bool, string){
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return false, err.Error()
	}
	defer response.Body.Close()
	
	bodyStr := string(bodyBytes)
	return true, bodyStr
}