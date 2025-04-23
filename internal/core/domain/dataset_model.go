package domain

type DatasetModel struct {
	ID               uint   `gorm:"primaryKey"`
	Class            string `gorm:"column:class"`
	Atom             string `gorm:"column:atom"`
	Snippet          string `gorm:"column:snippet"`
	Line             string `gorm:"column:line"`
	GithubLink       string `gorm:"column:github_link"`
	StatusCode       string `gorm:"column:status_code"`
	MarkedByAgentOne bool   `gorm:"column:marked_by_agent_one"`
	MarkedByAgentTwo bool   `gorm:"column:marked_by_agent_two"`
}
