package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Models
type Doctor struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Salary int    `json:"salary"`
}

type Patient struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	DoctorID int    `json:"doctor_id"`
}

// Storage with Mutex for concurrency safety
type Storage struct {
	doctors       []Doctor
	patients      []Patient
	doctorsMutex  sync.RWMutex
	patientsMutex sync.RWMutex
}

var store = Storage{}

// Logger
var (
	logger     *log.Logger
	logFile    *os.File
	loggerOnce sync.Once
)

// Configuration
const (
	AuthorizationHeader = "Authorization"
	AuthKey             =  "asdasdasd"
	LogFilePath         = "server.log"
)

// Initialize Logger
func initLogger() {
	var err error
	logFile, err = os.OpenFile(LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

// Utility Functions
func parseIDFromURL(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid path")
	}
	idStr := parts[len(parts)-1]
	return strconv.Atoi(idStr)
}

func (s *Storage) loadData() error {
	// Load Doctors
	if _, err := os.Stat("doctors.json"); err == nil {
		data, err := os.ReadFile("doctors.json")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &s.doctors); err != nil {
			return err
		}
	}

	// Load Patients
	if _, err := os.Stat("patients.json"); err == nil {
		data, err := os.ReadFile("patients.json")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &s.patients); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) saveDoctors() error {
	data, err := json.MarshalIndent(s.doctors, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("doctors.json", data, 0644)
}

func (s *Storage) savePatients() error {
	data, err := json.MarshalIndent(s.patients, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("patients.json", data, 0644)
}

// Middleware Functions

// LoggingMiddleware logs each incoming request to a log file.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure logger is initialized once
		loggerOnce.Do(initLogger)

		startTime := time.Now()
		// Proceed with the next handler
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)

		// Log format: [Timestamp] Method URL RemoteAddr Duration
		logger.Printf("%s %s %s %v\n", r.Method, r.URL.String(), r.RemoteAddr, duration)
	})
}

// AuthorizationMiddleware checks for a specific key in the request headers.
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authorization key from headers
		key := r.Header.Get(AuthorizationHeader)
		if key != AuthKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// Middleware Wrapper Helper
func applyMiddlewares(handler http.Handler) http.Handler {
	handler = LoggingMiddleware(handler)
	handler = AuthorizationMiddleware(handler)
	return handler
}

// Handlers for Doctors

// getDoctorsHandler retrieves all doctors, with optional filtering.
func getDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	store.doctorsMutex.RLock()
	defer store.doctorsMutex.RUnlock()

	// Parse query parameters for filtering
	nameFilter := r.URL.Query().Get("name")
	salaryFilter := r.URL.Query().Get("salary")

	filteredDoctors := filterDoctors(store.doctors, nameFilter, salaryFilter)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredDoctors)
}

// filterDoctors filters the list of doctors based on name and salary.
func filterDoctors(doctors []Doctor, name string, salary string) []Doctor {
	var filtered []Doctor
	for _, doctor := range doctors {
		matches := true
		if name != "" && !strings.Contains(strings.ToLower(doctor.Name), strings.ToLower(name)) {
			matches = false
		}
		if salary != "" {
			sal, err := strconv.Atoi(salary)
			if err != nil || doctor.Salary != sal {
				matches = false
			}
		}
		if matches {
			filtered = append(filtered, doctor)
		}
	}
	return filtered
}

func getDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	store.doctorsMutex.RLock()
	defer store.doctorsMutex.RUnlock()

	for _, doctor := range store.doctors {
		if doctor.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(doctor)
			return
		}
	}

	http.Error(w, "Doctor not found", http.StatusNotFound)
}

