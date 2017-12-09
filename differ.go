package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// DiffType is a difference type
type DiffType int

const (
	// Missing element
	MISSING DiffType = 1 + iota
	// Size difference
	SIZE_DIFFERENCE
	// Element that differs
	TYPE_DIFFERENCE
)

// diffTypeLabel provides a label for each DiffType
var diffTypeLabel = map[DiffType]string{
	MISSING:         "M",
	SIZE_DIFFERENCE: "S",
	TYPE_DIFFERENCE: "T",
}

// A Difference represents a detected difference : the type of difference and the element that differs
type Difference struct {
	Type    DiffType // Type of difference
	Element string   // Element that differs
}

// Diff compares two folders (recursively) and returns what differs and any error encountered.
func Diff(fromFolder string, toFolder string) ([]Difference, error) {
	from, err := os.Open(fromFolder)
	if err != nil {
		return nil, fmt.Errorf("error while opening %v: %v", fromFolder, err)
	}

	froms, err := from.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("error while listing folder %v: %v", fromFolder, err)
	}

	var diffs []Difference
	for _, from := range froms {
		subDiffs, err := compare(fromFolder, from.Name(), toFolder)
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
func compare(fromRoot string, fromRelative string, toRoot string) ([]Difference, error) {
	var err error

	// from
	fromAbs := filepath.Join(fromRoot, fromRelative)
	var fromInfo os.FileInfo
	var fromFile *os.File
	if fromFile, fromInfo, err = getInfo(fromAbs); err != nil {
		return nil, err
	}

	// to
	toAbs := filepath.Join(toRoot, fromRelative)
	var toInfo os.FileInfo
	if _, toInfo, err = getInfo(toAbs); err != nil {
		return nil, err
	}

	var diffs []Difference

	if toInfo == nil { // non existing
		return append(diffs, Difference{MISSING, fromRelative}), nil
	}

	if fromInfo.IsDir() && toInfo.IsDir() { // folders => recursive
		subs, err := fromFile.Readdir(-1)
		if err != nil {
			return nil, fmt.Errorf("error while reading folder %v: %v", fromAbs, err)
		}
		for _, sub := range subs {
			subDiff, err := compare(fromRoot, filepath.Join(fromRelative, sub.Name()), toRoot)
			if err != nil {
				return nil, err
			} else {
				diffs = append(diffs, subDiff...)
			}
		}
	} else if !fromInfo.IsDir() && !toInfo.IsDir() { // files
		if fromInfo.Size() != toInfo.Size() { // different size
			return append(diffs, Difference{SIZE_DIFFERENCE, fromRelative}), nil
		}
	} else { // file & folder
		return append(diffs, Difference{TYPE_DIFFERENCE, fromRelative}), nil
	}

	return diffs, nil

}

// getInfo returns a file metadata (or nil if it does not exist) and any error if encountered
func getInfo(file string) (*os.File, os.FileInfo, error) {
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		} else {
			return nil, nil, fmt.Errorf("error while opening %v: %v", file, err)
		}
	}
	if finfo, err := f.Stat(); err != nil {
		return nil, nil, fmt.Errorf("error while getting stats for %v: %v", file, err)
	} else {
		return f, finfo, nil
	}
}
