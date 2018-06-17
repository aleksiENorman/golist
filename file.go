package main

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	flock "github.com/theckman/go-flock"
)

// Save in file using series SHA1-sum as name
func (e *entry) save() error {
	filename := calcName(e.Series, true)

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if err := writeFile(calcName("index", false), []string{e.Time.String(), e.Series}); err != nil {
			return err
		}
	}

	return writeFile(filename, []string{e.Time.String(), e.Message})
}

func writeFile(filename string, record []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	lock := flock.NewFlock(filename)

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

	if csvWriter.Write(record) != nil {
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
	result := make([]entry, 0, 100)

	file, err := os.Open(calcName(series, series != "index"))
	if err != nil {
		fmt.Println(series)
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
		if len(record) == 2 && len(strings.TrimSpace(record[0])) != 0 {
			currentResult.Time.UnmarshalText([]byte(record[0]))
			currentResult.Series = series
			currentResult.Message = record[1]

			//fmt.Println(currentResult)
			if !currentResult.Time.IsZero() || currentResult.Message != "" {
				result = append(result, currentResult)
			}
		}
	}

	return result, nil
}

// Hash series to create name
func calcName(series string, doHash bool) string {
	var baseName []byte

	if doHash {
		hash := sha1.New()
		hash.Write([]byte(series))
		baseName = hash.Sum(nil)
	} else {
		baseName = []byte(series)
		return fmt.Sprintf("data/%s.csv", baseName)
	}
	return fmt.Sprintf("data/%x.csv", baseName)
}
