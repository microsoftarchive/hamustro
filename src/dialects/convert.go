package dialects

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
)

// Define converter functions based on the file extension type.
type Converter func(event *Event) (*bytes.Buffer, error)
type BatchConverter func(events []*Event) (*bytes.Buffer, error)

// Returns a single event converter function based on file extension.
func GetConverterFunction(fileFormat string) (Converter, error) {
	switch fileFormat {
	case "json":
		return ConvertJSON, nil
	case "csv":
		return ConvertCSV, nil
	}
	return nil, fmt.Errorf("Unsupported output `%s` file format (use `json` or `csv`)", fileFormat)
}

// Returns a batch event converter function based on file extension.
func GetBatchConverterFunction(fileFormat string) (BatchConverter, error) {
	switch fileFormat {
	case "json":
		return ConvertBatchJSON, nil
	case "csv":
		return ConvertBatchCSV, nil
	}
	return nil, fmt.Errorf("Unsupported output `%s` file format (use `json` or `csv`)", fileFormat)
}

// Dumps the Event into a JSON string
func ConvertJSON(event *Event) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(event); err != nil {
		return b, err
	}
	return b, nil
}

// Dumps the Event into a CSV string
func ConvertCSV(event *Event) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)
	writer.Comma = '\001'
	writer.Write(event.String())
	writer.Flush()
	return b, nil
}

// Converts multiple events into JSON string
func ConvertBatchJSON(events []*Event) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	for _, event := range events {
		buffer, err := ConvertJSON(event)
		if err != nil {
			return b, err
		}
		b.WriteString(buffer.String())
	}
	return b, nil
}

// Converts to CSV string for list of events
func ConvertBatchCSV(events []*Event) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	writer := csv.NewWriter(b)
	writer.Comma = '\001'
	for _, event := range events {
		writer.Write(event.String())
	}
	writer.Flush()
	return b, nil
}
