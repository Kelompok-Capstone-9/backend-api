package models

type Classes struct {
	ID uint
	Name string
	Description string
	ClassType ClassType
	InstructorID uint
	Instructor Instructor
}