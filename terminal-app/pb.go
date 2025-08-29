package main

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// user schema
type User struct {
	Username string `json:"username"`
	Status   string `json:"status"`
	Colour   string `json:"colour"`
}

// querys db for user (by username) and returns a (pointer to a)user
func GetUserByUsername(baseURL, adminEmail, adminPassword, username string) (*User, error) {
	client := resty.New()

	// admin auth
	var authResp struct {
		Token string `json:"token"`
	}

	resp, err := client.R().
		SetBody(map[string]string{"identity": adminEmail, "password": adminPassword}).
		SetResult(&authResp).
		Post(baseURL + "/api/collections/_superusers/auth-with-password")

	if err != nil {
		return nil, fmt.Errorf("connection failed - is PocketBase running at %s? Error: %w", baseURL, err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("admin auth failed (status %d) - check credentials", resp.StatusCode())
	}

	// get user
	var result struct {
		Items []User `json:"items"`
	}
	resp, err = client.R().
		SetAuthToken(authResp.Token).
		SetQueryParam("filter", fmt.Sprintf("username='%s'", username)).
		SetResult(&result).
		Get(baseURL + "/api/collections/users/records")

	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("search failed (status %d)", resp.StatusCode())
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user '%s' not found", username)
	}

	if len(result.Items) > 1 {
		return nil, fmt.Errorf("multiple users found with username '%s'", username)
	}

	return &result.Items[0], nil
}
