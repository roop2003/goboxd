package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thesouldev/goboxd/internal/languages"
	"github.com/thesouldev/goboxd/internal/runs"
)

func RunHandler(w http.ResponseWriter, r *http.Request) {
	var req runs.Request
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, runs.ErrorResponse{
			Error: runs.Error{Code: "bad_json", Message: "request body must be valid JSON"},
		})
		return
	}

	registry, err := languages.LoadDefaultRegistry()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, runs.ErrorResponse{
			Error: runs.Error{Code: "language_config_error", Message: "language registry could not be loaded"},
		})
		return
	}

	if apiErr := runs.ValidateRequest(req, registry); apiErr != nil {
		writeJSON(w, http.StatusBadRequest, runs.ErrorResponse{Error: *apiErr})
		return
	}

	executor := runs.Executor{Registry: registry}
	response, err := executor.Execute(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, runs.ErrorResponse{
			Error: runs.Error{
				Code:    "execution_error",
				Message: "code execution failed before a user-code result was available",
			},
		})
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
