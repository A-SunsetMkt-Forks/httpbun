package redirect

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sharat87/httpbun/ex"
	"github.com/sharat87/httpbun/response"
)

const MaxRedirectCount = 20

var RouteList = []ex.Route{
	ex.NewRoute(`/redirect(-to)?/?`, handleRedirectTo),
	ex.NewRoute(`/(?P<mode>relative-|absolute-)?redirect/(?P<count>\d+)`, handleRedirectCount),
}

func handleRedirectTo(ex *ex.Exchange) response.Response {
	query := ex.Request.URL.Query()
	urls := query["url"]
	if len(urls) < 1 || urls[0] == "" {
		return response.BadRequest("Need url parameter")
	}

	statusCodes := query["status_code"]
	if statusCodes == nil {
		statusCodes = query["status"]
	}

	statusCode := http.StatusFound
	if statusCodes != nil {
		var err error
		if statusCode, err = strconv.Atoi(statusCodes[0]); err != nil {
			return response.BadRequest("status_code must be an integer")
		}
		if statusCode < 300 || statusCode > 399 {
			statusCode = 302
		}
	}

	return response.Response{
		Status: statusCode,
		Header: http.Header{
			"Location": {urls[0]},
		},
	}
}

func handleRedirectCount(ex *ex.Exchange) response.Response {
	isAbsolute := ex.Field("mode") == "absolute-"
	n, _ := strconv.Atoi(ex.Field("count"))

	if n < 0 {
		return response.BadRequest("count must be a non-negative integer")

	} else if n > MaxRedirectCount {
		return response.BadRequest("count cannot be greater than %v", MaxRedirectCount)

	} else if n > 1 {
		target := fmt.Sprint(n - 1)
		if isAbsolute {
			target = "/absolute-redirect/" + target
		}
		return *ex.RedirectResponse(target)

	} else {
		var target string
		if isAbsolute {
			target = "/anything"
		} else {
			target = "../anything"
		}
		return *ex.RedirectResponse(target)

	}
}
