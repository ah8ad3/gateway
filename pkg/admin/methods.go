package admin

import (
	"net/http"

	"github.com/ah8ad3/gateway/pkg/proxy"
)

// Welcome just an sample welcome
func Welcome(w http.ResponseWriter, r *http.Request) {
	_ = r

	str, _ := proxy.AddPlugin("service1", "ipBlocker", nil)
	_, _ = w.Write([]byte(str))
}
