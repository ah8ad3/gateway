package admin

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ah8ad3/gateway/pkg/proxy"
)

func getServices() []byte {
	var jData []byte
	jData, _ = json.Marshal(proxy.Services)

	return jData
}

// Welcome just an sample welcome
func Welcome(w http.ResponseWriter, r *http.Request) {
	_ = r

	//str, _ := proxy.AddPlugin("service1", "rateLimiter", nil)
	_, _ = w.Write([]byte("Welcome To Gateway"))
	return
}

// GETService get all services in admin mode
func GETService(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, _ = w.Write(getServices())
	return
}

// PostService to add some proxy
func PostService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var serveice proxy.Service

	err := decoder.Decode(&serveice)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "the template of json is incorrect"}`))
		return
	}

	if serveice.Name == "" || serveice.Path == "" || serveice.Server == nil || serveice.Version == 0 || serveice.Urls == nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "input form is incomplete"}`))
		return
	}

	for _, val := range proxy.Services {
		if val.Name == serveice.Name {
			if val.Version == serveice.Version {
				w.WriteHeader(400)
				_, _ = w.Write([]byte(`{"error": "name and version of proxy must be unique"}`))
				return
			}
		}
	}

	proxy.Services = append(proxy.Services, serveice)

	proxy.SaveServices()

	jData, _ := json.Marshal(serveice)
	w.WriteHeader(201)
	w.Write(jData)
	return
}

// DeleteService to delete proxy of list
func DeleteService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var serveice proxy.Service

	err := decoder.Decode(&serveice)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "the template of json is incorrect"}`))
		return
	}

	if serveice.Name == "" || serveice.Version == 0 {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "input form is incomplete"}`))
		return
	}

	for id, val := range proxy.Services {
		if val.Name == serveice.Name && val.Version == serveice.Version {
			proxy.Services = append(proxy.Services[:id], proxy.Services[id+1:]...)
			proxy.SaveServices()
			w.WriteHeader(204)
			_, _ = w.Write([]byte(`{"ok": "deleted"}`))
			return
		}
	}

	w.WriteHeader(404)
	_, _ = w.Write([]byte(`{"error": "record not found"}`))
	return

}

// UpdateService to update proxy staff
// this is put method to replace whole object with new one
func UpdateService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	var serveice proxy.Service

	err := decoder.Decode(&serveice)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "the template of json is incorrect"}`))
		return
	}

	if serveice.Name == "" || serveice.Path == "" || serveice.Server == nil || serveice.Version == 0 || serveice.Urls == nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "input form is incomplete"}`))
		return
	}

	for id, val := range proxy.Services {
		if val.Name == serveice.Name {
			if val.Version == serveice.Version {
				proxy.Services[id] = serveice
				proxy.SaveServices()

				jData, _ := json.Marshal(serveice)
				w.Write(jData)
				return
			}
		}
	}

	w.WriteHeader(400)
	_, _ = w.Write([]byte(`{"error": "name and version of proxy must be unique"}`))
	return
}

// GETServiceSlug to query with service_name
func GETServiceSlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := chi.URLParam(r, "service_name")
	if name == "" {
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "url paramter not found"}`))
		return
	}
	var services []proxy.Service

	for _, val := range proxy.Services {
		if val.Name == name {
			services = append(services, val)
		}
	}

	if len(services) == 0 {
		w.WriteHeader(404)
		w.Write([]byte(`{"error": "no any service found"}`))
		return
	}

	data, _ := json.Marshal(services)
	w.Write(data)
	return
}
