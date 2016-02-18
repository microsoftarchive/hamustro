package dialects

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// Returns a Converter function for a single Event conversion
func TestFunctionGetConverterFunction(t *testing.T) {
	t.Log("Returns a Converter functions for single Event conversion")

	// Testing CSV format
	csvFn, err := GetConverterFunction("csv")
	if err != nil {
		t.Errorf("Converter csv should work")
	}
	if reflect.ValueOf(ConvertCSV).Pointer() != reflect.ValueOf(csvFn).Pointer() {
		t.Errorf("Not expected csv function was returned")
	}

	// Testing JSON format
	jsonFn, err := GetConverterFunction("json")
	if err != nil {
		t.Errorf("Converter json should work")
	}
	if reflect.ValueOf(ConvertJSON).Pointer() != reflect.ValueOf(jsonFn).Pointer() {
		t.Errorf("Not expected json function was returned")
	}

	// Testing mispelled (uppercase) targets
	_, err = GetConverterFunction("CSV")
	if err == nil {
		t.Errorf("Only lowercase FileFormat is supported")
	}

	// Testing not existing converters
	_, err = GetConverterFunction("")
	if err == nil {
		t.Errorf("Empty or not existing converter should return with error")
	}
}

// Returns a BatchConverter functions for multiple Event conversion
func TestFunctionGetBatchConverterFunction(t *testing.T) {
	t.Log("Returns a BatchConverter functions for multiple Event conversion")

	// Testing CSV format
	csvFn, err := GetBatchConverterFunction("csv")
	if err != nil {
		t.Errorf("Converter csv should work")
	}
	if reflect.ValueOf(ConvertBatchCSV).Pointer() != reflect.ValueOf(csvFn).Pointer() {
		t.Errorf("Not expected csv function was returned")
	}

	// Testing JSON format
	jsonFn, err := GetBatchConverterFunction("json")
	if err != nil {
		t.Errorf("Converter json should work")
	}
	if reflect.ValueOf(ConvertBatchJSON).Pointer() != reflect.ValueOf(jsonFn).Pointer() {
		t.Errorf("Not expected json function was returned")
	}

	// Testing mispelled (uppercase) targets
	_, err = GetBatchConverterFunction("CSV")
	if err == nil {
		t.Errorf("Only lowercase FileFormat is supported")
	}

	// Testing not existing converters
	_, err = GetBatchConverterFunction("")
	if err == nil {
		t.Errorf("Empty or not existing converter should return with error")
	}
}

// Converting a single event to CSV
func TestFunctionConvertCSV(t *testing.T) {
	t.Log("Converting a single event to CSV to check")
	event := GetTestEvent(97421193)
	b, err := ConvertCSV(event)
	if err != nil {
		t.Errorf("CSV conversion failed %s", err.Error())
	}
	if exp := strings.Join(event.String(), "\001") + "\n"; b.String() != exp {
		t.Errorf("Expected generated CSV line was `%s` but it was `%s` instead", exp, b.String())
	}
}

// Converting a single event to JSON
func TestFunctionConvertJSON(t *testing.T) {
	t.Log("Converting a single event to JSON and back to check integrity")
	event := GetTestEvent(97421193)
	b, err := ConvertJSON(event)
	if err != nil {
		t.Errorf("JSON conversion failed %s", err.Error())
	}
	jsonStr := b.String()
	if jsonStr[len(jsonStr)-1:] != "\n" {
		t.Errorf("Converted JSON should contains a new line")
	}
	var copyEvent Event
	if err := json.Unmarshal([]byte(jsonStr), &copyEvent); err != nil {
		t.Errorf("JSON conversion failed %s", err.Error())
	}
	if !reflect.DeepEqual(copyEvent, *event) {
		t.Errorf("JSON conversion is creating invalid JSON")
	}
}

// Converting multiple events to CSV
func TestFunctionConvertBatchCSV(t *testing.T) {
	t.Log("Converting a 3 events to CSV to check")

	events := []*Event{GetTestEvent(97421193), GetTestEvent(197421199), GetTestEvent(7421191)}
	b, err := ConvertBatchCSV(events)
	if err != nil {
		t.Errorf("Batch CSV conversion is failed: %s", err.Error())
	}
	for i := range events {
		line, _ := b.ReadString('\n')
		sb, _ := ConvertCSV(events[i])
		if sb.String() != line {
			t.Errorf("Line mismatch, expected was `%s` but it was `%s`", sb, line)
		}
	}
}

// Converting multiple events to JSON
func TestFunctionConvertBatchJSON(t *testing.T) {
	t.Log("Converting a 3 events to JSON to check")

	events := []*Event{GetTestEvent(97421193), GetTestEvent(197421199), GetTestEvent(7421191)}
	b, err := ConvertBatchJSON(events)
	if err != nil {
		t.Errorf("Batch JSON conversion is failed: %s", err.Error())
	}
	for i := range events {
		line, _ := b.ReadString('\n')
		sb, _ := ConvertJSON(events[i])
		if sb.String() != line {
			t.Errorf("Line mismatch, expected was `%s` but it was `%s`", sb, line)
		}
	}
}
