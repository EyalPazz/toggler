package toggl

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TimeEntry struct {
	Description string `json:"description"`
	CreatedWith string `json:"created_with"`
}

type StartRequest struct {
	Entry TimeEntry `json:"time_entry"`
}

func StartTimer(apiToken, description string) error {
	entry := StartRequest{
		Entry: TimeEntry{
			Description: description,
			CreatedWith: "toggler",
		},
	}
	body, _ := json.Marshal(entry)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.track.toggl.com/api/v8/time_entries/start", bytes.NewReader(body))
	req.SetBasicAuth(apiToken, "api_token")
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
