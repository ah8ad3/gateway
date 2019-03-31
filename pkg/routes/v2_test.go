package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ah8ad3/gateway/pkg/proxy"
)

func testAPIReverse(t *testing.T, query bool, path string, pxy Proxy) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Host = "localhost:8000"
	if query {
		req.URL.RawQuery = "test=hi"
		req.Host = "localhost:8000/?hi=there"
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(pxy.handleProxy)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
}

func TestNewProxy(t *testing.T) {
	res := singleJoiningSlash("a", "b")
	if res != "a/b" {
		t.Fatal("error in singleJoin method")
	}
	res = singleJoiningSlash("a", "/b/ac/")
	if res != "a/b/ac/" {
		t.Fatal("error in singleJoin method")
	}
	res = singleJoiningSlash("/", "/b/ac/")
	if res != "/b/ac/" {
		t.Fatal("error in singleJoin method")
	}

	getHost(proxy.Services[0])
	getHost(proxy.Services[2])

	_proxy := NewProxy(proxy.Services[2])
	testAPIReverse(t, false, "/google/foo", _proxy)
	testAPIReverse(t, true, "/google/", _proxy)

}
