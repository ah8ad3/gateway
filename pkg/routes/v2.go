package routes

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy to create path prefix routes
type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	prefix string
}

// NewProxy to create Proxy instance easier
func NewProxy(target string, prefix string) *Proxy {
	url, _ := url.Parse(target)
	fmt.Println(url)

	return &Proxy{target: url, proxy: singleHodtRewriteReverse(url, prefix)}
}

func (p *Proxy) handleProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	// p.proxy.Transport = &myTransport{}
	// also here can define all middleware stuff
	p.proxy.ServeHTTP(w, r)
}

func singleHodtRewriteReverse(target *url.URL, prefix string) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		if prefix == "/" {
			req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		} else {

			url := strings.Replace(req.URL.Path, prefix, "", 1)
			req.URL.Path = singleJoiningSlash(target.Path, url)
		}
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
