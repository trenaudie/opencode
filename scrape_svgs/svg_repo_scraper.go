package scrapesvgs

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func ScrapeSVG(query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 3
	}

	url := fmt.Sprintf("https://www.svgrepo.com/vectors/%s/", query)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	html := string(body)
	
	svgURLRegex := regexp.MustCompile(`src="(https://www\.svgrepo\.com/show/\d+/[^"]+\.svg)"`)
	matches := svgURLRegex.FindAllStringSubmatch(html, limit)
	
	var svgContents []string
	
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		
		svgURL := match[1]
		svgContent, err := downloadSVG(svgURL)
		if err != nil {
			fmt.Printf("Error downloading SVG from %s: %v\n", svgURL, err)
			continue
		}
		
		svgContents = append(svgContents, svgContent)
	}
	
	return svgContents, nil
}

func downloadSVG(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download SVG: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read SVG content: %v", err)
	}

	return strings.TrimSpace(string(body)), nil
}