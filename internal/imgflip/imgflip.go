package imgflip

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	Username string
	Password string
}

type Meme struct {
	ID string `json:"id"`
}

func New(username string, password string) *Config {
	return &Config{
		Username: username,
		Password: password,
	}
}

func GetMemes() ([]Meme, error) {
	resp, err := http.Get("https://api.imgflip.com/get_memes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body struct {
		Data struct {
			Memes []Meme `json:"memes"`
		} `json:"data"`
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	if !body.Success {
		return nil, fmt.Errorf("error occured")
	}

	return body.Data.Memes, nil
}

func (c *Config) MemeWithCaption(id, text0, text1 string) (string, error) {
	req, err := http.NewRequest("POST", "https://api.imgflip.com/caption_image", nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("template_id", id)
	q.Add("username", c.Username)
	q.Add("password", c.Password)
	q.Add("text0", text0)
	q.Add("text1", text1)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
		ErrorMessage string `json:"error_message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}

	if body.ErrorMessage != "" {
		return "", fmt.Errorf("imgflip error: %s", body.ErrorMessage)
	}

	return body.Data.URL, nil
}
