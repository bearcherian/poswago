package swagger

import "github.com/qri-io/jsonschema"

type OpenApiSpec struct {
	OpenAPI string               `json:"openapi"`
	Info    Info                 `json:"info"`
	Servers []*Server            `json:"servers,omitempty"`
	Paths   map[string]*PathItem `json:"paths"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type PathItem struct {
	Ref         *string        `json:"$ref,omitempty"`
	Summary     *string        `json:"summary"`
	Description *string        `json:"description"`
	Get         *PathOperation `json:"get,omitempty"`
	Put         *PathOperation `json:"put,omitempty"`
	Post        *PathOperation `json:"post,omitempty"`
	Delete      *PathOperation `json:"delete,omitempty"`
	Options     *PathOperation `json:"options,omitempty"`
	Head        *PathOperation `json:"head,omitempty"`
	Patch       *PathOperation `json:"patch,omitempty"`
	Trace       *PathOperation `json:"trace,omitempty"`
	Servers     []*Server      `json:"servers,omitempty"`
	Parameters  []Parameter    `json:"parameters,omitempty"`
}

type PathOperation struct { //todo
	Tags         []string              `json:"tags,omitempty"`
	Summary      string                `json:"summary,omitempty"`
	Description  string                `json:"description,omitempty"`
	ExternalDocs ExternalDocumentation `json:"externalDocs,omitempty"`
	OperationID  string                `json:"operationId,omitempty"`
	Parameters   []Parameter           `json:"parameters,omitempty"`
	Responses    Responses             `json:"responses"`
	Callbacks    map[string]Callback   `json:"callback,omitempty"`
	Deprecated   bool                  `json:"deprecated,omitempty"`
	Security     []Security            `json:"security,omitempty"`
	Servers      []*Server             `json:"servers,omitempty"`
}

type Callback map[string]PathItem

type Security map[string]string

type Header struct {
	Description     string             `json:"description"`
	Required        bool               `json:"required,omitempty"`
	Deprecated      bool               `json:"deprecated"`
	AllowEmptyValue bool               `json:"allowEmptyValue"`
	Style           string             `json:"style,omitempty"`
	Explode         bool               `json:"explode"`
	Schema          Schema             `json:"schema,omitempty"`
	Example         interface{}        `json:"example,omitempty"`
	Examples        map[string]Example `json:"examples,omitempty"`
}

type Example struct {
	Summary       string      `json:"summary"`
	Description   string      `json:"description"`
	Value         interface{} `json:"value"`
	ExternalValue string      `json:"externalValue"`
}

type Schema jsonschema.RootSchema

type Parameter struct {
	Header
	Name string `json:"name"`
	In   string `json:"in"`
}

type ExternalDocumentation struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Responses map[string]Response

type Response struct {
	Description string               `json:"description"`
	Headers     map[string]Header    `json:"headers"`
	Content     map[string]MediaType `json:"content"`
	Links       map[string]Link      `json:"links"`
}

type MediaType struct {
	jsonschema.RootSchema
	Example  interface{}         `json:"example"`
	Examples map[string]Example  `json:"examples"`
	Encoding map[string]Encoding `json:"encoding"`
}

type Encoding struct {
	ContentType   string
	Headers       map[string]Header
	Style         string
	Explode       bool
	AllowReserved bool
}

type Link struct {
	OperationRef string
	OperationId  string
	Parameters   map[string]interface{}
	RequestBody  interface{}
	Description  string
	Server       Server
}
