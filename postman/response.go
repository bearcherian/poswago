package postman

type Response struct {
	Name                    string
	OriginalRequest         OriginalRequest
	Status                  string
	Code                    int
	PostmanePreviewLanguage string
	Header                  *[]Header
	Cookie                  *[]Cookie
	Body                    string
}

type OriginalRequest struct {
	Method string
	Header *[]Header
	URL    URL
}
