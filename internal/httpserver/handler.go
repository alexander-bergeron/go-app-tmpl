package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	// jsonResp, err := json.Marshal(s.db.Health())
	// if err != nil {
	// 	log.Printf("Error marshaling JSON: %v", err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// _, _ = w.Write(jsonResp)
}

// Example Handlers
func (s *Server) exampleGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET example response"))
}

func (s *Server) examplePostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST example response"))
}

func (s *Server) exampleDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DELETE example response"))
}
