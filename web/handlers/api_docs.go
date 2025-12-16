package handlers

type ApiDocsData struct {
	ApiVers   string         `json:"-"`
	Endpoints []EndpointInfo `json:"endpoints"`
}

type EndpointInfo struct {
	Group           string      `json:"group"`
	Method          string      `json:"method"`
	Path            string      `json:"path"`
	Description     string      `json:"description"`
	ExampleRequest  interface{} `json:"example_request"`
	ExampleResponse interface{} `json:"example_response"`
}
