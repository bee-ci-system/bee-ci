package updater

type ghService struct{}

func (g ghService) NewGHService() *ghService {
	return &ghService{}
}
