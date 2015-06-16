package util

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &nt.Time)
	if err != nil {
		return err
	}
	nt.Valid = true
	return nil
}

func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

const JSONTimeLayout string = "\"2006-01-02T15:04:05.000Z\""

type JSONTime time.Time

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	tm, err := time.Parse(JSONTimeLayout, string(data))
	if err != nil {
		return err
	}
	*t = JSONTime(tm)
	return nil
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(JSONTimeLayout)), nil
}

func (t JSONTime) String() string {
	return time.Time(t).String()
}
