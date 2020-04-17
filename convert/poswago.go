package convert

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/bearcherian/poswago/postman"
	"github.com/bearcherian/poswago/swagger"
)

var urlParamReg *regexp.Regexp = regexp.MustCompile(`\{\{([^\{\}]*)\}\}`)

// PostmanJSONToOpenAPI converts JSON data exported from Postman to a valid OpenAPI v3.0.0 spec
func PostmanJSONToOpenAPI(postmanSource []byte) (*swagger.OpenApiSpec, error) {
	var postmanSpec postman.Spec
	json.Unmarshal(postmanSource, &postmanSpec)
	return PostmanToSwagger(postmanSpec)
}

// PostmanToSwagger converts a postman spec type to an OpenAPI v3.0.0 spec
func PostmanToSwagger(postmanSpec postman.Spec) (*swagger.OpenApiSpec, error) {

	openAPISpec := swagger.OpenApiSpec{
		OpenAPI: "3.0.0",
		Info: swagger.Info{
			Title:       postmanSpec.Info.Name,
			Description: postmanSpec.Info.Description,
			Version:     "",
		},
		Paths: make(map[string]*swagger.PathItem),
	}

	for _, item := range postmanSpec.Items {
		pathItems := convertItemsToPaths(item)

		for k, v := range pathItems {
			fmt.Printf("Adding %s pathItem\n", k)
			if existingPathItem, ok := openAPISpec.Paths[k]; ok {
				fmt.Printf("pathItem %s already exists. merging...\n", k)
				openAPISpec.Paths[k] = mergePathItems(existingPathItem, v)
			} else {
				fmt.Printf("pathItem %s is new\n", k)
				openAPISpec.Paths[k] = v
			}
		}
	}

	return &openAPISpec, nil
}

func convertItemsToPaths(i postman.Item) map[string]*swagger.PathItem {
	paths := make(map[string]*swagger.PathItem)
	if i.Request != nil {
		fmt.Printf("ITEM %s request being processed\n", i.Name)
		path, pathItem := postmanItemToOARequest(i)
		if path != "" {

			fmt.Printf("ITEM %s has path %s\n", i.Name, path)
			paths[path] = pathItem
		}
	}

	for _, item := range i.Items {

		morePaths := convertItemsToPaths(item)

		for k, v := range morePaths {
			if _, ok := paths[k]; ok {
				// TODO update existing
				paths[k] = mergePathItems(paths[k], v)
			} else {
				paths[k] = v
			}
		}
	}

	return paths
}

// merges two PathItems. values in item1 take precedence
func mergePathItems(item1 *swagger.PathItem, item2 *swagger.PathItem) *swagger.PathItem {
	if item1.Summary == nil {
		item1.Summary = item2.Summary
	}
	//if item1.Description == nil {
	//	item1.Description = item2.Description
	//}

	item1.Parameters = mergeParameters(item1.Parameters, item2.Parameters)
	item1.Servers = mergeServers(item1.Servers, item2.Servers)

	// OPERATIONS
	if item1.Get == nil {
		item1.Get = item2.Get
	}
	if item1.Put == nil {
		item1.Put = item2.Put
	}
	if item1.Post == nil {
		item1.Post = item2.Post
	}
	if item1.Delete == nil {
		item1.Delete = item2.Delete
	}
	if item1.Options == nil {
		item1.Options = item2.Options
	}
	if item1.Head == nil {
		item1.Head = item2.Head
	}
	if item1.Trace == nil {
		item1.Trace = item2.Trace
	}
	if item1.Patch == nil {
		item1.Patch = item2.Patch
	}

	return item1
}

func mergeParameters(parameters1 []swagger.Parameter, parameters2 []swagger.Parameter) []swagger.Parameter {

	// map to track added parameters
	paramMap := make(map[string]bool)
	var newParameters []swagger.Parameter

	// add first set of parameters to map
	for _, p1 := range parameters1 {
		k := p1.Name
		if _, ok := paramMap[k]; !ok {
			newParameters = append(newParameters, p1)
			paramMap[k] = true
		}
	}

	// add 2nd set of parameters to map, skipping same names added from first set
	for _, p2 := range parameters2 {
		k := p2.Name
		if _, ok := paramMap[k]; !ok {
			newParameters = append(newParameters, p2)
			paramMap[k] = true
		}
	}

	return newParameters

}

func mergeServers(servers1 []*swagger.Server, servers2 []*swagger.Server) []*swagger.Server {
	// map to track added parameters
	serverMap := make(map[string]bool)
	var newServers []*swagger.Server

	// add first set of parameters to map
	for _, s1 := range servers1 {
		k := s1.URL
		if _, ok := serverMap[k]; !ok {
			newServers = append(newServers, s1)
			serverMap[k] = true
		}
	}

	// add 2nd set of parameters to map, skipping same names added from first set
	for _, s2 := range servers2 {
		k := s2.URL
		if _, ok := serverMap[k]; !ok {
			newServers = append(newServers, s2)
			serverMap[k] = true
		}
	}

	return newServers
}

func postmanItemToOARequest(item postman.Item) (string, *swagger.PathItem) {
	if item.Request == nil {
		return "", nil
	}

	r := item.Request

	p := fmt.Sprintf("/%s", strings.Join(r.Url.Path, "/"))

	p = urlParamReg.ReplaceAllString(p, "{$1}")

	s := swagger.PathItem{}

	s.Description = &r.Description
	s.Summary = &item.Name

	switch strings.ToLower(item.Request.Method) {
	case "get":
		s.Get = requestToPathOperation(item.Request)
	case "put":
		s.Put = requestToPathOperation(item.Request)
	case "post":
		s.Post = requestToPathOperation(item.Request)
	case "patch":
		s.Patch = requestToPathOperation(item.Request)
	case "trace":
		s.Trace = requestToPathOperation(item.Request)
	case "head":
		s.Head = requestToPathOperation(item.Request)
	case "options":
		s.Options = requestToPathOperation(item.Request)
	case "delete":
		s.Delete = requestToPathOperation(item.Request)
	}

	return p, &s

}

func requestToPathOperation(request *postman.Request) *swagger.PathOperation {
	o := swagger.PathOperation{
		Description: request.Description,
		Responses:   swagger.Responses{},
		OperationID: fmt.Sprintf("%s_%s", request.Method, request.Url.Raw),
	}

	// convert url parameters to swagger params
	params := urlParamReg.FindAllString(request.Url.Raw, -1)
	for _, p := range params {
		// skip the parameterized host url
		if strings.Index(request.Url.Raw, p) == 0 {
			continue
		}

		// strip brackets from the param name
		p = urlParamReg.ReplaceAllString(p, "$1")
		param := swagger.Parameter{
			Name: p,
			In:   "path",
		}
		param.Required = true

		o.Parameters = append(o.Parameters, param)
	}

	// convert query params to parameters
	for _, q := range request.Url.Queries {
		param := swagger.Parameter{}

		param.Description = q.Description
		param.Name = q.Key
		param.Example = q.Value
		param.In = "query"

		o.Parameters = append(o.Parameters, param)
	}

	return &o
}
