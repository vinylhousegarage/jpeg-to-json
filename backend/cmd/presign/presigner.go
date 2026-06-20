package main

import (
	"encoding/json"
	"net/http"
	
	"github.com/vinylhousegarage/jpeg-to-json/backend/apierror"
)

func enableCORS(w http.ResponseWriter, r *http.Request, origin string) {
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

type PresignRequest struct {
	Filename string `json:"filename"`
	FileType string `json:"filetype"`
}

func validateRequest(r *http.Request) (*PresignRequest, error) {
	if r.Method != http.MethodPost {
		return nil, apierror.New(apierror.ErrorCodeInvalidMethod, http.StatusMethodNotAllowed, nil)
	}

	var req PresignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, apierror.New(apierror.ErrorCodeInvalidJSON, http.StatusBadRequest, err)
	}

	if req.Filename == "" {
		return nil, apierror.New(apierror.ErrorCodeMissingFilename, http.StatusBadRequest, nil)
	}

	return &req, nil
}
