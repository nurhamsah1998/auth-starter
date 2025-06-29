package auth

type (
	AuthService struct{}
)

// / service handler untuk menginject servis ke controller
func AuthHandler() *AuthService {
	return &AuthService{}
}
