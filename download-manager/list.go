package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type List struct {
	Items []struct {
		URL   string `json:"url"`
		Limit string `json:"limit"`
	} `json:"items"`
}

func FromList(path string) ([]*Download, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read list file: %w", err)
	}

	var list List
	if err = json.Unmarshal(content, &list); err != nil {
		return nil, fmt.Errorf("unmarshal list: %w", err)
	}

	var downloads []*Download
	for _, item := range list.Items {
		limit, err := Parse(item.Limit)
		if err != nil {
			return nil, fmt.Errorf("parse limit: %w", err)
		}
		downloads = append(downloads, &Download{
			URL:   item.URL,
			Limit: limit,
		})
	}

	return downloads, nil
}
