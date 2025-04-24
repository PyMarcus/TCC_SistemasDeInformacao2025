package repositories

type Question interface{
	FindAll() ([]*domain.Question, error)
}