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
	list(writer, nil) //TODO:Question, что за list?
	response := writer.Result()
	if response.StatusCode != http.StatusOK {
		test.Errorf("Status code is wrong. Have %d, want: %d", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		test.Fatal(err)
	}

	items := make([]Order, 10)
	if err = json.Unmarshal(jsonString, &items); err != nil {
		test.Errorf("Can't parse json response with error %v", err)
	}
}
