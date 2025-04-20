package main

import (
	"auto-ria-scraper/cmd/downloader/config"
	"auto-ria-scraper/cmd/downloader/util"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

var appConfig config.AppConfig

func getSearchParams(pageNumber int) url.Values {
	params := url.Values{}

	params.Add("abroad.not", appConfig.QueryParams.AbroadNot)
	params.Add("categories.main.id", appConfig.QueryParams.CategoriesMainId)
	params.Add("country.import.usa.not", appConfig.QueryParams.CountryImportUsaNot)
	params.Add("country.origin.id[1].not", appConfig.QueryParams.CountryOriginId1Not)
	params.Add("country.origin.id[2].not", appConfig.QueryParams.CountryOriginId2Not)
	params.Add("custom.not", appConfig.QueryParams.CustomNot)
	params.Add("indexName", appConfig.QueryParams.IndexName)
	params.Add("mileage.lte", appConfig.QueryParams.MileageLte)
	params.Add("price.currency", appConfig.QueryParams.PriceCurrency)
	params.Add("price.USD.gte", appConfig.QueryParams.PriceUSDGte)
	params.Add("price.USD.lte", appConfig.QueryParams.PriceUSDLte)
	params.Add("size", appConfig.QueryParams.Size)
	params.Add("sort[0].order", appConfig.QueryParams.SortOrder)
	params.Add("year[0].gte", appConfig.QueryParams.YearGte)
	params.Add("gearbox[0]", appConfig.QueryParams.Gearbox)
	params.Add("gearbox[1]", appConfig.QueryParams.Gearbox2)
	params.Add("gearbox[2]", appConfig.QueryParams.Gearbox3)
	params.Add("gearbox[3]", appConfig.QueryParams.Gearbox4)
	params.Add("page", strconv.Itoa(pageNumber))

	return params
}

func getPage(pageNumber int) {
	baseURL := "https://auto.ria.com/uk/search/"

	// Create a new URL object
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Add GET parameters to the URL
	u.RawQuery = getSearchParams(pageNumber).Encode()

	// Send HTTP GET request
	response, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Print the HTTP response status
	log.Printf("%s [%s]: %s", response.Request.URL.RawPath+response.Request.URL.RawQuery, response.Request.Method, response.Status)

	if response.StatusCode != http.StatusOK {
		log.Printf("Stopping requesting")
		os.Exit(1)
	}

	// Create a new file
	fileName := fmt.Sprintf("pages/autoria_page_%d.html", pageNumber)
	out, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Copy data from HTTP response to file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Print the HTTP response status
	log.Printf("Saved into %s", fileName)
}

func main() {
	util.Init()
	config.LoadConfig(&appConfig, "cmd/downloader/config.json")
	// Base URL of the page to download
	for i := appConfig.StartPage; i < appConfig.EndPage; i++ {
		getPage(i)
		time.Sleep(time.Duration(appConfig.RequestTimeout) * time.Second)
	}
}
