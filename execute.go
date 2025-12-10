package main

import (
	"log"
	"time"
)

var (
	timeStart = time.Now()
)

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

	if backlinks {
		log.Printf("Going from %s to %s via backlinks\n", start, end)
	} else {
		log.Printf("Going from %s to %s via links\n", start, end)
	}
	timeStart = time.Now()
	nextGuess([]string{start}, end, false, traversed, 5)

	// nextGuessBFS([]string{start}, end, false, 2)
	// for current != end {
	// 	traversed[current] = true
	// 	// for red text without library
	// 	// log.Println("\033[1;31mCurrent: ", current, "\033[0m")

	// 	// for red text with library
	// 	log.Println("Current: ", redBold(current))

	// 	s, err = fetchWikiLinks(current, backlinks)
	// 	if err != nil {
	// 		log.Println("Error: ", err)
	// 		return
	// 	}
	// 	topNchoicesWithActualSimilarity, err, current := checkSimilarity(end, s, traversed)
	// 	if current == "" {
	// 		log.Println("Error: Unable to find maximum similarit element || ", err)
	// 		return
	// 	}
	// 	if err != nil {
	// 		log.Println("Error: ", err)
	// 		return
	// 	}
	// 	s = nil
	// }
	// log.Println("Current: ", redBold(current))
	// log.Println()
}

func nextGuess(start []string, end string, backlinks bool, traversed map[string]bool, depth int) {
	var topNchoicesWithActualSimilarity []choiceWithSimilarity

	for _, current := range start {
		if current != end {
			traversed[current] = true
			log.Printf("Current: %s", current)
			s, err := fetchWikiLinks(current, backlinks)
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			topNchoicesWithActualSimilarity, err = checkSimilarity(end, s, traversed)
			if len(topNchoicesWithActualSimilarity) == 0 {
				return
			}
			current = topNchoicesWithActualSimilarity[0].Choice
			if current == "" {
				log.Println("Error: Unable to find maximum similarit element || ", err)
				return
			}
			if err != nil {
				log.Println("Error: ", err)
				return
			}
			s = nil
		} else {
			log.Println("REACHED!!!!!!")
			elapsed := time.Since(timeStart)
			log.Printf("Traversal time: %s\n", elapsed)
			return // return instead of os.Exit
		}
		log.Printf("Current: %s", current)
		log.Println()

		var nextChoices []string
		if depth > len(topNchoicesWithActualSimilarity) {
			depth = len(topNchoicesWithActualSimilarity)
		}
		for i := 0; i < depth; i++ {
			nextChoices = append(nextChoices, topNchoicesWithActualSimilarity[i].Choice)
		}
		log.Printf("Next choices are: %v\n", nextChoices)
		nextGuess(nextChoices, end, backlinks, traversed, depth)
		log.Println("!!!!!!!!!!!!! DEPTH 1 OVER !!!!!!!!!!!!!!!!!!")
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

		log.Printf("Current: %s (depth %d)", current, currDepth)

		if current == end {
			log.Println("REACHED!!!!!!")
			return // return instead of os.Exit
		}

		links, err := fetchWikiLinks(current, backlinks)
		if err != nil {
			log.Println("Error fetching links for", current, ":", err)
			continue
		}

		topChoices, err := checkSimilarity(end, links, traversed)
		if err != nil {
			log.Println("Error checking similarity:", err)
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

		log.Println()
	}
	log.Println("Target not found.")
}
