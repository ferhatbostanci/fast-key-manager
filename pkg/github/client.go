package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	baseURL string
}

type Key struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

func NewClient() *Client {
	return &Client{
		baseURL: "https://api.github.com",
	}
}

func (c *Client) GetUserKeys(username string) ([]Key, error) {
	url := fmt.Sprintf("%s/users/%s/keys", c.baseURL, username)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys from GitHub: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var keys []Key
	if err := json.Unmarshal(body, &keys); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %v", err)
	}

	return keys, nil
} 