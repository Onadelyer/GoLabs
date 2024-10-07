package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	students := make(map[string][]int)
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("Введіть номер опції: ")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			createStudent(students, reader)
		case "2":
			addGrade(students, reader)
		case "3":
			printStudentGrades(students, reader)
		case "4":
			printStudentAverage(students, reader)
		case "5":
			printAllStudents(students)
		case "6":
			fmt.Println("Вихід з програми.")
			return
		default:
			fmt.Println("Некоректний вибір. Спробуйте ще раз.")
		}
	}
}

func printMenu() {
	fmt.Println("\nОберіть опцію:")
	fmt.Println("1. Створити студента")
	fmt.Println("2. Додати оцінку студенту")
	fmt.Println("3. Вивести студента з оцінками")
	fmt.Println("4. Вивести середню оцінку студента")
	fmt.Println("5. Вивести всіх студентів")
	fmt.Println("6. Вийти з програми")
}

func createStudent(students map[string][]int, reader *bufio.Reader) {
	fmt.Print("Введіть ім'я студента: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if _, exists := students[name]; exists {
		fmt.Println("Студент з таким ім'ям вже існує.")
	} else {
		students[name] = []int{}
		fmt.Println("Студента успішно створено.")
	}
}

func addGrade(students map[string][]int, reader *bufio.Reader) {
	fmt.Print("Введіть ім'я студента: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if _, exists := students[name]; !exists {
		fmt.Println("Студента з таким ім'ям не знайдено.")
	} else {
		fmt.Print("Введіть оцінку (0-100): ")
		gradeStr, _ := reader.ReadString('\n')
		gradeStr = strings.TrimSpace(gradeStr)
		grade, err := strconv.Atoi(gradeStr)
		if err != nil || grade < 0 || grade > 100 {
			fmt.Println("Некоректне значення оцінки.")
		} else {
			students[name] = append(students[name], grade)
			fmt.Println("Оцінку додано.")
		}
	}
}

func printStudentGrades(students map[string][]int, reader *bufio.Reader) {
	fmt.Print("Введіть ім'я студента: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if grades, exists := students[name]; !exists {
		fmt.Println("Студента з таким ім'ям не знайдено.")
	} else {
		fmt.Printf("Оцінки студента %s: %v\n", name, grades)
	}
}

func printStudentAverage(students map[string][]int, reader *bufio.Reader) {
	fmt.Print("Введіть ім'я студента: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	if grades, exists := students[name]; !exists {
		fmt.Println("Студента з таким ім'ям не знайдено.")
	} else if len(grades) == 0 {
		fmt.Println("У студента немає оцінок.")
	} else {
		sum := 0
		for _, grade := range grades {
			sum += grade
		}
		average := float64(sum) / float64(len(grades))
		fmt.Printf("Середня оцінка студента %s: %.2f\n", name, average)
	}
}

func printAllStudents(students map[string][]int) {
	if len(students) == 0 {
		fmt.Println("Немає створених студентів.")
		return
	}
	fmt.Println("\nСписок всіх студентів:")
	fmt.Printf("%-20s %s\n", "Ім'я студента", "Оцінки")
	fmt.Println(strings.Repeat("-", 40))
	for name, grades := range students {
		fmt.Printf("%-20s %v\n", name, grades)
	}
}