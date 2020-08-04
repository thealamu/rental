package main

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	someData := struct {
		Data string
	}{
		Data: "Some String",
	}

	someDataBytes, err := json.Marshal(someData)
	if err != nil {
		t.Error(err)
	}

	rRecorder := httptest.NewRecorder()
	tag := "rjTest"

	respondJSON(tag, rRecorder, someData)

	respBodyStr := strings.TrimSpace(rRecorder.Body.String())
	someDataStr := strings.TrimSpace(string(someDataBytes))

	if respBodyStr != someDataStr {
		t.Errorf("respondJSON does not write passed in data, expected %v, got %v", someDataStr, respBodyStr)
	}
}
