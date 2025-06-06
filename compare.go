package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type WikiResponse struct {
	Continue struct {
		PlContinue string `json:"plcontinue,omitempty"`
		BlContinue string `json:"blcontinue,omitempty"`
	} `json:"continue,omitempty"`

	Query struct {
		Pages map[string]struct {
			Links []LinkItem `json:"links,omitempty"`
		} `json:"pages,omitempty"`

		Backlinks []LinkItem `json:"backlinks,omitempty"`
	} `json:"query"`
}

type LinkItem struct {
	Title  string `json:"title"`
	PageID int    `json:"pageid,omitempty"`
	NS     int    `json:"ns,omitempty"`
}

// make changes to incorporate backlinks
func fetchWikiLinks(title string, backlinks bool) ([]string, error) {
	var allLinks []string
	baseURL := "https://en.wikipedia.org/w/api.php"
	continueToken := ""

	for {
		params := url.Values{}
		params.Set("action", "query")
		params.Set("format", "json")

		if backlinks {
			// Fetch backlinks
			params.Set("list", "backlinks")
			params.Set("bltitle", title)
			params.Set("bllimit", "max")
			if continueToken != "" {
				params.Set("blcontinue", continueToken)
			}
		} else {
			// Fetch internal page links
			params.Set("titles", title)
			params.Set("prop", "links")
			params.Set("pllimit", "max")
			if continueToken != "" {
				params.Set("plcontinue", continueToken)
			}
		}

		resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var data WikiResponse
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, err
		}

		if backlinks {
			// Collect backlinks
			for _, link := range data.Query.Backlinks {
				allLinks = append(allLinks, link.Title)
			}
			if data.Continue.BlContinue == "" {
				break
			}
			continueToken = data.Continue.BlContinue
		} else {
			// Collect internal links
			for _, page := range data.Query.Pages {
				for _, link := range page.Links {
					allLinks = append(allLinks, link.Title)
				}
				break // Only one page is relevant
			}
			if data.Continue.PlContinue == "" {
				break
			}
			continueToken = data.Continue.PlContinue
		}
	}
	if backlinks {
		fmt.Println("Fetched all possible backlinks")
	} else {
		fmt.Println("Fetched all possible links")
	}
	return allLinks, nil
}

type similarityBulkRequest struct {
	Target string   `json:"target"`
	Inputs []string `json:"inputs"`
}

type similarityResponse struct {
	Similarity float64 `json:"similarity"`
}

type choiceWithSimilarity struct {
	Choice     string
	Similarity float64
}

func checkSimilarity(target string, choices []string, traversed map[string]bool) ([]choiceWithSimilarity, error) {
	// var max float64
	// maxElement := "test"
	// max = -1.1
	//check similarity for the choices
	similarity, err := getSimilarities(target, choices)
	if err != nil {
		return nil, fmt.Errorf("Unable to get Similarites: %w", err)
	}

	// storing choices with similarity
	var choicesWithSimilarity []choiceWithSimilarity
	for i, choice := range choices {
		choicesWithSimilarity = append(choicesWithSimilarity, choiceWithSimilarity{
			Choice:     choice,
			Similarity: similarity[i],
		})
	}

	// extracting top N similarity by sorting
	sort.Slice(choicesWithSimilarity, func(i, j int) bool {
		return choicesWithSimilarity[i].Similarity > choicesWithSimilarity[j].Similarity
	})

	topN := len(choices)
	if topN > 100 {
		topN = 99
	}
	var topNChoices []string
	for i := range topN {
		topNChoices = append(topNChoices, choicesWithSimilarity[i].Choice)
	}
	// fmt.Println(topNChoices)

	//getting definitions for topN choices
	topNchoicesDefinitionMap, err := getDefinitions(topNChoices)
	topNchoicesDefinition, err := getSliceFromMap(topNChoices, topNchoicesDefinitionMap)

	// for i := 0; i < 5; i++ {
	// 	fmt.Println("i: ", i)
	// 	fmt.Println(choices[i])
	// 	fmt.Println(choicesDefinition[i])
	// }
	// return nil, fmt.Errorf("Ending"), -1

	targetDefinitionMap, err := getDefinitions([]string{target})
	var targetDefinition string
	for _, def := range targetDefinitionMap {
		targetDefinition = def
	}
	// targetDefinition, err := getSliceFromMap(choices, targetDefinitionMap)

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch Definitions: %w", err)
	}

	// Get similarities
	similarityOfDefinition, err := getSimilarities(targetDefinition, topNchoicesDefinition)
	if err != nil {
		return nil, fmt.Errorf("Unable to get Similarites: %w", err)
	}

	// Find max similarity
	// for i, sim := range similarityOfDefinition {
	// 	actualSimilarity := choicesWithSimilarity[i].Similarity*0.5 + sim*0.5
	// 	// actualSimilarity := sim
	// 	if actualSimilarity > max && !traversed[choicesWithSimilarity[i].Choice] {
	// 		max = actualSimilarity
	// 		maxElement = choicesWithSimilarity[i].Choice
	// 	}
	// }

	// Calculate actual similarity
	var topNchoicesWithActualSimilarity []choiceWithSimilarity
	for i, sim := range similarityOfDefinition {
		if traversed[choicesWithSimilarity[i].Choice] {
			continue
		}
		topNchoicesWithActualSimilarity = append(topNchoicesWithActualSimilarity, choiceWithSimilarity{
			Choice:     choicesWithSimilarity[i].Choice,
			Similarity: choicesWithSimilarity[i].Similarity*0.5 + sim*0.5})
	}

	sort.Slice(topNchoicesWithActualSimilarity, func(i, j int) bool {
		return topNchoicesWithActualSimilarity[i].Similarity > topNchoicesWithActualSimilarity[j].Similarity // descending order
	})

	// maxElement = topNchoicesWithActualSimilarity[0].Choice
	// max = topNchoicesWithActualSimilarity[0].Similarity

	// topNchoicesWithActualSimilarity := choicesWithSimilarity[:topN]
	// for i := range topNchoicesWithActualSimilarity {
	// 	topNchoicesWithActualSimilarity[i].Similarity = topNchoicesWithActualSimilarity[i].Similarity*0.5 + similarityOfDefinition[i]*0.5
	// }
	// sort.Slice(topNchoicesWithActualSimilarity, func(i, j int) bool {
	// 	return topNchoicesWithActualSimilarity[i].Similarity > topNchoicesWithActualSimilarity[j].Similarity // descending order
	// })
	// i := 0
	// for traversed[topNchoicesWithActualSimilarity[i].Choice] && i < topN {
	// 	i = i + 1
	// }
	// maxElement = topNchoicesWithActualSimilarity[i].Choice
	// max = topNchoicesWithActualSimilarity[i].Similarity

	// fmt.Println(topNchoicesWithActualSimilarity[:5])

	fmt.Println("Compared all possible choices successfully")
	// if maxElement == "" {
	// 	return nil, fmt.Errorf("Unable to find max value"), maxElement
	// }
	// fmt.Printf("The max similarity is: %f\n", max)
	return topNchoicesWithActualSimilarity, nil
}

