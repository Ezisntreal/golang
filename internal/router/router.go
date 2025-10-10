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

	http.HandleFunc("/metric", func(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/sensor", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetSensorsHandler(w, r)
		case http.MethodPost:
			handler.CreateSensorHandler(w, r)
		case http.MethodPut:
			handler.UpdateSensorHandler(w, r)
		case http.MethodDelete:
			handler.DeleteSensorHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/box", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetBoxesHandler(w, r)
		case http.MethodPost:
			handler.CreateBoxHandler(w, r)
		case http.MethodPut:
			handler.UpdateBoxHandler(w, r)
		case http.MethodDelete:
			handler.DeleteBoxHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
