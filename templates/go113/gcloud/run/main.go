package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type runRequest struct {
	// Add request fields
}

type runResponse struct {
	// Add response fields
	Response string `json:"response"`
}

func init() {
	// Do any required initialisation
}

func validateRequest(req runRequest) error {
	// switch {
	// case req.RequiredField == "":
	// 	return errors.New("required_field is missing")
	// }
	return nil
}

// handler handles the request
func handler(w http.ResponseWriter, r *http.Request) {
	// Read and validate the request
	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validateRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add your code

	// Format and write the result
	rsp, err := json.Marshal(runResponse{
		// Set response fields
		Response: "Hello, world!",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp)
}

func main() {
	log.Print("starting server...")
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
