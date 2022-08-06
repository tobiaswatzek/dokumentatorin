package commands

import (
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_findAllMatchingFilePaths(t *testing.T) {
	assert := assert.New(t)
	memoryFs := afero.NewMemMapFs()
	dataDir := "/data"
	nestedDataDir := dataDir + "/foo"
	memoryFs.MkdirAll(nestedDataDir, os.ModeDir)
	expectedFilePaths := []string{
		nestedDataDir + "/bar.txt",
		nestedDataDir + "/baz.txt",
		dataDir + "/banana.txt",
	}
	createFiles(expectedFilePaths, memoryFs, t)

	additionalFilePaths := []string{
		dataDir + "/img.jpg",
		nestedDataDir + "hello.go",
	}
	createFiles(additionalFilePaths, memoryFs, t)
	pathRegex := regexp.MustCompile(`.*\.txt`)

	actualFilePaths, err := findAllMatchingFilePaths(dataDir, pathRegex, memoryFs)

	assert.Nil(err)
	assert.ElementsMatch(expectedFilePaths, actualFilePaths)
}

func createFiles(filePaths []string, appFs afero.Fs, t *testing.T) {
	for _, f := range filePaths {
		_, err := appFs.Create(f)
		require.Nil(t, err, "file must be created")
	}
}
