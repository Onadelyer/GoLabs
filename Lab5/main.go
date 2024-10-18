package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "strings"
    "sync"
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
    doctors      []Doctor
    patients     []Patient
    doctorsMutex sync.RWMutex
    patientsMutex sync.RWMutex
}

var store = Storage{}

// Utility Functions
func parseIDFromURL(path string) (int, error) {
    parts := strings.Split(strings.Trim(path, "/"), "/")
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

// Handlers for Doctors
func getDoctorsHandler(w http.ResponseWriter, r *http.Request) {
    store.doctorsMutex.RLock()
    defer store.doctorsMutex.RUnlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(store.doctors)
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
func getPatientsHandler(w http.ResponseWriter, r *http.Request) {
    store.patientsMutex.RLock()
    defer store.patientsMutex.RUnlock()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(store.patients)
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
    // Load existing data
    if err := store.loadData(); err != nil {
        fmt.Println("Error loading data:", err)
        return
    }

    // Doctor Routes
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

    // Patient Routes
    http.HandleFunc("/patients", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            fmt.Println("GET /patients")
            getPatientsHandler(w, r)
        case http.MethodPost:
            fmt.Println("POST /patients")
            createPatientHandler(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    http.HandleFunc("/patients/", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            fmt.Println("GET /patients/{id}")
            getPatientHandler(w, r)
        case http.MethodPut:
            fmt.Println("PUT /patients/{id}")
            updatePatientHandler(w, r)
        case http.MethodDelete:
            deletePatientHandler(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })

    fmt.Println("Server is running on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Server failed:", err)
    }
}