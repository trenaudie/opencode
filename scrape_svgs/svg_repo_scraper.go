package scrapesvgs

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/opencode-ai/opencode/internal/logging"
)

func ScrapeSVG(query string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 3
	}

	url := fmt.Sprintf("https://www.svgrepo.com/vectors/%s/", query)
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	
	resp, err := client.Do(req)
	if err != nil {
		logging.Info("Failed to fetch SVGs from SVGRepo:", err)
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Info("Failed to read response body:", err)
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
		logging.Info("Downloaded SVG from URL:", svgURL)

		if err != nil {
			fmt.Printf("Error downloading SVG from %s: %v\n", svgURL, err)
			continue
		}

		svgContents = append(svgContents, svgContent)
	}

	return svgContents, nil
}

func downloadSVG(url string) (string, error) {
	logging.Info("Downloading SVG from URL:", url)
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "image/svg+xml,image/*,*/*;q=0.8")
	
	resp, err := client.Do(req)
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
