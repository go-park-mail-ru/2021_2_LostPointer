package authorization

type AuthStorage interface {
	CreateSession(int64, string) error
	GetUserByPassword(string, string) (int64, error)
}
