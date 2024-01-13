package model

import (
	"fmt"
	"os"
	"strings"

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
	response     string      //FIXME:
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
	ref         string
	description string
}

type ExternalDoc struct {
	Description string `yaml:"description"`
	Url         string `yaml:"url"`
}

// --------------

type OpenApi struct { //FIXME:
	openapi     string
	title       string
	description string
	version     string
	servers     []Server
	tags        []Tag
	paths       []Path
}

// type response struct {
// 	refType     string
// 	ref         string
// 	description string
// }

func NewOpenApi(filename string) (*OpenApi, error) {
	// INFO: 読込み
	file, err := os.ReadFile("./openapi/orders/openapi.yaml")
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %s", err.Error())
	}

	// INFO: baseItem
	var apiBase ApiBase
	err = yaml.Unmarshal([]byte(file), &apiBase)
	if err != nil {
		return nil, err
	}
	// fmt.Println(apiBase) //WARNING:

	// INFO: PathItem
	yamlPath, err := yamlChild([]byte(file), "paths")
	if err != nil {
		return nil, err
	}
	paths, err := pathItems(yamlPath)
	if err != nil {
		return nil, err
	}
	fmt.Println(len(*paths)) //WARNING:

	// INFO: パース処理
	var mapData map[string]interface{}
	openapi, err := baseInfo(mapData)
	if err != nil {
		return nil, err
	}
	// paths, err = pathItemD(mapData)
	// if err != nil {
	// 	return nil, err
	// }
	// openapi.paths = *paths

	return openapi, nil
}

func (api *OpenApi) ToStr() {
	fmt.Println(api)
}

func (api *OpenApi) Openapi() string {
	return api.openapi
}

func (api *OpenApi) Info() (string, string, string) {
	return api.title, api.description, api.version
}

func (api *OpenApi) Servers() []Server {
	return api.servers
}

func (api *OpenApi) Tags() []Tag {
	return api.tags
}

func (api *OpenApi) Paths() []Path {
	return api.paths
}

// ----+----+----+----+----+----+----+----+----+----+----+

func (path *Path) ToArray() []string {
	return []string{
		strings.Join(path.Tags, ","),
		path.url,
		path.method,
		path.OperationId,
		path.Summary,
		path.Description,
		fmt.Sprint(10),
		"path.requestBody.description",
		path.response,
		"lo.Ternary(path.hasExternalDocs, ",
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

			requestBody, err := requestBody(methpdValue)
			if err != nil {
				return nil, err
			}
			path.requestBody = *requestBody

			path.url = url
			path.method = method

			// fmt.Println(methodMap) //WARNING:
			fmt.Println(path.requestBody) //WARNING:
			fmt.Println("***** *****")    //WARNING:

			paths = append(paths, *path)
		}
	}
	return &paths, nil
}

// ----+----+----+----+----+----+----+----+----+----+----+

func requestBody(methodValue map[string]interface{}) (*RequestBody, error) {
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

func baseInfo(yamldata map[string]interface{}) (*OpenApi, error) {

	return &OpenApi{
		openapi:     "openapi",
		title:       "title",
		description: "description",
		version:     "version",
		servers:     []Server{},
		tags:        []Tag{},
	}, nil
}

// func pathItemD(yamldata map[string]interface{}) (*[]Path, error) {

// 	// INFO: path階層
// 	for path, pathItem := range tPaths {

// 		// INFO: httpMethod階層
// 		for method, methodItem := range pathItem {

// 			// TODO: requestBody
// 			// requestBody := RequestBody{}
// 			// if tRequestBody, ok := methodItem["requestBody"]; ok { // requestBodyセクションがない場合は回避
// 			// 	tRequestBody, ok := tRequestBody.(map[interface{}]interface{})
// 			// 	if !ok {
// 			// 		return nil, errors.New("cannot parse yaml data [paths][method][requestBody]")
// 			// 	}

// 			// 	tDescription, ok := tRequestBody["description"]
// 			// 	if !ok {
// 			// 		return nil, errors.New("cannot parse yaml data [paths][method][requestBody][description]")
// 			// 	}
// 			// 	description, ok = tDescription.(string)
// 			// 	if !ok {
// 			// 		return nil, errors.New("cannot parse yaml data [paths][method][requestBody][description]")
// 			// 	}
// 			// 	requestBody.description = description

// 			// 	if ref, ok := tRequestBody["$ref"]; ok {
// 			// 		ref, ok := ref.(string)
// 			// 		if !ok {
// 			// 			return nil, errors.New("cannot parse yaml data [paths][method][requestBody][$ref]")
// 			// 		}
// 			// 		requestBody.refType = "requestBodies"
// 			// 		requestBody.ref = ref

// 			// 	} else if contentV, ok := tRequestBody["content"]; ok {
// 			// 		contentV, ok := contentV.(map[interface{}]interface{})
// 			// 		if !ok {
// 			// 			return nil, errors.New("cannot parse yaml data [paths][method][requestBody][content]")
// 			// 		}
// 			// 		appV, ok := contentV["application/json"]
// 			// 		if !ok {
// 			// 			return nil, errors.New("cannot parse yaml data [paths][method][requestBody][content]")
// 			// 		}
// 			// 		appVC, ok := appV.(map[interface{}]interface{})
// 			// 		if !ok {
// 			// 			return nil, errors.New("cannot parse yaml data [paths][method][requestBody][content]")
// 			// 		}
// 			// 		schemaV, ok := appVC["schema"]
// 			// 		if !ok {
// 			// 			return nil, errors.New("cannot parse yaml data [paths][method][requestBody][content]")
// 			// 		}
// 			// 		schemaVC, ok := schemaV.(map[interface{}]interface{})
// 			// 		if !ok {
// 			// 			return nil, errors.New("cannot parse yaml data [paths][method][requestBody][content]")
// 			// 		}
// 			// 		ref, ok := schemaVC["$ref"]
// 			// 		if ok {
// 			// 			ref, ok := ref.(string)
// 			// 			if !ok {
// 			// 				return nil, errors.New("cannot parse yaml data [paths][method][requestBody][content]")
// 			// 			}
// 			// 			requestBody.refType = "schemas"
// 			// 			requestBody.ref = ref
// 			// 		} else {
// 			// 			requestBody.refType = "none"
// 			// 		}

// 			// 	} else {
// 			// 		return nil, errors.New("cannot parse yaml data [paths][method][requestBody]")
// 			// 	}

// 			// 	fmt.Println(tRequestBody)
// 			// 	fmt.Println(requestBody)

// 			// }

// 			// FIXME: responses

// 			// INFO: externalDocs
// 			_, hasExternalDocs := methodItem["externalDocs"]

// 			paths = append(paths, Path{
// 				url:             path,
// 				method:          method,
// 				OperationId:     operationId,
// 				Summary:         summary,
// 				Description:     description,
// 				tags:            tags,
// 				parameters:      parameters,
// 				requestBody:     "requestBody",
// 				response:        "bbb",
// 				hasExternalDocs: hasExternalDocs,
// 			})
// 			// for k := range methodItem {
// 			// 	fmt.Println(k)
// 			// }
// 			// fmt.Println("***** *****")
// 		}
// 	}
// 	return &paths, nil
// }

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
