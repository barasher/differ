package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetInfoExisting(t *testing.T) {
	f, err := getInfo("testdata/")
	assert.Nil(t, err)
	assert.NotNil(t, f)
}

func TestGetInfoNonExisting(t *testing.T) {
	f, err := getInfo("nonExisting/")
	assert.Nil(t, err)
	assert.Nil(t, f)
}

func checkHasDiff(t *testing.T, diffs []Difference, expDiffType DiffType, expElement string) {
	for _, curDiff := range diffs {
		if curDiff.Type == expDiffType && curDiff.Element == expElement {
			return
		}
	}
	t.Errorf("Could not find difference %v, %v", expDiffType, expElement)
}

func logDiffs(t *testing.T, diffs []Difference) {
	for i, curDiff := range diffs {
		t.Logf("Difference %v: %v, %v", i, curDiff.Type, curDiff.Element)
	}
}

func TestCompareNonExistingTarget(t *testing.T) {
	diffs, err := compare("testdata/from", "folderNonExistingInTo", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 1, len(diffs), "wrong length")
	checkHasDiff(t, diffs, MISSING, "folderNonExistingInTo")
}

func TestCompareDifferentFile(t *testing.T) {
	diffs, err := compare("testdata/from", "different.txt", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 1, len(diffs), "wrong length")
	checkHasDiff(t, diffs, DIFFERENT, "different.txt")
}

func TestCompareRecursive(t *testing.T) {
	diffs, err := compare("testdata/from", "folder", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 2, len(diffs), "wrong length")
	checkHasDiff(t, diffs, DIFFERENT, "folder/different.txt")
	checkHasDiff(t, diffs, MISSING, "folder/nonExistingInTo.txt")
}

func TestDiff(t *testing.T) {
	diffs, err := Diff("testdata/from", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	checkHasDiff(t, diffs, DIFFERENT, "different.txt")
	checkHasDiff(t, diffs, MISSING, "nonExistingInTo.txt")
	checkHasDiff(t, diffs, DIFFERENT, "folder/different.txt")
	checkHasDiff(t, diffs, MISSING, "folder/nonExistingInTo.txt")
}