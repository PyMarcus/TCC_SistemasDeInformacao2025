package dto

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type RequestBody struct {
    Messages []Message `json:"messages"`
    Model    string    `json:"model"`
}