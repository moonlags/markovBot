package main

import (
	"encoding/gob"
	"os"
)

func (s *server) loadGobData() error {
	file, err := os.Open("data.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewDecoder(file).Decode(&s.chain.Chain)
}

func (s *server) saveGobData() error {
	file, err := os.Create("data.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(s.chain.Chain)
}
