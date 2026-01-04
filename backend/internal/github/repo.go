package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Repo struct {
	ID    int64 `json:"id"`
	Owner struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	FullName        string `json:"full_name"`
	HTMLURL         string `json:"html_url"`
	Homepage        string `json:"homepage"`
	Private         bool   `json:"private"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	Description     string `json:"description"`
	Permissions struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permissions"`
}

func (c *Client) GetRepo(ctx context.Context, accessToken string, fullName string) (Repo, error) {
	// fullName is owner/repo.
	owner, repo, err := splitFullName(fullName)
	if err != nil {
		return Repo{}, err
	}
	u := "https://api.github.com/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return Repo{}, err
	}
	if strings.TrimSpace(accessToken) != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Repo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Repo{}, fmt.Errorf("github repo fetch failed: status %d", resp.StatusCode)
	}

	var r Repo
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return Repo{}, err
	}
	if r.ID == 0 || r.FullName == "" {
		return Repo{}, fmt.Errorf("invalid github repo response")
	}
	return r, nil
}

func (c *Client) GetRepoLanguages(ctx context.Context, accessToken string, fullName string) (map[string]int64, error) {
	owner, repo, err := splitFullName(fullName)
	if err != nil {
		return nil, err
	}
	u := "https://api.github.com/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/languages"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(accessToken) != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("github repo languages fetch failed: status %d", resp.StatusCode)
	}

	var langs map[string]int64
	if err := json.NewDecoder(resp.Body).Decode(&langs); err != nil {
		return nil, err
	}
	if langs == nil {
		langs = map[string]int64{}
	}
	return langs, nil
}

// ReadmeResponse represents the GitHub API response for README content
type ReadmeResponse struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"` // Base64 encoded
	Encoding string `json:"encoding"`
}

// GetReadme fetches the README.md content from a GitHub repository
func (c *Client) GetReadme(ctx context.Context, accessToken string, fullName string) (string, error) {
	owner, repo, err := splitFullName(fullName)
	if err != nil {
		return "", err
	}
	// GitHub API endpoint for README (automatically finds README.md, README, etc.)
	u := "https://api.github.com/repos/" + url.PathEscape(owner) + "/" + url.PathEscape(repo) + "/readme"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(accessToken) != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("readme not found: status %d", resp.StatusCode)
	}

	var readme ReadmeResponse
	if err := json.NewDecoder(resp.Body).Decode(&readme); err != nil {
		return "", err
	}

	// Decode base64 content
	if readme.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(readme.Content)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}
	// If not base64, return as-is (shouldn't happen with GitHub API)
	return readme.Content, nil
}

func splitFullName(fullName string) (string, string, error) {
	s := strings.TrimSpace(fullName)
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repo full name (expected owner/repo)")
	}
	owner := strings.TrimSpace(parts[0])
	repo := strings.TrimSpace(parts[1])
	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("invalid repo full name (expected owner/repo)")
	}
	return owner, repo, nil
}


