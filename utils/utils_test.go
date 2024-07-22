package utils

import (
	"encoding/json"
	"errors"
	"gitbeam/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteHTTPError(t *testing.T) {
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	err := errors.New("test error")

	// Call the function with the ResponseRecorder and the error.
	WriteHTTPError(rr, http.StatusInternalServerError, err)

	// Check the status code.
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Parse the response body.
	var result models.Result
	_ = json.NewDecoder(rr.Body).Decode(&result)

	// Check the response body.
	assert.False(t, result.Success)
	assert.Equal(t, "test error", result.Message)
	assert.Nil(t, result.Data)
}

func TestWriteHTTPSuccess(t *testing.T) {
	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	data := map[string]any{"key": "value"}
	message := "Success"

	// Call the function with the ResponseRecorder, message, and data.
	WriteHTTPSuccess(rr, message, data)

	// Check the status code.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse the response body.
	var result models.Result
	_ = json.NewDecoder(rr.Body).Decode(&result)

	// Check the response body.
	assert.True(t, result.Success)
	assert.Equal(t, message, result.Message)
	assert.EqualValues(t, data, result.Data)
}

// Define a test struct
type TestStruct struct {
	Name string
	Age  int
}

func TestUnPack(t *testing.T) {
	// Test case 1: Input as []byte
	t.Run("input as []byte", func(t *testing.T) {
		input := []byte(`{"Name":"Alice","Age":30}`)
		var target TestStruct

		err := UnPack(input, &target)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", target.Name)
		assert.Equal(t, 30, target.Age)
	})

	// Test case 2: Input as struct
	t.Run("input as struct", func(t *testing.T) {
		input := TestStruct{Name: "Bob", Age: 25}
		var target TestStruct

		err := UnPack(input, &target)
		assert.NoError(t, err)
		assert.Equal(t, "Bob", target.Name)
		assert.Equal(t, 25, target.Age)
	})

	// Test case 3: Input as map
	t.Run("input as map", func(t *testing.T) {
		input := map[string]interface{}{
			"Name": "Charlie",
			"Age":  20,
		}
		var target TestStruct

		err := UnPack(input, &target)
		assert.NoError(t, err)
		assert.Equal(t, "Charlie", target.Name)
		assert.Equal(t, 20, target.Age)
	})

	// Test case 4: Invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		input := []byte(`{"Name":"Alice","Age":}`) // Invalid JSON
		var target TestStruct

		err := UnPack(input, &target)
		assert.Error(t, err)
	})

	// Test case 5: Unmarshal into incompatible type
	t.Run("unmarshal into incompatible type", func(t *testing.T) {
		input := []byte(`{"Name":"Alice","Age":"thirty"}`) // Age is a string
		var target TestStruct

		err := UnPack(input, &target)
		assert.Error(t, err)
	})
}
