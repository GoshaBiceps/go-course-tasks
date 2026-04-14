package main

import "fmt"

// TODO: объяви структуру ContactInfo с полями Phone (string), Email (string)
type ContactInfo struct {
	Phone string
	Email string
}

// TODO: объяви структуру Address с полями City (string), Street (string)
type Address struct {
	City   string
	Street string
}

// TODO: объяви структуру Client с полями:
//       ID (string), Address Address, и встроенным ContactInfo (embedding)
type Client struct {
	ID string
	Address Address
	ContactInfo
}

func main() {
	// TODO: создай значение Client:
	// ID: "C-01"
	// Address: {City: "Moscow", Street: "Tverskaya"}
	// ContactInfo: {Phone: "+7-999-000", Email: "ivan@example.com"}
	c := Client{
		ID: "C-01",
		Address: Address{
			City:   "Moscow",
			Street: "Tverskaya",
		},
		ContactInfo: ContactInfo{
			Phone: "+7-999-000",
			Email: "ivan@example.com",
		},
	}

	fmt.Println("ID:", c.ID)
	fmt.Println("City:", c.Address.City)
	fmt.Println("Email:", c.Email) // доступ через embedding
}