package availability

import (
	"autoassigner/config"
	"encoding/json"
	"net/http"
)

type InOutChecker struct{}

func (c *InOutChecker) IsAvailable(username string) (bool, error) {
	url := config.Settings.Availability.InOutApiUrlPrefix + username
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, err
	}

	status, ok := result["inOutLocation"].(string)
	if !ok {
		return true, nil
	}

	for _, unavailable := range config.Settings.Availability.InOutUnavailableStatuses {
		if status == unavailable {
			return false, nil
		}
	}
	return true, nil
}
