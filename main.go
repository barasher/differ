package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

const ARG_COUNT = 2
const ARG_FROM = 0
const ARG_TO = 1

const DEFAULT_CONFFILE = ""

const CODE_NO_DIFF = 0
const CODE_DIFF = 1
const CODE_CONF_ERROR = 2
const CODE_EXEC_ERROR = 3

func unmarshallConfiguration(confFile string) (*DifferConf, error) {
	// configuration
	conf := DifferConf{}
	if confFile != DEFAULT_CONFFILE {
		confReader, err := os.Open(confFile)
		if err != nil {
			return nil, fmt.Errorf("error while opening configuration file %v: %v", confFile, err)
		}
		err = json.NewDecoder(confReader).Decode(&conf)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshaling configuration file %v : %v", confFile, err)
			os.Exit(CODE_CONF_ERROR)
		}
	}
	return &conf, nil
}

func execute(confFile string, from string, to string) int {
	var err error

	// configuration
	conf := &DifferConf{}
	if confFile != DEFAULT_CONFFILE {
		conf, err = unmarshallConfiguration(confFile)
		if err != nil {
			logrus.Errorf("%v", err)
			return CODE_CONF_ERROR
		}
	}

	// execution
	differ, err := NewDiffer(*conf)
	if err != nil {
		logrus.Errorf("Error while creating a Differ: %v", err)
		return CODE_EXEC_ERROR
	}
	diffs, err := differ.Diff(flag.Arg(ARG_FROM), flag.Arg(ARG_TO))
	if err != nil {
		logrus.Errorf("%v", err)
		return CODE_EXEC_ERROR
	}

	// results
	if len(diffs) > 0 {
		for _, d := range diffs {
			logrus.Warnf("[%v] %v", diffTypeLabel[d.Type], d.Element)
		}
		return CODE_DIFF
	}

	return CODE_NO_DIFF
}

func main() {
	confFile := flag.String("c", DEFAULT_CONFFILE, "configuration file")
	flag.Parse()

	// from & to check
	if len(flag.Args()) != ARG_COUNT {
		logrus.Errorf("Wrong argument count, %v/%v provided", len(flag.Args()), ARG_COUNT)
		os.Exit(CODE_CONF_ERROR)
	}

	code := execute(*confFile, flag.Arg(ARG_FROM), flag.Arg(ARG_TO))
	os.Exit(code)
}
