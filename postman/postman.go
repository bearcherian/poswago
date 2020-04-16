package postman

type Spec struct {
	Info  Info   `json:"info"`
	Items []Item `json:"item"`
}
type Info struct {
	PostmanID   string `json:"_postman_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
}

type Item struct {
	Name    string  `json:"name"`
	Request *Request `json:"request,omitempty"`
	Items   []Item  `json:"item"`
}
