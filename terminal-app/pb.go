package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Profile struct {
	Username string `json:"username"`
	Status   string `json:"status"`
	Colour   string `json:"colour"`
}

// GetProfileByUsername fetches a public profile from PocketBase.
// No admin credentials required — the profiles collection is publicly readable.
func GetProfileByUsername(baseURL, username string) (*Profile, error) {
	client := resty.New()

	var result struct {
		Items []Profile `json:"items"`
	}

	resp, err := client.R().
		SetQueryParam("filter", fmt.Sprintf("username='%s'", username)).
		SetResult(&result).
		Get(baseURL + "/api/collections/profiles/records")

	if err != nil {
		return nil, fmt.Errorf("connection failed - is PocketBase running at %s? Error: %w", baseURL, err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("search failed (status %d)", resp.StatusCode())
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user '%s' not found", username)
	}

	return &result.Items[0], nil
}
