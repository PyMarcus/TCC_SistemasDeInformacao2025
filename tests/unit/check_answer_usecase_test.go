package unit

import (
	"strings"
	"testing"

	"github.com/PyMarcus/TCC_SistemasDeInformacao2025/internal/core/usecase"
	"github.com/stretchr/testify/assert"
)

type testCaseCheckAnswer struct {
	validAnswer bool
	answer      string
}

func createTestCasesToCheckAnswer() []*testCaseCheckAnswer{
	return []*testCaseCheckAnswer{
		{validAnswer: true, answer: "Yes.Conditional Operator,repurposed variables"},
		{validAnswer: true, answer: "Yes.Conditional Operator,repurposed variables"},
		{validAnswer: true, answer: "Yes.Conditional Operator,repurposed variables"},
		{validAnswer: true, answer: "Yes.Conditional Operator,repurposed variables"},
		{validAnswer: true, answer: "Yes.Conditional Operator,repurposed variables"},
		{validAnswer: false, answer: "No."},
		{validAnswer: false, answer: "das.Conditional Operator,repurposed variables"},
		{validAnswer: false, answer: "Yes."},
		{validAnswer: false, answer: "Yes.Chapolim, Fuzz, doxing"},
	}
}

func TestCheckIfAnswerContainsAtomOfConfusion(t *testing.T) {
	
	result := usecase.CheckIfAnswerContainsAtomOfConfusion("Yes.Conditional Operator,repurposed variables", "repurposed variables")
	assert.True(t, result)

	result = usecase.CheckIfAnswerContainsAtomOfConfusion("Yes.Conditional Operator,repurposed variables", "doxing")
	assert.False(t, result)

	result = usecase.CheckIfAnswerContainsAtomOfConfusion("Yes. Type Conversion, Constant Variables.", "Type Conversion")
	assert.True(t, result)
	
}

func TestCheckWhatAtomOfConfusion(t *testing.T) {
	for _, test := range createTestCasesToCheckAnswer(){
		result := usecase.CheckWhatAtomOfConfusion(test.answer)
		var value string
		if strings.Contains(strings.ToLower(test.answer), "yes"){
			value = strings.ToLower(strings.Split(test.answer, ".")[1])
		}else{
			value = strings.ToLower(test.answer)
		}
		assert.Equal(t, result, value)
	}
}
