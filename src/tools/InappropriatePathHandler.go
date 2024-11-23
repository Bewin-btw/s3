package tools

import "net/http"

func BadPathHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "There is no such endpoint", http.StatusBadRequest)
}
