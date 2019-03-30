package routes

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/ah8ad3/gateway/pkg/proxy"
)

// Proxy to create path prefix routes
type Proxy struct {
	proxy   httputil.ReverseProxy
	service proxy.Service
}

// NewProxy to create Proxy instance easier
func NewProxy(service proxy.Service) Proxy {
	return Proxy{service: service, proxy: singleHostRewriteReverse(service)}
}

func (p Proxy) handleProxy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")
	// p.proxy.Transport = &myTransport{}
	// also here can define all middleware stuff
	p.proxy.ServeHTTP(w, r)
}

func getHost(proxy proxy.Service) *url.URL {
	rand.Seed(time.Now().Unix())
	if proxy.UPHostsCount == 0 {
		host, _ := url.Parse(proxy.Server[rand.Intn(len(proxy.Server))].Server)
		return host
	}
	host, _ := url.Parse(proxy.UPHosts[rand.Intn(len(proxy.UPHosts))])
	return host

}

func singleHostRewriteReverse(proxy proxy.Service) httputil.ReverseProxy {

	director := func(req *http.Request) {
		target := getHost(proxy)
		targetQuery := target.RawQuery

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		if proxy.Path == "/" {
			req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		} else {

			url := strings.Replace(req.URL.Path, proxy.Path, "", 1)
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
	return httputil.ReverseProxy{Director: director}
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
