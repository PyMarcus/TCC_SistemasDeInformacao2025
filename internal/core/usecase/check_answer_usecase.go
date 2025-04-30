package usecase

import "strings"

var (
	atoms = []string{
		"infix operator precedence",
		"post-increment/decrement",
		"pre-increment/decrement",
		"constant variables",
		"remove indentation",
		"conditional operator",
		"arithmetic as logic",
		"logic as control flow",
		"repurposed variables",
		"dead, unreachable, repeated",
		"change of literal encoding",
		"omitted curly braces",
		"type conversion",
		"indentation",
	}
)


func CheckIfAnswerContainsAtomOfConfusion(answer string) bool{
	answersForm := strings.ToLower(answer)

	if strings.Contains(answersForm, "yes"){
		for _, atom := range atoms{
			if strings.Contains(answersForm, atom){
				return true
			}
		}
	}
	return false
}

func CheckWhatAtomOfConfusion(answer string) string{
	answersForm := strings.ToLower(answer)
	if strings.Contains(answersForm, "yes."){
		return strings.Split(answersForm, ".")[1]
	}
	return answersForm
}