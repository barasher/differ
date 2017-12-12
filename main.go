package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const argCount = 2
const argFrom = 0
const argTo = 1

const defaultConfFile = ""

const codeNoDiff = 0
const codeDiff = 1
const codeConfError = 2
const codeExecError = 3

// Unmarshal configuration file
func unmarshalConfiguration(confFile string) (*DifferConf, error) {
	conf := DifferConf{}
	if confFile != defaultConfFile {
		confReader, err := os.Open(confFile)
		if err != nil {
			return nil, fmt.Errorf("error while opening configuration file %v: %v", confFile, err)
		}
		err = json.NewDecoder(confReader).Decode(&conf)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshaling configuration file %v : %v", confFile, err)
		}
	}
	return &conf, nil
}

// Executes diff, returns code that represents execution state
func execute(confFile string, from string, to string) int {
	var err error

	// configuration
	conf := &DifferConf{}
	if confFile != defaultConfFile {
		conf, err = unmarshalConfiguration(confFile)
		if err != nil {
			logrus.Errorf("%v", err)
			return codeConfError
		}
	}

	// execution
	differ, err := NewDiffer(*conf)
	if err != nil {
		logrus.Errorf("Error while creating a Differ: %v", err)
		return codeExecError
	}
	diffs, err := differ.Diff(from, to)
	if err != nil {
		logrus.Errorf("%v", err)
		return codeExecError
	}

	// results
	if len(diffs) > 0 {
		for _, d := range diffs {
			logrus.Warnf("[%v] %v", diffTypeLabel[d.Type], d.Element)
		}
		return codeDiff
	}

	return codeNoDiff
}

// Main function
func main() {
	logrus.SetOutput(os.Stdout)
	confFile := flag.String("c", defaultConfFile, "configuration file")
	flag.Parse()
	if len(flag.Args()) != argCount {
		logrus.Errorf("Wrong argument count, %v/%v provided", len(flag.Args()), argCount)
		os.Exit(codeConfError)
	}
	code := execute(*confFile, flag.Arg(argFrom), flag.Arg(argTo))
	os.Exit(code)
}
