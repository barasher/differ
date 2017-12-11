package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func logConf(t *testing.T, c *DifferConf) {
	for i, cur := range c.BlacklistPatterns {
		t.Logf("blacklist pattern %v: %v", i, cur)
	}
}

func checkHasBlacklistedPattern(t *testing.T, c *DifferConf, p string) {
	for _, cur := range c.BlacklistPatterns {
		if cur == p {
			return
		}
	}
	t.Errorf("Could not find pattern %v", p)
}

func TestUnmarshalConfigurationNominal(t *testing.T) {
	c, err := unmarshallConfiguration("testdata/conf/blacklist.json")
	assert.Nil(t, err)
	assert.NotNil(t, c)
	logConf(t, c)
	checkHasBlacklistedPattern(t, c, "^.*txt$")
	checkHasBlacklistedPattern(t, c, "^.*doc$")
}

func TestExecuteFailOnUnmarshalingConfiguration(t *testing.T) {
	t.Skip()
}

func TestExecuteFailOnDifferCreation(t *testing.T) {
	t.Skip()
}

func TestExecuteFailOnDifferExecution(t *testing.T) {
	t.Skip()
}

func TestExecuteNominalWithoutDifference(t *testing.T) {
	t.Skip()
}

func TestExecuteNominalWithoutDifferences(t *testing.T) {
	t.Skip()
}
