package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WikiResponse struct {
	Continue struct {
		PlContinue string `json:"plcontinue"`
	} `json:"continue"`
	Query struct {
		Pages map[string]struct {
			Links []struct {
				Title string `json:"title"`
			} `json:"links"`
		} `json:"pages"`
	} `json:"query"`
}

func fetchWikiLinks(title string) ([]string, error) {
	var allLinks []string
	baseURL := "https://en.wikipedia.org/w/api.php"
	plContinue := ""

	for {
		// Prepare request parameters
		params := url.Values{}
		params.Set("action", "query")
		params.Set("format", "json")
		params.Set("titles", title)
		params.Set("prop", "links")
		params.Set("pllimit", "max")
		if plContinue != "" {
			params.Set("plcontinue", plContinue)
		}

		// Make HTTP request
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

		// Extract links from dynamic page ID
		for _, page := range data.Query.Pages {
			for _, link := range page.Links {
				allLinks = append(allLinks, link.Title)
			}
			break // Only one page is relevant
		}

		// Check if there's more to fetch
		if data.Continue.PlContinue == "" {
			break
		}
		plContinue = data.Continue.PlContinue
	}
	fmt.Println("Fetched all possible links/choices")
	return allLinks, nil

}

type similarityBulkRequest struct {
	Target string   `json:"target"`
	Inputs []string `json:"inputs"`
}

type similarityResponse struct {
	Similarity float64 `json:"similarity"`
}

func checkSimilarity(target string, choices []string, traversed map[string]bool) ([]float64, error, int) {
	var max float64
	maxIndex := -1
	max = -1.1

	//getting definitions
	choicesDefinitionMap, err := getDefinitions(choices)
	choicesDefinition, err := getSliceFromMap(choices, choicesDefinitionMap)
	// for i := 0; i < 5; i++ {
	// 	fmt.Println("i: ", i)
	// 	fmt.Println(choices[i])
	// 	fmt.Println(choicesDefinition[i])
	// }
	// return nil, fmt.Errorf("Ending"), -1

	targetDefinitionMap, err := getDefinitions([]string{target})
	targetDefinition, err := getSliceFromMap(choices, targetDefinitionMap)

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch Definitions: %w", err), maxIndex
	}

	// Get similarities
	similarity, err := getSimilarities(target, choices)
	similarityOfDefinition, err := getSimilarities(targetDefinition[0], choicesDefinition)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch Definitions: %w", err), maxIndex
	}

	// for i := 0; i < 5; i++ {
	// 	fmt.Println("i: ", i)
	// 	fmt.Println("choices: ", choices[i])
	// 	fmt.Println("sim: ", similarity[i])
	// 	fmt.Println("sim of def: ", similarityOfDefinition[i])
	// 	fmt.Println("choices def", choicesDefinition[i][:10])
	// }
	// return nil, fmt.Errorf("Ending"), -1

	// Find max similarity
	for i, sim := range similarity {
		actualSimilarity := sim*0.5 + similarityOfDefinition[i]*0.5
		// actualSimilarity := sim
		if actualSimilarity > max && !traversed[choices[i]] {
			max = actualSimilarity
			maxIndex = i
		}
	}

	fmt.Println("Compared all possible choices successfully")
	if maxIndex == -1 {
		return nil, fmt.Errorf("Unable to find max value"), maxIndex
	}
	fmt.Printf("The max similarity is: %f\n", max)
	return similarity, nil, maxIndex
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
