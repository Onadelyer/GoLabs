package main

import (
	"fmt"
	"net/http"
)

func main() {
	if err := loadData(); err != nil {
		fmt.Println("Error loading data:", err)
		return
	}

	http.HandleFunc("/doctors", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getDoctorsHandler(w, r)
		case http.MethodPost:
			createDoctorHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/doctors/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getDoctorHandler(w, r)
		case http.MethodPut:
			updateDoctorHandler(w, r)
		case http.MethodDelete:
			deleteDoctorHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/patients", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPatientsHandler(w, r)
		case http.MethodPost:
			createPatientHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/patients/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPatientHandler(w, r)
		case http.MethodPut:
			updatePatientHandler(w, r)
		case http.MethodDelete:
			deletePatientHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
