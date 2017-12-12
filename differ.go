package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"regexp"
)

// DiffType is a difference type
type DiffType int

const (
	missing DiffType = 1 + iota
	sizeDifference
	typeDifference
)

// diffTypeLabel provides a label for each DiffType
var diffTypeLabel = map[DiffType]string{
	missing:        "M",
	sizeDifference: "S",
	typeDifference: "T",
}

// Difference represents a detected difference : the type of difference and the element that differs
type Difference struct {
	Type    DiffType // Type of difference
	Element string   // Element that differs
}

// DifferConf represents the configuration for a Differ execution
type DifferConf struct {
	BlacklistedPatterns []string
}

// Differ is the "main" structure of Differ
type Differ struct {
	blacklistPatterns []*regexp.Regexp
}

// NewDiffer returns a new configured Differ and any error encountered
func NewDiffer(conf DifferConf) (*Differ, error) {
	d := Differ{}

	// blacklisted patterns
	for _, p := range conf.BlacklistedPatterns {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("error while compiling pattern %v: %v", p, err)
		}
		d.blacklistPatterns = append(d.blacklistPatterns, r)
	}

	return &d, nil
}

// Diff compares two folders (recursively) and returns what differs and any error encountered
func (d *Differ) Diff(fromFolder string, toFolder string) ([]Difference, error) {
	from, err := os.Open(fromFolder)
	if err != nil {
		return nil, fmt.Errorf("error while opening %v: %v", fromFolder, err)
	}
	defer from.Close()

	froms, err := from.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("error while listing folder %v: %v", fromFolder, err)
	}

	var diffs []Difference
	for _, from := range froms {
		logrus.Infof("Checking %v...", from.Name())
		subDiffs, err := d.compare(fromFolder, from.Name(), toFolder)
		if err != nil {
			return nil, err
		} else {
			diffs = append(diffs, subDiffs...)
		}

	}
	return diffs, nil
}

// compare compares a specific item (recursively) between two folders and returns what differs and any error
// encountered.
func (d *Differ) compare(fromRoot string, fromRelative string, toRoot string) ([]Difference, error) {
	// check blacklist
	for _, bl := range d.blacklistPatterns {
		if bl.MatchString(fromRelative) {
			return nil, nil
		}
	}

	var err error
	var diffs []Difference

	// from
	fromAbs := filepath.Join(fromRoot, fromRelative)
	var fromFile *os.File
	var fromInfo os.FileInfo
	if fromFile, err = os.Open(fromAbs); err != nil {
		return nil, err
	}
	defer fromFile.Close()
	if fromInfo, err = fromFile.Stat(); err != nil {
		return nil, err
	}

	// to
	toAbs := filepath.Join(toRoot, fromRelative)
	var toInfo os.FileInfo
	if toInfo, err = os.Stat(toAbs); err != nil {
		if os.IsNotExist(err) {
			return append(diffs, Difference{missing, fromRelative}), nil
		}
		return nil, err
	}

	if fromInfo.IsDir() && toInfo.IsDir() { // folders => recursive
		subs, err := fromFile.Readdir(-1)
		if err != nil {
			return nil, fmt.Errorf("error while reading folder %v: %v", fromAbs, err)
		}
		for _, sub := range subs {
			subDiff, err := d.compare(fromRoot, filepath.Join(fromRelative, sub.Name()), toRoot)
			if err != nil {
				return nil, err
			} else {
				diffs = append(diffs, subDiff...)
			}
		}
	} else if !fromInfo.IsDir() && !toInfo.IsDir() { // files
		if fromInfo.Size() != toInfo.Size() { // different size
			return append(diffs, Difference{sizeDifference, fromRelative}), nil
		}
	} else { // file & folder
		return append(diffs, Difference{typeDifference, fromRelative}), nil
	}

	return diffs, nil
}
