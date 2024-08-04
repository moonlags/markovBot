package main

import (
	"encoding/gob"
	"os"

	"github.com/moonlags/markovBot/internal/markov"
)

type gobData struct {
	Chain  map[string][]string
	Images []string
}

func (s *server) loadGobData(chain *markov.Chain) error {
	file, err := os.Open("data.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	var data gobData
	if err := gob.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	s.images = data.Images
	chain.Chain = data.Chain

	return nil
}

func (s *server) saveGobData(chain *markov.Chain) error {
	file, err := os.Create("data.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	data := gobData{
		Chain:  chain.Chain,
		Images: s.images,
	}

	return gob.NewEncoder(file).Encode(data)
}
