package main

import "github.com/sirupsen/logrus"

func main() {
	diffs, err := Diff("testdata/from", "testdata/to")
	if err != nil {
		logrus.Errorf("%v", err)
	}
	for _, d := range diffs {
		logrus.Warnf("[%v] %v" , diffTypeLabel[d.Type], d.Element)
	}
}