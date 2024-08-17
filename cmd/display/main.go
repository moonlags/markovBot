package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("data.gob")
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}
	defer file.Close()

	var data map[string][]string
	if err := gob.NewDecoder(file).Decode(&data); err != nil {
		log.Fatal("Failed to decode gob:", err)
	}

	for key, value := range data {
		fmt.Printf("%s: %v\n", key, value)
	}
}
