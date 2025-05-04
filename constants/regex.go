package constants

const (
	SINGLE_COMMENTS string = `(?m)//[^\n]*\n`
	MULTI_COMMENTS string = `(?s)/\*.*?\*/`
	IMPORTS string = `(?m)^import\s+[^\n]+;\s*`
	HEADERS string = `(?m)^package\s+[^\n]+;\s*`
	ANSWER  string = `(?i)\byes\.\s+[A-Za-z ]+(?:,\s*[A-Za-z ]+)*\.`
)