package configparser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"utils"
)

type InputConfig struct {
	Filename      string
	IncludeFields []string
}

func (config *InputConfig) SetItem(key, value string) error {
	if key == "file" {
		config.Filename = value
	} else if key == "include_fields" {
		fields, err := parseConfigArrayItems(value)
		if err != nil {
			return err
		}
		config.IncludeFields = fields
	}
	return nil
}
func (config *InputConfig) ParseConfig(data string) (err error) {
	input, err := findConfigRootData(data, "input")
	if err != nil {
		return err
	}
	lines := utils.Map(splitStringIntoLines(input), func(s string) string { return strings.TrimSpace(s) })
	for i := range lines {
		key, value, err := parseConfigLine(lines[i])
		if err != nil {
			return err
		}
		err = config.SetItem(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

type OutputConfig struct {
	StatFields  []string
	StatsFields []string
}

func (config *OutputConfig) SetItem(key, value string) error {
	if key == "stat_fields" {
		fields, err := parseConfigArrayItems(value)
		if err != nil {
			return err
		}
		config.StatFields = fields
	} else if key == "stats_fields" {
		fields, err := parseConfigArrayItems(value)
		if err != nil {
			return err
		}
		config.StatsFields = fields
	}
	return nil
}
func (config *OutputConfig) ParseConfig(data string) (err error) {
	output, err := findConfigRootData(data, "output")
	if err != nil {
		return err
	}
	lines := utils.Map(splitStringIntoLines(output), func(s string) string { return strings.TrimSpace(s) })
	for i := range lines {
		key, value, err := parseConfigLine(lines[i])
		if err != nil {
			return err
		}
		err = config.SetItem(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	Input  InputConfig
	Output OutputConfig
}

func (config *Config) readConfig(file *os.File) (string, error) {
	data, err := ioutil.ReadAll(file)
	return string(data), err
}
func (config *Config) ParseConfig(file *os.File) (err error) {
	data, err := config.readConfig(file)
	if err != nil {
		return err
	}

	err = config.Input.ParseConfig(data)
	if err != nil {
		return err
	}
	err = config.Output.ParseConfig(data)
	if err != nil {
		return err
	}
	return nil
}

func splitStringIntoLines(data string) (lines []string) {
	return strings.Split(data, "\n")
}

func parseConfigLine(line string) (key, value string, err error) {
	sepIndex := strings.Index(line, ":")
	if sepIndex == -1 {
		return "", "", errors.New(fmt.Sprintf("cannot find separating sign in line: '%s'", line))
	}
	key = strings.TrimSpace(line[:sepIndex])
	value = strings.TrimSpace(line[sepIndex+1:])
	return key, value, nil
}

func parseConfigArrayItems(value string) ([]string, error) {
	items, err := findDataBetween(value, "[", "]")
	if err != nil {
		return nil, nil
	}
	return utils.Map(strings.Split(items, ","), strings.TrimSpace), nil
}

func findDataBetween(s, start, end string) (value string, err error) {
	startIndex := strings.Index(s, start)
	if startIndex == -1 {
		return value, errors.New(fmt.Sprintf("cannot find '%s' in string: '%s'", start, s))
	}
	endIndex := strings.Index(s, end)
	if endIndex == -1 {
		return value, errors.New(fmt.Sprintf("cannot find '%s' in string: '%s'", end, s))
	}
	return s[startIndex+1 : endIndex], nil
}

func findConfigRootData(data, rootName string) (value string, err error) {
	rootIndex := strings.Index(data, rootName)

	root := data[rootIndex:]
	if rootIndex == -1 {
		return value, errors.New(fmt.Sprintf("cannot find index of '%s' in config", rootName))
	}

	value, err = findDataBetween(root, "{", "}")
	if err != nil {
		return value, err
	}
	return strings.TrimSpace(value), nil
}
