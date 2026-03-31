package main

import (
	"fmt"
	"net/http"
)

func healthzHandler(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(fmt.Sprintf(`
<html>
<body>
    <h1>Welcome, Beste Zimmer Admin</h1>
    <p>The home page has been visited %d times!</p>
  </body>
</html>`, cfg.server_hits.Load())))
}
