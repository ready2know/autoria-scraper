package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type AppConfig struct {
	QueryParams struct {
		AbroadNot           string `json:"abroad.not"`
		CategoriesMainId    string `json:"categories.main.id"`
		CountryImportUsaNot string `json:"country.import.usa.not"`
		CountryOriginId1Not string `json:"country.origin.id[1].not"`
		CountryOriginId2Not string `json:"country.origin.id[2].not"`
		CustomNot           string `json:"custom.not"`
		IndexName           string `json:"indexName"`
		MileageLte          string `json:"mileage.lte"`
		PriceCurrency       string `json:"price.currency"`
		PriceUSDGte         string `json:"price.USD.gte"`
		PriceUSDLte         string `json:"price.USD.lte"`
		Size                string `json:"size"`
		SortOrder           string `json:"sort[0].order"`
		YearGte             string `json:"year[0].gte"`
	} `json:"queryParams"`
	StartPage      int `json:"startPage"`
	EndPage        int `json:"endPage"`
	RequestTimeout int `json:"requestTimeout"`
}

func LoadConfig(appConfig *AppConfig, filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&appConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}
