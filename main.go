package main

import "fmt"

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

	start := "Silver"
	end := "Light"
	current := start
	var s []string
	var err error

	for current != end {
		fmt.Println(current)
		s, err = fetchWikiLinks(current)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		_, err, current = checkSimilarity(end, s, current)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		s = nil
	}
	fmt.Println(current)
}