func createDoctorHandler(w http.ResponseWriter, r *http.Request) {
	var doctor Doctor
	if err := json.NewDecoder(r.Body).Decode(&doctor); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	store.doctorsMutex.Lock()
	defer store.doctorsMutex.Unlock()

	// Assign ID
	if len(store.doctors) > 0 {
		doctor.ID = store.doctors[len(store.doctors)-1].ID + 1
	} else {
		doctor.ID = 1
	}
	store.doctors = append(store.doctors, doctor)

	if err := store.saveDoctors(); err != nil {
		http.Error(w, "Failed to save doctor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctor)
}

func updateDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	var updatedDoctor Doctor
	if err := json.NewDecoder(r.Body).Decode(&updatedDoctor); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedDoctor.ID = id

	store.doctorsMutex.Lock()
	defer store.doctorsMutex.Unlock()

	for i, doctor := range store.doctors {
		if doctor.ID == id {
			store.doctors[i] = updatedDoctor
			if err := store.saveDoctors(); err != nil {
				http.Error(w, "Failed to save doctor", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedDoctor)
			return
		}
	}

	http.Error(w, "Doctor not found", http.StatusNotFound)
}

func deleteDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	store.doctorsMutex.Lock()
	defer store.doctorsMutex.Unlock()

	for i, doctor := range store.doctors {
		if doctor.ID == id {
			store.doctors = append(store.doctors[:i], store.doctors[i+1:]...)
			if err := store.saveDoctors(); err != nil {
				http.Error(w, "Failed to delete doctor", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Doctor not found", http.StatusNotFound)
}

// Handlers for Patients

// getPatientsHandler retrieves all patients, with optional filtering.
func getPatientsHandler(w http.ResponseWriter, r *http.Request) {
	store.patientsMutex.RLock()
	defer store.patientsMutex.RUnlock()

	// Parse query parameters for filtering
	nameFilter := r.URL.Query().Get("name")
	ageFilter := r.URL.Query().Get("age")
	doctorIDFilter := r.URL.Query().Get("doctor_id")

	filteredPatients := filterPatients(store.patients, nameFilter, ageFilter, doctorIDFilter)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredPatients)
}

// filterPatients filters the list of patients based on name, age, and doctor_id.
func filterPatients(patients []Patient, name string, age string, doctorID string) []Patient {
	var filtered []Patient
	for _, patient := range patients {
		matches := true
		if name != "" && !strings.Contains(strings.ToLower(patient.Name), strings.ToLower(name)) {
			matches = false
		}
		if age != "" {
			ageInt, err := strconv.Atoi(age)
			if err != nil || patient.Age != ageInt {
				matches = false
			}
		}
		if doctorID != "" {
			docID, err := strconv.Atoi(doctorID)
			if err != nil || patient.DoctorID != docID {
				matches = false
			}
		}
		if matches {
			filtered = append(filtered, patient)
		}
	}
	return filtered
}

func getPatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	store.patientsMutex.RLock()
	defer store.patientsMutex.RUnlock()

	for _, patient := range store.patients {
		if patient.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(patient)
			return
		}
	}

	http.Error(w, "Patient not found", http.StatusNotFound)
}

func createPatientHandler(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	store.patientsMutex.Lock()
	defer store.patientsMutex.Unlock()

	// Assign ID
	if len(store.patients) > 0 {
		patient.ID = store.patients[len(store.patients)-1].ID + 1
	} else {
		patient.ID = 1
	}
	store.patients = append(store.patients, patient)

	if err := store.savePatients(); err != nil {
		http.Error(w, "Failed to save patient", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

func updatePatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	var updatedPatient Patient
	if err := json.NewDecoder(r.Body).Decode(&updatedPatient); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedPatient.ID = id

	store.patientsMutex.Lock()
	defer store.patientsMutex.Unlock()

	for i, patient := range store.patients {
		if patient.ID == id {
			store.patients[i] = updatedPatient
			if err := store.savePatients(); err != nil {
				http.Error(w, "Failed to save patient", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedPatient)
			return
		}
	}

	http.Error(w, "Patient not found", http.StatusNotFound)
}

func deletePatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	store.patientsMutex.Lock()
	defer store.patientsMutex.Unlock()

	for i, patient := range store.patients {
		if patient.ID == id {
			store.patients = append(store.patients[:i], store.patients[i+1:]...)
			if err := store.savePatients(); err != nil {
				http.Error(w, "Failed to delete patient", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Patient not found", http.StatusNotFound)
}

// Main Function
func main() {
	// Initialize Logger
	loggerOnce.Do(initLogger)
	defer logFile.Close()

	// Load existing data
	if err := store.loadData(); err != nil {
		fmt.Println("Error loading data:", err)
		return
	}

	// Doctor Routes
	doctorsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getDoctorsHandler(w, r)
		case http.MethodPost:
			createDoctorHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.Handle("/doctors", applyMiddlewares(doctorsHandler))

	doctorIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	http.Handle("/doctors/", applyMiddlewares(doctorIDHandler))

	// Patient Routes
	patientsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPatientsHandler(w, r)
		case http.MethodPost:
			createPatientHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.Handle("/patients", applyMiddlewares(patientsHandler))

	patientIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	http.Handle("/patients/", applyMiddlewares(patientIDHandler))

	fmt.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}