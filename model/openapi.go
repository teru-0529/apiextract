package model

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
)

type RefType string

const (
	SCHEMAS       RefType = "schemas"
	RESPONSES     RefType = "responses"
	PARAMETERS    RefType = "parameters"
	REQUESTBODIES RefType = "requestBodies"
	HEADERS       RefType = "headers"
	NONE          RefType = "raw"
)

type ApiBase struct {
	Openapi string
	Info    Info     `yaml:"info"`
	Servers []Server `yaml:"servers"`
	Tags    []Tag    `yaml:"tags"`
}

type Info struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

type Server struct {
	Url         string `yaml:"url"`
	Description string `yaml:"description"`
}

type Tag struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// --------------

type Path struct {
	url          string
	method       string
	OperationId  string      `yaml:"operationId"`
	Summary      string      `yaml:"summary"`
	Description  string      `yaml:"description"`
	Tags         []string    `yaml:"tags"`
	Parameters   []Parameter `yaml:"parameters"`
	requestBody  RequestBody
	responses    []Response
	ExternalDocs ExternalDoc `yaml:"externalDocs"`
}

type Parameter struct {
	Description string `yaml:"description"`
	Ref         string `yaml:"$ref"`
}

type RequestBody struct {
	exist       bool
	description string
	ref         string
	content     Content
}

type Response struct {
	status      string
	description string
	ref         string
	content     Content
	headers     []ResponseHeader
}

type Content struct {
	exist bool
	ref   string
}

type ResponseHeader struct {
	name        string
	description string
	ref         string
}

type ExternalDoc struct {
	Description string `yaml:"description"`
	Url         string `yaml:"url"`
}

type Component struct {
}

// --------------

func NewOpenApi(filename string) (*ApiBase, *[]Path, *[]Component, error) {
	// INFO: read
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cannot read file: %s", err.Error())
	}

	// INFO: baseItem
	var apiBase ApiBase
	err = yaml.Unmarshal([]byte(file), &apiBase)
	if err != nil {
		return nil, nil, nil, err
	}

	// INFO: PathItem
	yamlPath, err := yamlChild([]byte(file), "paths")
	if err != nil {
		return nil, nil, nil, err
	}
	paths, err := pathItems(yamlPath)
	if err != nil {
		return nil, nil, nil, err
	}

	// TODO: componentTtem

	return &apiBase, paths, nil, nil
}

// ----+----+----+----+----+----+----+----+----+----+----+

func (path *Path) ToPathArray() []string {
	return []string{
		strings.Join(path.Tags, ","),
		path.url,
		path.method,
		path.OperationId,
		path.Summary,
		path.Description,
		strconv.Itoa(len(path.Parameters)),
		strconv.FormatBool(path.requestBody.exist),
		strings.Join(lo.Map(path.responses, func(x Response, _ int) string { return x.status }), ","),
		strconv.FormatBool(path.ExternalDocs.Url != ""),
	}
}

// func (path *Path) ToIoArray() [][]string { // TODO:
// 	return []string{
// 		strings.Join(path.Tags, ","),
// 		path.url,
// 		path.method,
// 		path.OperationId,
// 		path.Summary,
// 		path.Description,
// 		fmt.Sprint(10),
// 		"path.requestBody.description",
// 		"path.responses",
// 		"lo.Ternary(path.hasExternalDocs, ",
// 	}
// }

func (parameter *Parameter) refType() RefType {
	if parameter.Ref != "" {
		return PARAMETERS
	} else {
		return NONE
	}
}

func (requestBody *RequestBody) refType() RefType {
	if requestBody.ref != "" {
		return REQUESTBODIES
	} else {
		return requestBody.content.refType()
	}
}

func (response *Response) refType() RefType {
	if response.ref != "" {
		return RESPONSES
	} else {
		return response.content.refType()
	}
}

func (content *Content) refType() RefType {
	if content.ref != "" {
		return SCHEMAS
	} else {
		return NONE
	}
}

func (header *ResponseHeader) refType() RefType {
	if header.ref != "" {
		return HEADERS
	} else {
		return NONE
	}
}

// ----+----+----+----+----+----+----+----+----+----+----+

func pathItems(yamlData []byte) (*[]Path, error) {
	var pathMap map[string]interface{}
	err := yaml.Unmarshal(yamlData, &pathMap)
	if err != nil {
		return nil, err
	}

	var paths []Path
	for url := range pathMap {
		methodMap, err := mapValue(pathMap, url)
		if err != nil {
			return nil, err
		}
		for method := range methodMap {
			path, err := newPath(methodMap, method)
			if err != nil {
				return nil, err
			}

			methpdValue, err := mapValue(methodMap, method)
			if err != nil {
				return nil, err
			}

			requestBody, err := newRequestBody(methpdValue)
			if err != nil {
				return nil, err
			}
			path.requestBody = *requestBody

			responses, err := newResponses(methpdValue)
			if err != nil {
				return nil, err
			}
			path.responses = *responses

			path.url = url
			path.method = method
			paths = append(paths, *path)
		}
	}

	sort.SliceStable(paths, func(x, y int) bool {
		a, b := paths[x], paths[y]
		switch {
		case a.url < b.url:
			return true
		case a.url > b.url:
			return false
		default:
			return a.method < b.method
		}
	})

	return &paths, nil
}

