package usecase

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http/dto"
)

func ResponseParser(response *http.Response) (bool, dto.ClientResponseDTO){
	var responseDTO dto.ClientResponseDTO

	err := json.NewDecoder(response.Body).Decode(&responseDTO)
	if err != nil {
		log.Println(err)
		return false, responseDTO
	}
	return true, responseDTO
}