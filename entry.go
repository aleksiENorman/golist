package main

import (
	"errors"
	"strings"
	"time"
)

// List entry type
type entry struct {
	time    time.Time
	series  string
	message string
}

// Validate entry and default optional fields
func (e *entry) saveValidate() error {
	errs := make([]string, 0)

	if e.time.IsZero() {
		e.time = time.Now()
	}

	if e.series == "" {
		errs = append(errs, "Field Series is a required field")
	}

	if e.message == "" {
		errs = append(errs, "Field Message is a required field")
	}

	if len(errs) > 0 {
		return errors.New("One or more errors where encountered:\n" + strings.Join(errs, "\n"))
	}

	return nil
}

// Set index based on POST values or csv
func (e *entry) setIndex(index string, value []byte) error {
	switch index {
	case "time":
		e.time.UnmarshalText(value)
	case "series":
		e.series = string(value)
	case "message":
		e.message = string(value)
	default:
		return errors.New("Form field " + index + " is not regonized")
	}

	return nil
}
