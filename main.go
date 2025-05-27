package main

import (
	"fmt"

	"github.com/fatih/color"
)

func main() {
	var err error
	// Take input during runtime
	// fmt.Printf("Title of starting page: ")
	// start, err := takeInput()
	// fmt.Printf("Title of ending page: ")
	// end, err := takeInput()
	// if err != nil {
	// 	os.Exit(1)
	// }

	start := "Acid"
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
	end := "Geophysics"
	// end := "Light"
	current := start
	// current := "Silver"
	// current := "Dimethylaminoethanol"
	// current := "Diethylaminoethanol"

	flag := false

	var s []string

	traversed := make(map[string]bool)
	fmt.Printf("Going from %s to %s\n", start, end)
	for current != end {
		traversed[current] = true
		// for red text without library
		// fmt.Println("\033[1;31mCurrent: ", current, "\033[0m")

		// for red text with library
		color.New(color.FgRed, color.Bold).Println("Current: ", current)

		s, err = fetchWikiLinks(current)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		if flag {
			_, err, current = checkSimilarity(end+" and "+start, s, traversed)
			if current == "" {
				fmt.Println("Error: Index out of range")
				return
			}
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			flag = true
		} else {
			_, err, current = checkSimilarity(end, s, traversed)
			if current == "" {
				fmt.Println("Error: Index out of range")
				return
			}
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}

		s = nil
	}
	color.New(color.FgRed, color.Bold).Println("Current: ", current)

}
