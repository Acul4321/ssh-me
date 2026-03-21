package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Profile struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	DisplayName  string `json:"display_name"`
	Status       string `json:"status"`
	Colour       string `json:"colour"`
	AccentColour string `json:"accent_colour"`
	Bio          string `json:"bio"`
	Pronouns     string `json:"pronouns"`
	Location     string `json:"location"`
	Layout       string `json:"layout"`
	AvatarURL    string `json:"avatar_url"`
}

type Link struct {
	Platform  string `json:"platform"`
	Label     string `json:"label"`
	URL       string `json:"url"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
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

// GetLinksByProfileId fetches visible links for a profile, sorted by sort_order.
func GetLinksByProfileId(baseURL, profileId string) ([]Link, error) {
	client := resty.New()

	var result struct {
		Items []Link `json:"items"`
	}

	resp, err := client.R().
		SetQueryParam("filter", fmt.Sprintf("profile='%s' && visible=true", profileId)).
		SetQueryParam("sort", "sort_order").
		SetResult(&result).
		Get(baseURL + "/api/collections/links/records")

	if err != nil {
		return nil, fmt.Errorf("links fetch failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("links fetch failed (status %d)", resp.StatusCode())
	}

	return result.Items, nil
}
