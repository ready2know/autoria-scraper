package main

import (
	"auto-ria-scraper/cmd/scrapper/util"
	"auto-ria-scraper/internal/models"
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// AutoProduct defining a data structure to store the scraped data

func openHTMLFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func parseHTMLFile(file *os.File) []models.AutoProduct {
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		log.Fatal(err)
	}

	autoLots := []models.AutoProduct{}
	doc.Find(".ticket-item ").Each(func(i int, s *goquery.Selection) {
		mainDiv := s.Find("[data-advertisement-data]")
		id, _ := mainDiv.Attr("data-id")
		url, _ := mainDiv.Attr("data-link-to-view")
		brand, _ := mainDiv.Attr("data-mark-name")
		model, _ := mainDiv.Attr("data-model-name")
		year, _ := mainDiv.Attr("data-year")
		userId, _ := mainDiv.Attr("data-user-id")
		generation, _ := mainDiv.Attr("data-generation-name")
		equipment, _ := mainDiv.Attr("data-equipment-name")
		modification, _ := mainDiv.Attr("data-modification-name")

		photo, _ := s.Find(".ticket-photo img").Attr("src")
		price, _ := s.Find(".price-ticket").Attr("data-main-price")
		currency, _ := s.Find(".price-ticket").Attr("data-main-currency")

		carPrice := 0
		if len(price) != 0 {
			if currency != "USD" {
				convertedPrice, _ := strconv.Atoi(price)
				carPrice = convertedPrice / 39
			} else {
				carPrice, _ = strconv.Atoi(price)
			}
		}

		definitionBlock := s.Find(".definition-data")
		mileageBlock := definitionBlock.Find(".icon-mileage").Parent()
		mileage := strings.TrimSpace(mileageBlock.Text())

		if mileage == "без пробігу" {
			mileage = "0"
		} else {
			mileage = strings.Replace(mileage, " тис. км", "", -1)
		}
		mileageNumber, _ := strconv.Atoi(mileage)

		cityBlock := definitionBlock.Find(".icon-location").Parent()
		city := cityBlock.Text()
		city = strings.Replace(city, "( від )", "", -1)

		gearIcon := definitionBlock.Find("[title='Тип коробки передач']")

		gearBlock := gearIcon.Parent()
		gear := gearBlock.Text()

		batteryBlock := definitionBlock.Find(".icon-battery").Parent()
		var fuel, engine string
		if batteryBlock.Length() > 0 {
			fuel = "Електро"
		} else {
			fuelEngineBlock := definitionBlock.Find(".icon-fuel").Parent()
			fuelEngine := fuelEngineBlock.Text()
			fuelEngineParts := strings.Split(fuelEngine, ",")
			fuel = fuelEngineParts[0]
			if len(fuelEngineParts) > 1 {
				engine = fuelEngineParts[1]
			}
		}

		vinBlock := definitionBlock.Find(".label-vin")
		var vin string
		if vinBlock.Length() > 0 {
			vin = strings.TrimSpace(vinBlock.Text())
			if len(vin) > 17 {
				vin = vin[:17]
				if strings.Contains(vin, " ") {
					vin = "checked"
				}
			}
		}

		licensePlateBlock := definitionBlock.Find(".state-num")
		var licensePlateNumber string
		if licensePlateBlock.Length() > 0 {
			licensePlateNumber = licensePlateBlock.Text()
			licensePlateNumber = strings.Replace(licensePlateNumber, "Ми розпізнали держномер авто на фото та перевірили його за реєстрами МВС.", "", -1)
			licensePlateNumber = strings.TrimSpace(licensePlateNumber)
		}

		autoLot := models.AutoProduct{
			Id:           id,
			Url:          strings.TrimSpace(url),
			Mileage:      mileageNumber,
			Fuel:         strings.TrimSpace(fuel),
			Gear:         strings.TrimSpace(gear),
			Price:        carPrice,
			City:         strings.TrimSpace(city),
			Brand:        strings.TrimSpace(brand),
			Model:        strings.TrimSpace(model),
			Year:         strings.TrimSpace(year),
			Engine:       strings.TrimSpace(engine),
			Generation:   strings.TrimSpace(generation),
			UserId:       strings.TrimSpace(userId),
			Image:        strings.TrimSpace(photo),
			Equipment:    strings.TrimSpace(equipment),
			Modification: strings.TrimSpace(modification),
			LicensePlate: licensePlateNumber,
			Vin:          vin,
		}

		autoLots = append(autoLots, autoLot)
	})

	return autoLots
}

func writeToCSV(path string, autoProducts []models.AutoProduct) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	err = writer.Write([]string{
		"id",
		"url",
		"mileage",
		"fuel",
		"gear",
		"price",
		"city",
		"brand",
		"model",
		"year",
		"engine",
		"generation",
		"userId",
		"image",
		"equipment",
		"modification",
		"licensePlate",
		"vin",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Write data
	for _, product := range autoProducts {
		err := writer.Write([]string{
			product.Id,
			product.Url,
			strconv.Itoa(product.Mileage),
			product.Fuel,
			product.Gear,
			strconv.Itoa(product.Price),
			product.City,
			product.Brand,
			product.Model,
			product.Year,
			product.Engine,
			product.Generation,
			product.UserId,
			product.Image,
			product.Equipment,
			product.Modification,
			product.LicensePlate,
			product.Vin,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	util.Init()
	dirPath := "./pages"

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	autoLots := []models.AutoProduct{}
	autoChannel := make(chan []models.AutoProduct)
	channelsOpened := 0
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".html" {
			channelsOpened++
			go func(filepath string, c chan []models.AutoProduct) {
				fileReader := openHTMLFile(filepath)
				c <- parseHTMLFile(fileReader)
				log.Printf("Finished parsing file: %s", filepath)
			}(filepath.Join(dirPath, file.Name()), autoChannel)
		}
	}

	for range channelsOpened {
		autoLots = append(autoLots, <-autoChannel...)
	}

	writeToCSV("./csv/auto.csv", autoLots)
}
