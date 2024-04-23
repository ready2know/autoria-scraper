package main

import (
	"auto-ria-scraper/cmd/pgImporter/config"
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var appConfig config.AppConfig

func main() {

	config.LoadConfig(&appConfig, "./config.json")

	postgresConfig := appConfig.Postgres

	// Open the CSV file
	f, err := os.Open("./csv/auto.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader
	reader := csv.NewReader(f)

	// Read and discard the header line
	_, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the PostgreSQL database
	dataSourceName := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", postgresConfig.User, postgresConfig.Password, postgresConfig.Dbname)
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Begin a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the SQL statement
	stmt, err := tx.Prepare(`INSERT INTO 
    	autoria(id,url,mileage,fuel,gear,price,city,brand,model,year,engine,generation,userId,image,equipment,modification,licensePlate,vin) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18) 
		ON CONFLICT (id) DO UPDATE SET (price, updated)=(excluded.price, DEFAULT);`)
	if err != nil {
		log.Fatal(err)
	}

	// Read the CSV file line by line
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Execute the prepared SQL statement with the data from the CSV line
		_, err = stmt.Exec(
			record[0], record[1],
			record[2], record[3],
			record[4], record[5],
			record[6], record[7],
			record[8], record[9],
			record[10], record[11],
			record[12], record[13],
			record[14], record[15],
			record[16], record[17],
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Commit the transaction
	tx.Commit()
}
