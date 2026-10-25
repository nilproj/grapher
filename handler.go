package grapher

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/graph-gophers/graphql-go"
)

type Handler struct {
	Schema      *graphql.Schema
	Explorer    Explorer
	BaseContext func(req *http.Request) context.Context
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGET(w, r)
	case http.MethodPost:
		h.handlePOST(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGET(w http.ResponseWriter, r *http.Request) {
	if e := h.Explorer; e != nil {
		e.ServeHTTP(w, r)
		return
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func (h *Handler) handlePOST(w http.ResponseWriter, r *http.Request) {
	var params graphqlParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	if h.BaseContext != nil {
		ctx = h.BaseContext(r)
	}
	result := h.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type graphqlParams struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName"`
	Variables     map[string]any `json:"variables"`
}
