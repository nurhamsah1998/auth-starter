package profile

type (
	ProfileService struct{}
)

// / service handler untuk menginject servis ke controller
func ProfileHandler() *ProfileService {
	return &ProfileService{}
}
