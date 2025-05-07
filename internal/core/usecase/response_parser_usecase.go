package usecase

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/constants"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/adapters/http/dto"
	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/domain"
)

func ResponseParser(response *http.Response, logUsecase *LoggerUsecase, question int, atom *domain.Atom) (bool, dto.ClientResponseDTO){
	var responseDTO dto.ClientResponseDTO

	err := json.NewDecoder(response.Body).Decode(&responseDTO)
	if err != nil {
		log.Println(err)
		return false, responseDTO
	}
	text := responseDTO.Candidates[0].Content.Parts[0].Text
	logUsecase.Info("[+] RESPONSE OK: " + text)
	responseDTO.Candidates[0].Content.Parts[0].Text = text

	if question == 1{
		responseRe := regexp.MustCompile(constants.ANSWER)
		results := responseRe.FindAllString(text, -1)
		concat := strings.Join(results, " ")
		atom.AtomFinded = concat
	}else{
		atom.AtomFinded = text
	}

	return true, responseDTO
}