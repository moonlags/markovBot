package imgflip

import (
	"bytes"
	"encoding/json"
)

type captionBody struct {
	TemplateID string `json:"template_id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Text0      string `json:"text0"`
	Text1      string `json:"text1"`
}

func (c *Config) captionBody(id, text0, text1 string) (*bytes.Reader, error) {
	body := captionBody{
		TemplateID: id,
		Username:   c.Username,
		Password:   c.Password,
		Text0:      text0,
		Text1:      text1,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
