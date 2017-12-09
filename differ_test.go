package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestCompareSizeDifference(t *testing.T) {
	diffs, err := compare("testdata/from", "different.txt", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 1, len(diffs), "wrong length")
	checkHasDiff(t, diffs, SIZE_DIFFERENCE, "different.txt")
}

func TestCompareRecursive(t *testing.T) {
	diffs, err := compare("testdata/from", "folder", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 2, len(diffs), "wrong length")
	checkHasDiff(t, diffs, SIZE_DIFFERENCE, "folder/different.txt")
	checkHasDiff(t, diffs, MISSING, "folder/nonExistingInTo.txt")
}

func TestCompareTypeDifference(t *testing.T) {
	diffs, err := compare("testdata/from", "differentType", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 1, len(diffs), "wrong length")
	checkHasDiff(t, diffs, TYPE_DIFFERENCE, "differentType")
}

func TestDiff(t *testing.T) {
	diffs, err := Diff("testdata/from", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 6, len(diffs), "wrong length")
	checkHasDiff(t, diffs, SIZE_DIFFERENCE, "different.txt")
	checkHasDiff(t, diffs, MISSING, "nonExistingInTo.txt")
	checkHasDiff(t, diffs, TYPE_DIFFERENCE, "differentType")
	checkHasDiff(t, diffs, MISSING, "folderNonExistingInTo")
	checkHasDiff(t, diffs, SIZE_DIFFERENCE, "folder/different.txt")
	checkHasDiff(t, diffs, MISSING, "folder/nonExistingInTo.txt")
}
