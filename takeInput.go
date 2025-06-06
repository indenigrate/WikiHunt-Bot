package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func fetchWikipediaTitles(query string) ([]string, error) {
	escapedQuery := url.QueryEscape(query)

	url := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=opensearch&format=json&search=%s", escapedQuery)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	var result []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON decoding failed: %w", err)
	}

	titlesRaw, ok := result[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	var titles []string
	for _, t := range titlesRaw {
		if titleStr, ok := t.(string); ok {
			titles = append(titles, titleStr)
		}
	}

	return titles, nil
}

func takeInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		// fmt.Print("Enter search term for Wikipedia: ")
		input, _ := reader.ReadString('\n')
		query := strings.TrimSpace(input)

		if query == "" {
			fmt.Println("Please enter a non-empty search term.")
			continue
		}

		titles, err := fetchWikipediaTitles(query)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if len(titles) == 0 {
			fmt.Println("No results found.")
			continue
		}

		prompt := promptui.Select{
			Label: "Select Wikipedia Page",
			Items: titles,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Println("Prompt failed:", err)
			return "", err
		}
		return titles[i], nil
	}
}
