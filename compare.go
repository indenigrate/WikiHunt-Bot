package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type similarityRequest struct {
	Input1 string `json:"input1"`
	Input2 string `json:"input2"`
}

type similarityResponse struct {
	Similarity float64 `json:"similarity"`
}

func checkSimilarity(target string, choices []string, current string) ([]float64, error, string) {
	var results []float64
	var max float64
	var maxName string
	maxName = "No value found yet"
	max = -1
	url := "http://127.0.0.1:8000/similarity"

	for _, choice := range choices {
		reqBody := similarityRequest{
			Input1: target,
			Input2: choice,
		}

		// Encode the request body
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to encode JSON: %w", err), maxName
		}

		// Send HTTP POST request
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("HTTP request failed: %w", err), maxName
		}
		defer resp.Body.Close()

		// Decode response
		var respData similarityResponse
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err), maxName
		}

		results = append(results, respData.Similarity)
		if respData.Similarity > max && choice != current {
			max = respData.Similarity
			maxName = choice
			fmt.Println("maxx: ", maxName)
		}
	}
	fmt.Println("Compared all possible choices succesfully")
	return results, nil, maxName
}
