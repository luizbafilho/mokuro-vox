package mokurovox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Style struct {
	ID   int
	Name string
}

type Speaker struct {
	Name   string
	Styles []Style
}

func ListSpeakers() ([]Speaker, error) {
	resp, err := http.Get("http://localhost:50021/speakers")
	if err != nil {
		return nil, fmt.Errorf("failed listing speakers: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading speakers response: %w", err)
	}

	var result []Speaker
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		return nil, fmt.Errorf("failed speakers unmarshal JSON: %w", err)
	}

	return result, nil
}
