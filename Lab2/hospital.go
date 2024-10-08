package main

import (
	"fmt"
)

type Doctor struct {
	Name   string
	Salary int
}

type Patient struct {
	Name   string
	Age    int
	Doctor Doctor
}

type Hospital struct {
	Name     string
	Location string
	Patients []Patient
	Doctors  []Doctor
}

func (d Doctor) display() string {
	return fmt.Sprintf("Doctor: %s, Salary: %d", d.Name, d.Salary)
}

func (p Patient) display() string {
	return fmt.Sprintf("Patient: %s, Age: %d, Assigned Doctor: %s", p.Name, p.Age, p.Doctor.display())
}

func (h Hospital) display() string {
	return fmt.Sprintf("Hospital: %s, Location: %s", h.Name, h.Location)
}