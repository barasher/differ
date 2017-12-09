package main

import (
	"flag"
	"github.com/sirupsen/logrus"
)

const ARG_COUNT = 2
const ARG_FROM = 0
const ARG_TO = 1

func main() {
	flag.Parse()

	if len(flag.Args()) != ARG_COUNT {
		logrus.Errorf("Wrong argument count, %v/%v provided", len(flag.Args()), ARG_COUNT)
	}

	diffs, err := Diff(flag.Arg(ARG_FROM), flag.Arg(ARG_TO))
	if err != nil {
		logrus.Errorf("%v", err)
	}
	for _, d := range diffs {
		logrus.Warnf("[%v] %v", diffTypeLabel[d.Type], d.Element)
	}
}
