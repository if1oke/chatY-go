package auth

type IAuthService interface {
	Login(username, password string) (bool, string)
}
