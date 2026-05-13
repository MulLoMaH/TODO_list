package web_service

type WebService struct {
	webRepoitory WebRepoitory
}

type WebRepoitory interface {
	GetFile(filePath string) ([]byte, error)
}

func NewWebService(webRepoitory WebRepoitory) *WebService {
	return &WebService{
		webRepoitory: webRepoitory,
	}
}
