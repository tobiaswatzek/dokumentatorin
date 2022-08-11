package commands

import (
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

func Test_readTemplate(t *testing.T) {
	testCases := []struct {
		desc            string
		templateContent string
		expectErr       bool
	}{
		{
			desc:            "successfully read template",
			templateContent: "Hello {{ .name }}",
			expectErr:       false,
		},
		{
			desc:            "fail for non existent file",
			templateContent: "",
			expectErr:       true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert := assert.New(t)
			appFs := afero.NewMemMapFs()
			path := "/template.txt"
			if tC.templateContent != "" {
				afero.WriteFile(appFs, path, []byte(tC.templateContent), fs.FileMode(os.O_TRUNC))
			}

			tmpl, err := readTemplate(path, appFs)

			if tC.expectErr {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
				assert.NotNil(tmpl)
			}
		})
	}
}

func createFiles(filePaths []string, appFs afero.Fs, t *testing.T) {
	for _, f := range filePaths {
		_, err := appFs.Create(f)
		require.Nil(t, err, "file must be created")
	}
}
