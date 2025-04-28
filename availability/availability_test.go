package availability

import (
	"autoassigner/config"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlwaysAvailable(t *testing.T) {
	checker := &AlwaysAvailable{}

	tests := []struct {
		name     string
		username string
	}{
		{
			name:     "regular user",
			username: "alice",
		},
		{
			name:     "empty username",
			username: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available, err := checker.IsAvailable(tt.username)
			if err != nil {
				t.Errorf("AlwaysAvailable.IsAvailable() error = %v", err)
				return
			}
			if !available {
				t.Error("AlwaysAvailable.IsAvailable() = false, want true")
			}
		})
	}
}

func TestInOutChecker(t *testing.T) {
	// Create a test server to mock the InOut API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Path[len("/status/"):]
		var response map[string]interface{}

		switch username {
		case "available":
			response = map[string]interface{}{
				"inOutLocation": "OFFICE",
			}
		case "away":
			response = map[string]interface{}{
				"inOutLocation": "AWAY",
			}
		case "ooo":
			response = map[string]interface{}{
				"inOutLocation": "OOO",
			}
		case "no-status":
			response = map[string]interface{}{}
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Configure the InOut API URL
	config.Settings.Availability.InOutApiUrlPrefix = server.URL + "/status/"
	config.Settings.Availability.InOutUnavailableStatuses = []string{"OOO", "AWAY"}

	checker := &InOutChecker{}

	tests := []struct {
		name        string
		username    string
		want        bool
		wantErr     bool
		errContains string
	}{
		{
			name:     "available user",
			username: "available",
			want:     true,
		},
		{
			name:     "away user",
			username: "away",
			want:     false,
		},
		{
			name:     "out of office user",
			username: "ooo",
			want:     false,
		},
		{
			name:     "user with no status",
			username: "no-status",
			want:     true,
		},
		{
			name:        "non-existent user",
			username:    "non-existent",
			want:        false,
			wantErr:     true,
			errContains: "404",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available, err := checker.IsAvailable(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("InOutChecker.IsAvailable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if available != tt.want {
				t.Errorf("InOutChecker.IsAvailable() = %v, want %v", available, tt.want)
			}
		})
	}
}

func TestCheckerInterface(t *testing.T) {
	var _ Checker = &AlwaysAvailable{} // Verify AlwaysAvailable implements Checker
	var _ Checker = &InOutChecker{}    // Verify InOutChecker implements Checker
}
