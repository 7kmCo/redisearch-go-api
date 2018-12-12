package handler

import (
	"net/http"
	"io"
	"fmt"
	"encoding/json"

	"github.com/7kmCo/redisearch-go-api/redisearch"
	"github.com/go-chi/chi"	
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"response"`
}

type handler struct {
	client *redisearch.Client
}

func New(client *redisearch.Client) chi.Router {
	h := handler{client}
	router := chi.NewRouter()
	router.Get("/info/{idx}", requestHandler(h.Info))
	// router.Post("/search", requestHandler(h.Search))

	// router.Post("/create", createHandler)
	// router.Post("/add", addHandler)
	// router.Delete("/drop/{idx}", dropHandler)
	// router.Route("/suggestion", func(r chi.Router) {
	// 	r.Post("/add", addSuggestion)
	// 	r.Post("/get", getSuggestion)
	// })
	http.ListenAndServe(":8080", router)
	return router
}

func requestHandler(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		err = json.NewEncoder(w).Encode(response{Data: data, Success: err == nil})
		if err != nil {
			fmt.Printf("could not encode response to output: %v", err)
		}
	}
}
