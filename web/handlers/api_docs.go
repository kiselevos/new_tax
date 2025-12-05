package handlers

type ApiDocsData struct {
	ApiVers   string         `json:"-"`
	Endpoints []EndpointInfo `json:"endpoints"`
}

type EndpointInfo struct {
	Group              string `json:"group"`
	Method             string `json:"method"`
	Path               string `json:"path"`
	Description        string `json:"description"`
	ExampleRequestLite string `json:"example_request_lite"`
	ExampleRequestFull string `json:"example_request_full"`
	ExampleResponse    string `json:"example_response"`
}
