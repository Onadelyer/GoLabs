package main

import (
	"fmt"
)

func main() {
	doctors := []Doctor{
		{Name: "Dr. Smith", Salary: 5000},
		{Name: "Dr. Johnson", Salary: 6000},
		{Name: "Dr. Brown", Salary: 4500},
	}

	patients := []Patient{
		{Name: "Alice", Age: 30, Doctor: doctors[0]},
		{Name: "Bob", Age: 45, Doctor: doctors[1]},
		{Name: "Charlie", Age: 25, Doctor: doctors[0]},
	}

	doctorStream := CreateStream(doctors)
	doctorStream.
		Filter(func(d Doctor) bool { return d.Salary > 4500 }).
		Map(func(d Doctor) Doctor { d.Salary += 500; return d }).
		Display()

	patientStream := CreateStream(patients)
	patientStream.
		Distinct().
		Display()

	maxDoctor := doctorStream.Max(func(d1, d2 Doctor) bool { return d1.Salary < d2.Salary })
	if maxDoctor != nil {
		fmt.Println("Doctor with max salary:", maxDoctor.display())
	}

	totalSalary := doctorStream.Reduce(0, func(acc int, d Doctor) int { return acc + d.Salary })
	fmt.Println("Total salary of doctors:", totalSalary)
}