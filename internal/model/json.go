package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON is a JSON object.
type JSON map[string]interface{}

// Value implements driver.Valuer interface for JSON.
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}

	return json.Marshal(j)
}

// Scan implements sql.Scanner interface for JSON.
func (j JSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if data, ok := value.([]byte); ok {
		return json.Unmarshal(data, &j)
	}

	return fmt.Errorf("could not not decode type %T -> %T", value, j)
}