func getSimilarities(target string, choices []string) ([]float64, error) {

	//timimng execution
	start := time.Now()

	url := "http://127.0.0.1:8000/similarity"
	reqBody := similarityBulkRequest{
		Target: target,
		Inputs: choices,
	}
	// Encode JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %w", err)
	}

	// Send POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var respData struct {
		Similarities []float64 `json:"similarities"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	//timing execution
	elapsed := time.Since(start)
	fmt.Printf("Execution of getSimilarities took %s\n", elapsed)
	return respData.Similarities, nil
}

func getDefinitions(choices []string) (map[string]string, error) {
	const endpoint = "https://en.wikipedia.org/w/api.php"
	definitionMap := make(map[string]string)

	for i := 0; i < len(choices); i += 50 {
		end := i + 50
		if end > len(choices) {
			end = len(choices)
		}
		batch := choices[i:end]
		titles := strings.Join(batch, "|")

		// Prepare query parameters
		params := url.Values{}
		params.Set("action", "query")
		params.Set("format", "json")
		params.Set("prop", "extracts")
		params.Set("exintro", "true")
		params.Set("explaintext", "true")
		params.Set("titles", titles)

		// fmt.Println(endpoint + "?" + params.Encode())
		// return nil, fmt.Errorf("ending")

		// Send the HTTP GET request
		resp, err := http.Get(endpoint + "?" + params.Encode())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch definitions: %v", err)
		}
		defer resp.Body.Close()

		// Decode the response
		var result struct {
			Query struct {
				Pages map[string]struct {
					Title   string `json:"title"`
					Extract string `json:"extract"`
				} `json:"pages"`
			} `json:"query"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("failed to decode JSON: %v", err)
		}

		// Collect definitions in arbitrary page order
		// maxLen := 1000
		// for _, page := range result.Query.Pages {
		// 	if page.Extract != "" {
		// 		// if len(page.Extract) > maxLen {
		// 		// 	page.Extract = page.Extract[:maxLen]
		// 		// }
		// 		definitions = append(definitions, page.Extract)
		// 	} else {
		// 		// definitions = append(definitions, fmt.Sprintf("No definition found for %s.", page.Title))
		// 		definitions = append(definitions, page.Title)
		// 	}
		// }
		for _, page := range result.Query.Pages {
			if page.Extract != "" {
				definitionMap[page.Title] = page.Extract
			} else {
				definitionMap[page.Title] = page.Title
			}
		}

	}
	fmt.Println("Fetched all definitions")
	return definitionMap, nil
}

func getSliceFromMap(choices []string, definition map[string]string) ([]string, error) {
	var s []string
	for _, j := range choices {
		s = append(s, definition[j])
	}
	return s, nil
}
