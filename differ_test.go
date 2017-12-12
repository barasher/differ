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
	differ, err := NewDiffer(DifferConf{})
	assert.Nil(t, err)
	diffs, err := differ.compare("testdata/from", "folderNonExistingInTo", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 1, len(diffs), "wrong length")
	checkHasDiff(t, diffs, missing, "folderNonExistingInTo")
}

func TestCompareSizeDifference(t *testing.T) {
	differ, err := NewDiffer(DifferConf{})
	assert.Nil(t, err)
	diffs, err := differ.compare("testdata/from", "different.txt", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 1, len(diffs), "wrong length")
	checkHasDiff(t, diffs, sizeDifference, "different.txt")
}

func TestCompareRecursive(t *testing.T) {
	differ, err := NewDiffer(DifferConf{})
	assert.Nil(t, err)
	diffs, err := differ.compare("testdata/from", "folder", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 2, len(diffs), "wrong length")
	checkHasDiff(t, diffs, sizeDifference, "folder/different.txt")
	checkHasDiff(t, diffs, missing, "folder/nonExistingInTo.txt")
}

func TestCompareTypeDifference(t *testing.T) {
	differ, err := NewDiffer(DifferConf{})
	assert.Nil(t, err)
	diffs, err := differ.compare("testdata/from", "differentType", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 1, len(diffs), "wrong length")
	checkHasDiff(t, diffs, typeDifference, "differentType")
}

func TestCompareBlacklisted(t *testing.T) {
	conf := DifferConf{
		BlacklistedPatterns: []string{"^.*\\.txt$"},
	}
	differ, err := NewDiffer(conf)
	assert.Nil(t, err)
	diffs, err := differ.compare("testdata/from", "different.txt", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Nil(t, diffs)
}

func TestDiff(t *testing.T) {
	differ, err := NewDiffer(DifferConf{})
	assert.Nil(t, err)
	diffs, err := differ.Diff("testdata/from", "testdata/to")
	assert.Nil(t, err)
	logDiffs(t, diffs)
	assert.Equal(t, 6, len(diffs), "wrong length")
	checkHasDiff(t, diffs, sizeDifference, "different.txt")
	checkHasDiff(t, diffs, missing, "nonExistingInTo.txt")
	checkHasDiff(t, diffs, typeDifference, "differentType")
	checkHasDiff(t, diffs, missing, "folderNonExistingInTo")
	checkHasDiff(t, diffs, sizeDifference, "folder/different.txt")
	checkHasDiff(t, diffs, missing, "folder/nonExistingInTo.txt")
}

func TestNewDifferFailOnRegexpCompile(t *testing.T) {
	conf := DifferConf{}
	conf.BlacklistedPatterns = []string{"("}
	_, err := NewDiffer(conf)
	assert.NotNil(t, err, "should have failed")
}
