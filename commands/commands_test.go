package commands

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseArguments(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		desc              string
		rawArgs           []string
		expectedArguments arguments
		expectedErr       error
	}{
		{
			desc:              "parse dataRoot successfully",
			rawArgs:           []string{"./data"},
			expectedArguments: arguments{DataRoot: "./data"},
			expectedErr:       nil,
		},
		{
			desc:              "error when no arguments given",
			rawArgs:           []string{},
			expectedArguments: arguments{},
			expectedErr:       errors.New("no arguments given"),
		},
		{
			desc:              "error when dataRoot empty",
			rawArgs:           []string{""},
			expectedArguments: arguments{},
			expectedErr:       errors.New("dataRoot is required"),
		},
		{
			desc:              "error when dataRoot whitespace only",
			rawArgs:           []string{"   "},
			expectedArguments: arguments{},
			expectedErr:       errors.New("dataRoot is required"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actualArguments, actualErr := parseArguments(tC.rawArgs)

			if tC.expectedErr == nil {
				assert.Equal(tC.expectedArguments, actualArguments)
			} else {
				assert.Equal(tC.expectedErr, actualErr)
			}
		})
	}
}
