package sessions

//go:generate moq -out ../mock/sessions_repo_db_mock.go -pkg mock . SessionRepository:MockSessionRepository
type SessionRepository interface {
	CreateSession(int, string) error
	DeleteSession(string) error
	GetUserIdByCookie(string) (int, error)
}
