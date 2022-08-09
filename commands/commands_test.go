package commands

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"watzek.dev/apps/dokumentatorin/util"
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
			desc:              "parse args successfully",
			rawArgs:           []string{"./data", "./schema.json"},
			expectedArguments: arguments{DataRoot: "./data", SchemaPath: "./schema.json"},
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
		{
			desc:              "no error when schemaPath empty",
			rawArgs:           []string{"./data", ""},
			expectedArguments: arguments{DataRoot: "./data", SchemaPath: ""},
			expectedErr:       nil,
		},
		{
			desc:              "no error when schemaPath whitespace only",
			rawArgs:           []string{"./data", "   "},
			expectedArguments: arguments{DataRoot: "./data", SchemaPath: ""},
			expectedErr:       nil,
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

func Test_parseDataFile(t *testing.T) {
	assert := assert.New(t)
	memoryFs := afero.NewMemMapFs()
	filePath := "/data/foo.yaml"
	content := "bar: 1\nbaz: bam\n"
	afero.WriteFile(memoryFs, filePath, []byte(content), fs.FileMode(os.O_TRUNC))

	actualData, err := parseDataFile(filePath, memoryFs)

	assert.Nil(err)
	expectedData := map[string]interface{}{
		"bar": 1,
		"baz": "bam",
	}
	assert.Equal(expectedData, actualData)
}

func Test_parseDataFiles(t *testing.T) {
	assert := assert.New(t)
	memoryFs := afero.NewMemMapFs()

	expectedFiles := []struct {
		parsed  parsedDataFile
		content string
	}{
		{parsed: parsedDataFile{
			FileName: "foo",
			Data: map[string]interface{}{
				"bar": 1,
				"baz": "bam",
			}},
			content: "bar: 1\nbaz: bam\n",
		},
		{
			parsed: parsedDataFile{
				FileName: "banana",
				Data: map[string]interface{}{
					"bar": 42,
					"baz": "yeah man",
				}},
			content: "bar: 42\nbaz: yeah man\n",
		},
	}

	filePaths := make([]string, len(expectedFiles))

	for i, f := range expectedFiles {
		path := filepath.Join("/data", f.parsed.FileName+".yaml")
		err := afero.WriteFile(memoryFs, path, []byte(f.content), fs.FileMode(os.O_TRUNC))
		require.Nil(t, err)
		filePaths[i] = path
	}

	actualParsedData, err := parseDataFiles(filePaths, memoryFs)

	assert.Nil(err)
	expectedParsedData := util.Map(expectedFiles, func(d struct {
		parsed  parsedDataFile
		content string
	}) parsedDataFile {
		return d.parsed
	})
	assert.ElementsMatch(expectedParsedData, actualParsedData)
}

func createFiles(filePaths []string, appFs afero.Fs, t *testing.T) {
	for _, f := range filePaths {
		_, err := appFs.Create(f)
		require.Nil(t, err, "file must be created")
	}
}
