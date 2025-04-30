package dto

type Part struct {
    Text string `json:"text"`
}

type Content struct {
    Parts []Part `json:"parts"`
}

type Body struct {
    Contents []Content `json:"contents"`
}