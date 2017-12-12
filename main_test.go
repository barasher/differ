package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func logConf(t *testing.T, c *DifferConf) {
	for i, cur := range c.BlacklistedPatterns {
		t.Logf("blacklist pattern %v: %v", i, cur)
	}
}

func checkHasBlacklistedPattern(t *testing.T, c *DifferConf, p string) {
	for _, cur := range c.BlacklistedPatterns {
		if cur == p {
			return
		}
	}
	t.Errorf("Could not find pattern %v", p)
}

func TestUnmarshalConfigurationNominal(t *testing.T) {
	c, err := unmarshalConfiguration("testdata/conf/blacklist.json")
	assert.Nil(t, err)
	assert.NotNil(t, c)
	logConf(t, c)
	checkHasBlacklistedPattern(t, c, "^.*txt$")
	checkHasBlacklistedPattern(t, c, "^.*doc$")
}

func TestUnmarshalFailOnUnmarshalingConfiguration(t *testing.T) {
	c, err := unmarshalConfiguration("testdata/conf/unmarshalable.json")
	assert.NotNil(t, err)
	assert.Nil(t, c)
}

func TestExecuteFailOnUnmarshalingConfiguration(t *testing.T) {
	var testCases = []struct {
		testCaseId string
		confFile   string
		from       string
		to         string
		expReturn  int
	}{
		{"executeFailOnMarshalingConfiguration", "testdata/conf/unmarshalable.json",
			"testdata/from", "testdata/to", codeConfError},
		{"executeFailOnMarshalingConfiguration", "testdata/conf/nonExisting.json",
			"testdata/from", "testdata/to", codeConfError},
		{"executeFailOnDifferCreation", "testdata/conf/blacklistError.json",
			"testdata/from", "testdata/to", codeExecError},
		{"executeFailOnDifferExecution", "testdata/conf/blacklist.json",
			"unknownFolder", "testdata/to", codeExecError},
		{"executeNominalWithDifferences", "testdata/conf/blacklist.json",
			"testdata/from", "testdata/to", codeDiff},
		{"executeNominalWithoutDifference", "testdata/conf/blacklist.json",
			"testdata/from", "testdata/from", codeNoDiff},
		{"executeNominalWithoutConfiguration", defaultConfFile,
			"testdata/from", "testdata/to", codeDiff},
	}

	for _, tc := range testCases {
		t.Run(tc.testCaseId, func(t *testing.T) {
			code := execute(tc.confFile, tc.from, tc.to)
			assert.Equalf(t, tc.expReturn, code, "wrong return code for %v, %v expected but got %v",
				tc.testCaseId, tc.expReturn, code)
		})
	}
}
