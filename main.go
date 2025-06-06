package main

import (
	"fmt"
	"os"
)

func main() {
	// Take input during runtime

	fmt.Printf("Title of starting page: ")
	start, err := takeInput()
	fmt.Printf("Title of ending page: ")
	end, err := takeInput()
	if err != nil {
		os.Exit(1)
	}

	// start := "Acid"
	// start := "Malaria"
	// start := "Cat"
	// start := "Electromagnetism"
	// start := "Art Deco"
	// start := "Diethylaminoethanol"
	// start := "Silver"
	// end := "Lysergic acid diethylamide"
	// end := "Al-Qaeda"
	// end := "Eiffel Tower"
	// end := "Spacecraft"
	// end := "Geophysics"
	// end := "Light"
	// current := "Silver"
	// current := "Dimethylaminoethanol"
	// current := "Diethylaminoethanol"

	wikiHunt(start, end, false)
	// wikiHunt(start, end, true)
}
