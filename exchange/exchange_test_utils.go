package exchange

import (
	"github.com/sharat87/httpbun/response"
	"github.com/sharat87/httpbun/server/spec"
	"github.com/sharat87/httpbun/util"
	"net/http"
	"net/url"
)

func InvokeHandlerForTest(path string, req http.Request, routePat string, fn HandlerFn) response.Response {
	if req.URL != nil {
		panic("req.URL must be nil")
	}
	req.URL, _ = url.Parse("http://localhost/" + path)

	ex := New(
		nil,
		&req,
		spec.Spec{
			PathPrefix: "",
		},
	)

	var isMatch bool
	ex.fields, isMatch = util.MatchRoutePat(routePat, ex.RoutedPath)
	if !isMatch {
		panic("Route pattern did not match path")
	}

	return fn(ex)
}
