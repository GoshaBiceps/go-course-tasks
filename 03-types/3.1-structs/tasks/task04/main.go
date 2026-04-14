package main

import "fmt"

// TODO: объяви структуру Package с полями ID (string), Weight (int)
type Package struct {
	ID     string
	Weight int
}

// TODO: объяви структуру Destination с полями City (string), Zip (string)
type Destination struct {
	City string
	Zip  string
}

// TODO: объяви структуру Shipment с полями Package (Package), Destination (Destination)
type Shipment struct {
	Package     Package
	Destination Destination
}

func main() {
	// TODO: создай значение Shipment:
	// Package: {ID: "PKG-001", Weight: 5}
	// Destination: {City: "Moscow", Zip: "101000"}
	s := Shipment{
		Package: Package{
			ID:     "PKG-001",
			Weight: 5,
		},
		Destination: Destination{
			City: "Moscow",
			Zip:  "101000",
		},
	}

	fmt.Printf("Package: %s, City: %s\n", s.Package.ID, s.Destination.City)
}