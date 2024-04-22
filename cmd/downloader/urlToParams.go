package main

import (
	"fmt"
	"net/url"
)

func main() {
	// The URL you want to parse
	rawURL := "?indexName=auto,order_auto,newauto_search&categories.main.id=1&country.origin.id[1].not=643&country.origin.id[2].not=804&country.import.usa.not=-1&price.USD.gte=4000&price.currency=1&abroad.not=0&custom.not=1&page=0&size=100"

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Get the parameters from the parsed URL
	params := parsedURL.Query()
	for key, value := range params {
		fmt.Printf("params.Add(\"%s\", \"%s\")\n", key, value[0])
	}

	// Print the parameters
}
