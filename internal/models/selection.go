package models

import "encoding/json"

type Selection struct {
	Tracks []string `json:"tracks,omitempty"`
}

func (s *Selection) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Selection) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	return nil
}