package operatorai

import (
	"encoding/json"
	"net/http"
)

type functionRequest struct {
	// Add request fields
}

type functionResponse struct {
	// Add response fields
	Response string `json:"response"`
}

func init() {
	// Do any required initialisation
}

func validateRequest(req functionRequest) error {
	// switch {
	// case req.RequiredField == "":
	// 	return errors.New("api_key is missing")
	// }
	return nil
}

// {{.FunctionName}} add docstring here
func {{.FunctionName}}(w http.ResponseWriter, r *http.Request) {
	// Read and validate the request
	var req functionRequest
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
	rsp, err := json.Marshal(functionResponse{
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

