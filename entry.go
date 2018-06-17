package main

import (
	"errors"
	"strings"
	"time"
)

// List entry type
type entry struct {
	Time    time.Time
	Series  string
	Message string
}

// Validate entry and default optional fields
func (e *entry) saveValidate() error {
	errs := make([]string, 0)

	if e.Time.IsZero() {
		e.Time = time.Now()
	}

	if e.Series == "" {
		errs = append(errs, "Field Series is a required field")
	}

	if e.Message == "" {
		errs = append(errs, "Field Message is a required field")
	}

	if len(errs) > 0 {
		return errors.New("One or more errors where encountered:\n" + strings.Join(errs, "\n"))
	}

	return nil
}

func (e *entry) readValidate() error {
	errs := make([]string, 0)

	if e.Time.IsZero() {
		errs = append(errs, "Field Time is a required field")
	}

	if e.Series == "" {
		errs = append(errs, "Field Series is a required field")
	}

	if e.Message == "" {
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
		e.Time.UnmarshalText(value)
	case "series":
		e.Series = string(value)
	case "message":
		e.Message = string(value)
	default:
		return errors.New("Form field " + index + " is not regonized")
	}

	return nil
}