// ----+----+----+----+----+----+----+----+----+----+----+

func newRequestBody(methodValue map[string]interface{}) (*RequestBody, error) {
	// requestBody属性が未定義
	if _, ok := methodValue["requestBody"]; !ok {
		return &RequestBody{}, nil
	}

	requestMap, err := mapValue(methodValue, "requestBody")
	if err != nil {
		return nil, err
	}
	requestbody := RequestBody{exist: true}
	if _, ok := requestMap["description"]; ok {
		requestbody.description = requestMap["description"].(string)
	}

	if _, ok := requestMap["$ref"]; ok { // requestBody:$ref
		requestbody.ref = requestMap["$ref"].(string)
	} else {
		content, err := newContent(requestMap)
		if err != nil {
			return nil, err
		}
		requestbody.content = *content
	}

	return &requestbody, nil
}

func newResponses(methodValue map[string]interface{}) (*[]Response, error) {
	// reponses属性が未定義
	if _, ok := methodValue["responses"]; !ok {
		return nil, errors.New("responses must exist")
	}
	responsesMap, err := mapValue(methodValue, "responses")
	if err != nil {
		return nil, err
	}

	var responses []Response
	for status := range responsesMap {
		responseMap, err := mapValue(responsesMap, status)
		if err != nil {
			return nil, err
		}

		response := Response{status: status}
		if _, ok := responseMap["description"]; ok {
			response.description = responseMap["description"].(string)
		}

		if _, ok := responseMap["$ref"]; ok { // response:$ref
			response.ref = responseMap["$ref"].(string)
		} else {
			content, err := newContent(responseMap)
			if err != nil {
				return nil, err
			}
			response.content = *content

			responseheader, err := newResponseHeaders(responseMap)
			if err != nil {
				return nil, err
			}
			response.headers = *responseheader
		}

		responses = append(responses, response)
	}

	sort.Slice(responses, func(x, y int) bool {
		return responses[x].status < responses[y].status
	})

	return &responses, nil
}

// ----+----+----+----+----+----+----+----+----+----+----+

func mapValue(mapData map[string]interface{}, key string) (map[string]interface{}, error) {
	yamlData, err := yaml.Marshal(mapData[key])
	if err != nil {
		return nil, err
	}
	var distData map[string]interface{}
	err = yaml.Unmarshal(yamlData, &distData)
	if err != nil {
		return nil, err
	}
	return distData, nil
}

func yamlChild(yamlData []byte, key string) ([]byte, error) {
	var mapData map[string]interface{}
	err := yaml.Unmarshal(yamlData, &mapData)
	if err != nil {
		return nil, err
	}

	distData, err := yaml.Marshal(mapData[key])
	if err != nil {
		return nil, err
	}
	return distData, nil
}

func newPath(mapData map[string]interface{}, key string) (*Path, error) {
	yamlData, err := yaml.Marshal(mapData[key])
	if err != nil {
		return nil, err
	}
	var path Path
	err = yaml.Unmarshal(yamlData, &path)
	if err != nil {
		return nil, err
	}
	return &path, nil
}

func newContent(mapData map[string]interface{}) (*Content, error) {
	if _, ok := mapData["content"]; !ok {
		return &Content{exist: false}, nil // no_content
	}

	contentMap, err := mapValue(mapData, "content")
	if err != nil {
		return nil, err
	}
	jsonMap, err := mapValue(contentMap, "application/json")
	if err != nil {
		return nil, err
	}
	schemaMap, err := mapValue(jsonMap, "schema")
	if err != nil {
		return nil, err
	}

	if _, ok := schemaMap["$ref"]; ok { // schema:$ref
		return &Content{exist: true, ref: schemaMap["$ref"].(string)}, nil
	} else {
		return &Content{exist: true}, nil
	}
}

func newResponseHeaders(mapData map[string]interface{}) (*[]ResponseHeader, error) {
	var headers []ResponseHeader
	if _, ok := mapData["headers"]; !ok {
		return &headers, nil // no_headers
	}

	headerMap, err := mapValue(mapData, "headers")
	if err != nil {
		return nil, err
	}
	for name := range headerMap {
		header := ResponseHeader{name: name}
		headerValue, err := mapValue(headerMap, name)
		if err != nil {
			return nil, err
		}

		if _, ok := headerValue["description"]; ok {
			header.description = headerValue["description"].(string)
		}
		if _, ok := headerValue["$ref"]; ok {
			header.ref = headerValue["$ref"].(string)
		}
		headers = append(headers, header)
	}

	return &headers, nil
}
