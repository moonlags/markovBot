package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type data struct {
	Chain  map[string][]string
	Images []string
}

func main() {
	file, err := os.Open("data.gob")
	if err != nil {
		log.Fatal("Failed to open file:", err)
	}
	defer file.Close()

	data := new(data)
	if err := gob.NewDecoder(file).Decode(data); err != nil {
		log.Fatal("Failed to decode gob:", err)
	}

	for key, value := range data.Chain {
		fmt.Printf("%s: %v\n", key, value)
	}
}
