package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

var (
	doctors      []Doctor
	patients     []Patient
	doctorsMutex sync.Mutex
	patientsMutex sync.Mutex
)

func loadData() error {
	doctorsMutex.Lock()
	defer doctorsMutex.Unlock()
	patientsMutex.Lock()
	defer patientsMutex.Unlock()

	if _, err := os.Stat("doctors.json"); err == nil {
		data, err := ioutil.ReadFile("doctors.json")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &doctors); err != nil {
			return err
		}
	}

	if _, err := os.Stat("patients.json"); err == nil {
		data, err := ioutil.ReadFile("patients.json")
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &patients); err != nil {
			return err
		}
	}

	return nil
}

func saveDoctors() error {
	doctorsMutex.Lock()
	defer doctorsMutex.Unlock()

	data, err := json.MarshalIndent(doctors, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("doctors.json", data, 0644)
}

func savePatients() error {
	patientsMutex.Lock()
	defer patientsMutex.Unlock()

	data, err := json.MarshalIndent(patients, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("patients.json", data, 0644)
}
