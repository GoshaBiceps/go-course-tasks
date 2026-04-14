package main

import "fmt"

// TODO: объяви структуру ContactInfo с полями Phone (string), Email (string)
type ContactInfo struct {
	Phone string
	Email string
}

// TODO: объяви структуру Course с полями ID (string), Title (string)
type Course struct {
	ID    string
	Title string
}

// TODO: объяви структуру Student с полями Name (string) и встроенным ContactInfo (embedding)
type Student struct {
	Name string
	ContactInfo
}

// TODO: объяви структуру CourseEnrollment с полями:
//       EnrollmentID (string), Status (string)
//       Student Student, Course Course
type CourseEnrollment struct {
	EnrollmentID string
	Status       string
	Student      Student
	Course       Course
}

func main() {
	// TODO: создай значение CourseEnrollment с осмысленными данными
	e := CourseEnrollment{
		EnrollmentID: "E-001",
		Status:       "active",
		Student: Student{
			Name: "Alice",
			ContactInfo: ContactInfo{
				Phone: "+123456789",
				Email: "alice@example.com",
			},
		},
		Course: Course{
			ID:    "C-101",
			Title: "Go Basics",
		},
	}

	fmt.Printf("Enrollment %s: %s enrolled in %s [%s]\n",
		e.EnrollmentID, e.Student.Name, e.Course.Title, e.Status)
}