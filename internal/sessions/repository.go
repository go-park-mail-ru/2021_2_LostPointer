package sessions

type SessionRepository interface {
	CreateSession(int, string) error
	DeleteSession(string) error
	GetUserIdByCookie(string) (int, error)
}
