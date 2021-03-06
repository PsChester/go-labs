package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

//Answer тесты чтоб запустить: go test go-projects
func TestList(test *testing.T) {
	writer := httptest.NewRecorder()
	list(writer, nil)
	response := writer.Result()
	if response.StatusCode != http.StatusOK {
		test.Errorf("Status code is wrong. Have %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonString, error := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if error != nil {
		t.Fatal(error)
	}

	items := make([]Order, 10)
	if error = json.Unmarshal(jsonString, &items); error != nil {
		t.Errorf("Can't parse json response with error %v", error)
	}
}
