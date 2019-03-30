package routes

import (
	"github.com/ah8ad3/gateway/pkg/proxy"
	"testing"
)

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

}
