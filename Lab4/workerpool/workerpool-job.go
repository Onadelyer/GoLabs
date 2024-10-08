package workerpool

type Job struct {
	Id          int
	Description string
	Run         func()
}