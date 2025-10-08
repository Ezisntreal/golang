package router

import (
	"api/internal/handler"
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetUsersHandler(w, r)
		case http.MethodPost:
			handler.CreateUserHandler(w, r)
		case http.MethodPut:
			handler.UpdateUserHandler(w, r)
		case http.MethodDelete:
			handler.DeleteUserHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetMetricsHandler(w, r)
		case http.MethodPost:
			handler.CreateMetricHandler(w, r)
		case http.MethodPut:
			handler.UpdateMetricHandler(w, r)
		case http.MethodDelete:
			handler.DeleteMetricHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
