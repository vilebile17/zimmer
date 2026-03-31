package main

import (
	"fmt"
	"net/http"
)

func handler_get_endpoint_path(url string) http.FileSystem {
	return http.Dir("./app" + url)
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(handler_get_endpoint_path("")))
	mux.HandleFunc("/healthz", healthzHandler)

	fmt.Printf("Hosting Beste Zimmer at http://localhost:%s\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
