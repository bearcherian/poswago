package postman

type Request struct {
	Method      string   `json:"method"`
	Headers     []Header `json:"header"`
	Url         URL      `json:"url"`
	Description string   `json:"description"`
}

type Header struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Cookie struct {
	
}

type URL struct {
	Raw     string   `json:"raw"`
	Host    []string `json:"host"`
	Path    []string `json:"path"`
	Queries []Query  `json:"query"`
}

type Query struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Disabled    bool   `json:"disabled"`
}
