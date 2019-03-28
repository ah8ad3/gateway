package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testAPIGet(t *testing.T, method string, url string)  {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProxyHttp)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
}

func testAPIPost(t *testing.T, method string, path string)  {
	// Create a request to pass to our handler.
	payload := []byte(`{"name": "test product", "price": 11}`)


	req, err := http.NewRequest(method, path, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err.Error())
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postProxyHttp)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
}

func testAPIPut(t *testing.T, method string, path string) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(putProxyHttp)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
}

func testAPIDelete(t *testing.T, method string, path string) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteProxyHttp)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
}

func testAPIIntegrate(t *testing.T, method string, path string) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Host = "localhost:8000"

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(integrateProxyHttp)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
}


func TestGetService(t *testing.T) {
	GetService("google.com", "/about", "")
	GetService("google.com", "/about", "?q=there")
	GetService("localhost:8000", "/about", "?q=there")
	findService("/foo")
	findService("foo")

	testAPIGet(t, "GET", "/hi/asj")
	testAPIGet(t, "GET", "/hi")

	testAPIIntegrate(t, "GET", "/agg/")
	testAPIIntegrate(t, "GET", "/agg")
}

func TestPostService(t *testing.T) {
	PostService("localhost:8000", "/test", nil)
	PostService("reqres.in", "/api/users", nil)

	testAPIPost(t, "POST", "/as")
	testAPIPost(t, "POST", "/as/as")

	testAPIPut(t, "PUT", "/as")
	testAPIDelete(t, "PUT", "/as")
}
