package handlers

import "net/http"

func RegisterHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return 
	}
}

