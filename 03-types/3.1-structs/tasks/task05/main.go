package main

import "fmt"

// TODO: объяви структуру Audit с полями CreatedAt (string), UpdatedAt (string)
type Audit struct {
	CreatedAt string
	UpdatedAt string
}

// TODO: объяви структуру Article с полем Title (string) и встроенным полем Audit (embedding)
type Article struct {
	Title string
	Audit
}

func main() {
	// TODO: создай значение Article:
	// Title: "Go Basics"
	// Audit: {CreatedAt: "2024-01-01", UpdatedAt: "2024-06-01"}
	a := Article{
		Title: "Go Basics",
		Audit: Audit{
			CreatedAt: "2024-01-01",
			UpdatedAt: "2024-06-01",
		},
	}

	fmt.Println("Title:", a.Title)
	fmt.Println("CreatedAt:", a.CreatedAt)
	fmt.Println("UpdatedAt:", a.UpdatedAt)
}