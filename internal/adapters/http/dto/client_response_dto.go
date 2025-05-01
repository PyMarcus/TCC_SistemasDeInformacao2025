package dto

type ClientResponseOpenAIDTO struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"` 
		} `json:"message"`
	} `json:"choices"`
}

type ClientResponseDTO struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}
