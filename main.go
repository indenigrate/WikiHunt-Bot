package main

import (
	"fmt"

	"github.com/fatih/color"
)

func main() {
	// // s, err := fetchWikiLinks("Silver")
	// // s, err := fetchWikiLinks("Litharge")
	// // s, err := fetchWikiLinks("Gold")
	// s, err := fetchWikiLinks("Color")
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }

	// similar, err, maxName := checkSimilarity("light", s)

	// for index, value := range s {
	// 	if similar[index] > 0.4 {
	// 		// fmt.Printf("Value%d:Similarity:%f\n", index, similar[index])
	// 		fmt.Printf("Value%d: %s Similarity:%f\n", index, value, similar[index])

	// 	}
	// }
	// fmt.Println("Max name is : ", maxName)

	// start := "Poison"
	start := "Malaria"
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
	var err error
	var currentIndex int
	traversed := make(map[string]bool)
	fmt.Printf("From %s to %s\n", start, end)
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
			_, err, currentIndex = checkSimilarity(end+" and "+start, s, traversed)
			if currentIndex == -1 {
				fmt.Println("Error: Index out of range")
				return
			}
			current = s[currentIndex]
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			flag = true
		} else {
			_, err, currentIndex = checkSimilarity(end, s, traversed)
			if currentIndex == -1 {
				fmt.Println("Error: Index out of range")
				return
			}
			current = s[currentIndex]
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
		}

		s = nil
	}
	color.New(color.FgRed, color.Bold).Println("Current: ", current)

}
