package dto

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestOpenAiBody struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Temperature float32   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
	Stream      bool      `json:"stream"`
}

type RequestBody struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}