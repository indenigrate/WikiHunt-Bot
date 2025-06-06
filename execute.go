package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	blueStyler = color.New(color.FgBlue, color.Bold).SprintFunc()
	redStyler  = color.New(color.FgRed, color.Bold).SprintFunc()
)

func blueBold(text string) {
	fmt.Println(blueStyler(text))
}

func redBold(text string) {
	fmt.Println(redStyler(text))
}

// backlinks is true when we try to find a path from end to start
// it is false when we go in the usual way from start to end
func wikiHunt(start string, end string, backlinks bool) {
	// initialising error and slice s that stores all links
	// current store the active page title

	// var err error
	// var s []string

	// current begins from the starting point

	if backlinks {
		temp := start
		start = end
		end = temp
	}
	// current := start
	traversed := make(map[string]bool)

	// text color
	blueBold := color.New(color.FgBlue, color.Bold).SprintFunc()
	// redBold := color.New(color.FgRed, color.Bold).SprintFunc()

	if backlinks {
		fmt.Printf("Going from %s to %s via backlinks\n", blueBold(start), blueBold(end))
	} else {
		fmt.Printf("Going from %s to %s via links\n", blueBold(start), blueBold(end))
	}

	nextGuess([]string{start}, end, false, traversed, 5)
	// nextGuessBFS([]string{start}, end, false, 2)
	// for current != end {
	// 	traversed[current] = true
	// 	// for red text without library
	// 	// fmt.Println("\033[1;31mCurrent: ", current, "\033[0m")

	// 	// for red text with library
	// 	fmt.Println("Current: ", redBold(current))

	// 	s, err = fetchWikiLinks(current, backlinks)
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}
	// 	topNchoicesWithActualSimilarity, err, current := checkSimilarity(end, s, traversed)
	// 	if current == "" {
	// 		fmt.Println("Error: Unable to find maximum similarit element || ", err)
	// 		return
	// 	}
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}
	// 	s = nil
	// }
	// fmt.Println("Current: ", redBold(current))
	// fmt.Println()
}

func nextGuess(start []string, end string, backlinks bool, traversed map[string]bool, depth int) {
	var topNchoicesWithActualSimilarity []choiceWithSimilarity

	for _, current := range start {
		if current != end {
			traversed[current] = true
			redBold(fmt.Sprintf("Current: %s", current))
			s, err := fetchWikiLinks(current, backlinks)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			topNchoicesWithActualSimilarity, err = checkSimilarity(end, s, traversed)
			if len(topNchoicesWithActualSimilarity) == 0 {
				return
			}
			current = topNchoicesWithActualSimilarity[0].Choice
			if current == "" {
				fmt.Println("Error: Unable to find maximum similarit element || ", err)
				return
			}
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}
			s = nil
		} else {
			fmt.Println("REACHED!!!!!!")
			os.Exit(0)
		}
		redBold(fmt.Sprintf("Current: %s", current))
		fmt.Println()

		var nextChoices []string
		if depth > len(topNchoicesWithActualSimilarity) {
			depth = len(topNchoicesWithActualSimilarity)
		}
		for i := 0; i < depth; i++ {
			nextChoices = append(nextChoices, topNchoicesWithActualSimilarity[i].Choice)
		}
		blueBold(fmt.Sprintf("Next choices are: ") + fmt.Sprintln(nextChoices))
		nextGuess(nextChoices, end, backlinks, traversed, depth)
		fmt.Println("!!!!!!!!!!!!! DEPTH 1 OVER !!!!!!!!!!!!!!!!!!")
	}
}

// GPT code
type queueItem struct {
	Title string
	Depth int
}

func nextGuessBFS(start []string, end string, backlinks bool, depthLimit int) {
	queue := []queueItem{}
	traversed := map[string]bool{}

	// Initialize queue with start nodes
	for _, title := range start {
		queue = append(queue, queueItem{Title: title, Depth: 0})
		traversed[title] = true
	}

	for len(queue) > 0 {
		// Dequeue
		item := queue[0]
		queue = queue[1:]

		current := item.Title
		currDepth := item.Depth

		redBold(fmt.Sprintf("Current: %s (depth %d)", current, currDepth))

		if current == end {
			fmt.Println("REACHED!!!!!!")
			os.Exit(0)
		}

		links, err := fetchWikiLinks(current, backlinks)
		if err != nil {
			fmt.Println("Error fetching links for", current, ":", err)
			continue
		}

		topChoices, err := checkSimilarity(end, links, traversed)
		if err != nil {
			fmt.Println("Error checking similarity:", err)
			continue
		}

		// Limit number of next choices to consider
		for i := 0; i < depthLimit && i < len(topChoices); i++ {
			choice := topChoices[i].Choice
			if !traversed[choice] {
				traversed[choice] = true
				queue = append(queue, queueItem{Title: choice, Depth: currDepth + 1})
			}
		}

		fmt.Println()
	}
	fmt.Println("Target not found.")
}
