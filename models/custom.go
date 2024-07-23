package models

import (
	"encoding/json"
	"time"
)

type Date struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface.
func (ct *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Format(time.DateOnly))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ct *Date) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	t, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// MarshalText implements the text Marshaler interface.
func (ct *Date) MarshalText() ([]byte, error) {
	return []byte(ct.Format(time.DateOnly)), nil
}

// UnmarshalText implements the text Unmarshaler interface.
func (ct *Date) UnmarshalText(data []byte) error {
	t, err := time.Parse(time.DateOnly, string(data))
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

// String returns the time formatted using the custom layout.
func (ct *Date) String() string {
	return ct.Format(time.DateOnly)
}

// Parse transforms the time.DateTime string format to Date
func Parse(layout string) (*Date, error) {
	t, err := time.Parse(time.DateOnly, layout)
	if err != nil {
		return nil, err
	}

	return &Date{t}, nil
}
