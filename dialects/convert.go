package dialects

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Define a converter function type
type Converter func(event *Event) (string, error)

// Returns a converter function based on file extension
func GetConverterFunction(fileFormat string) (Converter, error) {
	switch fileFormat {
	case "json":
		return ConvertToJSON, nil
	}
	return nil, fmt.Errorf("Unsupported output `%s` file format (use `json` or `csv`)", fileFormat)
}

// Dumps the Event into a JSON string
func ConvertToJSON(event *Event) (string, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(event); err != nil {
		return "", err
	}
	return b.String(), nil
}
