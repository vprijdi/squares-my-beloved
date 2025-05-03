package api

import (
	"net/http"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	testApp := newTestApplication(t, nil)
	testMux := testApp.Mount()

	t.Run("returns status OK", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/health", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, testMux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}
