package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/lubouski/go-marathon/receipt-processor/utils"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("GET /users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		w.Write([]byte("User ID: " + userID))
	})

	router.HandleFunc("GET /receipts/{id}/points", s.getReceipt)
	router.HandleFunc("POST /receipts/process", s.postReceipt)

	middlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		//RequireAuthMiddleware,
	)

	server := http.Server{
		Addr: s.addr,
		//Handler: RequireAuthMiddleware(RequestLoggerMiddleware(router)),
		Handler: middlewareChain(router),
	}

	log.Printf("Server has started %s", s.addr)

	return server.ListenAndServe()
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "Bearer token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next.ServeHTTP
	}
}

type Receipt struct {
	Receipts string `json:"receipts"`
	Total    float64 `json:"total,string"`
}

func (s *APIServer) postReceipt(w http.ResponseWriter, r *http.Request) {
	var res Receipt
	uuid := utils.GenerateUUID()
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	receiptDB[uuid] = res
	fmt.Fprintf(w, "Receipt: %+v, id: %s\n", res, uuid)
//	w.Write([]byte("User ID: " + userID))
}


func (s *APIServer) getReceipt(w http.ResponseWriter, r *http.Request) {
	var res Receipt
	res = receiptDB[r.PathValue("id")]
	fmt.Fprintf(w, "Receipt: %+v", res)
}
