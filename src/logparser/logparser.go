package logparser

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"units"
)

func getReader(file *os.File) (reader *bufio.Reader, err error) {
	reader = bufio.NewReader(file)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func convertLineToUnit(line []byte) (units.Unit, error) {
	var unit units.Unit
	err := json.Unmarshal(line, &unit)
	if err != nil {
		return unit, err
	}
	return unit, nil
}

func ParseFile(file *os.File) (units []units.Unit, err error) {
	reader, err := getReader(file)
	if err != nil {
		return nil, err
	}

	var line []byte
	for ; err == nil; line, _, err = reader.ReadLine() {
		if len(line) < 1 {
			continue
		}
		unit, err := convertLineToUnit(line)
		if err != nil {
			log.Fatal("Cannot parse log line")
		}
		units = append(units, unit)
	}
	return units, nil
}
