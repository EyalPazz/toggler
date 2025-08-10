package toggl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.track.toggl.com/api/v9"

type Client struct {
	HTTPClient *http.Client
	APIToken   string
}

type TimeEntry struct {
	ID          int    `json:"id,omitempty"`
	Description string `json:"description"`
	Start       string `json:"start,omitempty"`
	Stop        string `json:"stop,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	ProjectID   int    `json:"project_id,omitempty"`
	WorkspaceID int    `json:"workspace_id"`
	CreatedWith string `json:"created_with"`
}

type StartRequest struct {
	Description string `json:"description"`
	CreatedWith string `json:"created_with"`
	WorkspaceID int    `json:"workspace_id"`
}

type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewClient(apiToken string) *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		APIToken:   apiToken,
	}
}

func (c *Client) makeRequest(method, endpoint string, body []byte) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", baseURL, endpoint)
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	
	req.SetBasicAuth(c.APIToken, "api_token")
	req.Header.Set("Content-Type", "application/json")
	
	return c.HTTPClient.Do(req)
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	resp, err := c.makeRequest("GET", "/workspaces", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}
	
	var workspaces []Workspace
	if err := json.NewDecoder(resp.Body).Decode(&workspaces); err != nil {
		return nil, err
	}
	
	return workspaces, nil
}

func (c *Client) StartTimer(workspaceID int, description string) (*TimeEntry, error) {
	entry := map[string]interface{}{
		"description":  description,
		"workspace_id": workspaceID,
		"created_with": "toggler",
		"start":        time.Now().UTC().Format(time.RFC3339),
		"duration":     -1,
	}
	
	body, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}
	
	resp, err := c.makeRequest("POST", fmt.Sprintf("/workspaces/%d/time_entries", workspaceID), body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}
	
	var timeEntry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntry); err != nil {
		return nil, err
	}
	
	return &timeEntry, nil
}

func (c *Client) StopTimer(workspaceID int) (*TimeEntry, error) {
	current, err := c.GetCurrentTimer(workspaceID)
	if err != nil {
		return nil, err
	}
	
	if current == nil {
		return nil, nil
	}
	
	resp, err := c.makeRequest("PATCH", fmt.Sprintf("/workspaces/%d/time_entries/%d/stop", workspaceID, current.ID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}
	
	var timeEntry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntry); err != nil {
		return nil, err
	}
	
	return &timeEntry, nil
}

func (c *Client) GetCurrentTimer(workspaceID int) (*TimeEntry, error) {
	resp, err := c.makeRequest("GET", "/me/time_entries/current", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if string(bodyBytes) == "null" {
			return nil, nil
		}
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}
	
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if string(bodyBytes) == "null" {
		return nil, nil
	}
	
	var timeEntry TimeEntry
	if err := json.Unmarshal(bodyBytes, &timeEntry); err != nil {
		return nil, err
	}
	
	if timeEntry.ID == 0 {
		return nil, nil
	}
	
	return &timeEntry, nil
}

func (c *Client) GetTimeEntries(workspaceID int, startDate, endDate string) ([]TimeEntry, error) {
	endpoint := fmt.Sprintf("/me/time_entries?start_date=%s&end_date=%s", startDate, endDate)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}
	
	var timeEntries []TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntries); err != nil {
		return nil, err
	}
	
	return timeEntries, nil
}

func StartTimer(apiToken, description string) error {
	client := NewClient(apiToken)
	
	workspaces, err := client.GetWorkspaces()
	if err != nil {
		return err
	}
	
	if len(workspaces) == 0 {
		return fmt.Errorf("no workspaces found")
	}
	
	_, err = client.StartTimer(workspaces[0].ID, description)
	return err
}
