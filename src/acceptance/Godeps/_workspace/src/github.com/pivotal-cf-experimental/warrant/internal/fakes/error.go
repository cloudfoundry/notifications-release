package fakes

import (
	"fmt"
	"net/http"
)

func (s *UAAServer) Error(w http.ResponseWriter, status int, message, errorType string) {
	output := fmt.Sprintf(`{"message":"%s","error":"%s"}`, message, errorType)

	w.WriteHeader(status)
	w.Write([]byte(output))
}
