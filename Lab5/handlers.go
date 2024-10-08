package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func parseIDFromURL(path string) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	idStr := parts[len(parts)-1]
	return strconv.Atoi(idStr)
}

func getDoctorsHandler(w http.ResponseWriter, r *http.Request) {
	doctorsMutex.Lock()
	defer doctorsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doctors)
}

func getDoctorHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
		return
	}

	doctorsMutex.Lock()
	defer doctorsMutex.Unlock()

	for _, doctor := range doctors {
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &doctor); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	doctorsMutex.Lock()
	defer doctorsMutex.Unlock()

	if len(doctors) > 0 {
		doctor.ID = doctors[len(doctors)-1].ID + 1
	} else {
		doctor.ID = 1
	}
	doctors = append(doctors, doctor)
	saveDoctors()

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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &updatedDoctor); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedDoctor.ID = id

	doctorsMutex.Lock()
	defer doctorsMutex.Unlock()

	for i, doctor := range doctors {
		if doctor.ID == id {
			doctors[i] = updatedDoctor
			saveDoctors()
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

	doctorsMutex.Lock()
	defer doctorsMutex.Unlock()

	for i, doctor := range doctors {
		if doctor.ID == id {
			doctors = append(doctors[:i], doctors[i+1:]...)
			saveDoctors()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Doctor not found", http.StatusNotFound)
}

func getPatientsHandler(w http.ResponseWriter, r *http.Request) {
	patientsMutex.Lock()
	defer patientsMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

func getPatientHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromURL(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid patient ID", http.StatusBadRequest)
		return
	}

	patientsMutex.Lock()
	defer patientsMutex.Unlock()

	for _, patient := range patients {
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &patient); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	patientsMutex.Lock()
	defer patientsMutex.Unlock()

	if len(patients) > 0 {
		patient.ID = patients[len(patients)-1].ID + 1
	} else {
		patient.ID = 1
	}
	patients = append(patients, patient)
	savePatients()

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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &updatedPatient); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedPatient.ID = id

	patientsMutex.Lock()
	defer patientsMutex.Unlock()

	for i, patient := range patients {
		if patient.ID == id {
			patients[i] = updatedPatient
			savePatients()
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

	patientsMutex.Lock()
	defer patientsMutex.Unlock()

	for i, patient := range patients {
		if patient.ID == id {
			patients = append(patients[:i], patients[i+1:]...)
			savePatients()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Patient not found", http.StatusNotFound)
}
