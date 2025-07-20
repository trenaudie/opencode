package main

import (
	"fmt"
	"log"

	"github.com/opencode-ai/opencode/scrape_svgs"
)

func main() {
	fmt.Println("Testing ScrapeSVG function with query 'hospital'...")
	
	svgs, err := scrapesvgs.ScrapeSVG("hospital", 3)
	if err != nil {
		log.Fatalf("Error scraping SVGs: %v", err)
	}

	fmt.Printf("Successfully scraped %d SVG files:\n\n", len(svgs))
	
	for i, svg := range svgs {
		fmt.Printf("=== SVG %d ===\n", i+1)
		fmt.Printf("Length: %d characters\n", len(svg))
		fmt.Printf("Preview (first 200 chars): %s...\n\n", truncate(svg, 200))
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}