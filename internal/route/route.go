package route

import "net/http"

type HTTPVerb string

const (
	GET     HTTPVerb = "GET"
	POST    HTTPVerb = "POST"
	PUT     HTTPVerb = "PUT"
	PATCH   HTTPVerb = "PATCH"
	DELETE  HTTPVerb = "DELETE"
	HEAD    HTTPVerb = "HEAD"
	OPTIONS HTTPVerb = "OPTIONS"
)

type Route struct {
	name     string
	httpVerb HTTPVerb
	handler  func(responseWriter http.ResponseWriter, request *http.Request)
}

func NewRoute(name string, httpVerb HTTPVerb, handler func(responseWriter http.ResponseWriter, request *http.Request)) Route {
	return Route{
		name:     name,
		httpVerb: httpVerb,
		handler:  handler,
	}
}

func (r *Route) Name() string {
	return r.name
}

func (r *Route) Handler() func(responseWriter http.ResponseWriter, request *http.Request) {
	return r.handler
}
