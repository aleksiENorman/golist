package main

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	flock "github.com/theckman/go-flock"
)

// Save in file using series SHA1-sum as name
func (e *entry) save() error {
	filename := calcName(e.series)
	records := []string{e.time.String(), e.message}
	lock := flock.NewFlock(filename)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	for {
		locked, err := lock.TryLock()
		if err != nil {
			return err
		}
		if locked {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}

	if csvWriter.Write(records) != nil {
		lock.Unlock()
		return err
	}
	csvWriter.Flush()
	lock.Unlock()

	return nil
}

// Load from file. (Untested 2017-05-03)
func load(series string) ([]entry, error) {
	var currentResult entry
	result := make([]entry, 1, 100)

	file, err := os.Open(calcName(series))
	if err != nil {
		return []entry{}, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []entry{}, err
		}

		currentResult.time.UnmarshalText([]byte(record[0]))
		currentResult.series = series
		currentResult.message = record[1]

		result = append(result, currentResult)
	}

	return result, nil
}

// Hash series to create name
func calcName(series string) string {
	hash := sha1.New()
	hash.Write([]byte(series))
	return fmt.Sprintf("data/%x.csv", hash.Sum(nil))
}
