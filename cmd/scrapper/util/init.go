package util

import (
	"log"
	"os"
)

func Init() {
	dirs := []string{"./pages", "./csv"}

	for _, dir := range dirs {
		// Check if the directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// The directory does not exist, create it
			err := os.Mkdir(dir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
