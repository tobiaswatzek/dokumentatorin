package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func Execute(args []string) error {
	parsedArgs, err := parseArguments(args)
	if err != nil {
		return err
	}

	appFs := afero.NewOsFs()

	yamlRegex, err := regexp.Compile(`.*\.(ya?ml)`)
	if err != nil {
		return err
	}

	filePaths, err := findAllMatchingFilePaths(parsedArgs.DataRoot, yamlRegex, appFs)
	if err != nil {
		return err
	}

	parsedFiles, err := parseDataFiles(filePaths, appFs)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", parsedFiles)

	return nil
}

type arguments struct {
	DataRoot   string
	SchemaPath string
}

func parseArguments(rawArgs []string) (arguments, error) {
	if len(rawArgs) == 0 {
		return arguments{}, errors.New("no arguments given")
	}

	dataRoot := strings.TrimSpace(rawArgs[0])
	if dataRoot == "" {
		return arguments{}, errors.New("dataRoot is required")
	}

	schemaPath := strings.TrimSpace(rawArgs[1])

	return arguments{DataRoot: dataRoot, SchemaPath: schemaPath}, nil
}

func findAllMatchingFilePaths(root string, pattern *regexp.Regexp, appFs afero.Fs) ([]string, error) {
	var files []string
	err := afero.Walk(appFs, root, func(path string, d fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if pattern.MatchString(d.Name()) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

type parsedDataFile struct {
	FileName string
	Data     interface{}
}

func parseDataFiles(paths []string, appFs afero.Fs) ([]parsedDataFile, error) {
	parsedFiles := make([]parsedDataFile, len(paths))

	for i, path := range paths {
		fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

		data, err := parseDataFile(path, appFs)
		if err != nil {
			return nil, err
		}

		parsedFiles[i] = parsedDataFile{FileName: fileName, Data: data}
	}

	return parsedFiles, nil
}

func parseDataFile(path string, appFs afero.Fs) (interface{}, error) {
	contents, err := afero.ReadFile(appFs, path)

	if err != nil {
		return nil, err
	}

	var data interface{}
	err = yaml.Unmarshal(contents, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func buildJsonSchema(schemaPath string, appFs afero.Fs) (*jsonschema.Schema, error) {
	schemaFile, err := appFs.Open(schemaPath)
	if err != nil {
		return nil, err
	}

	compiler := jsonschema.NewCompiler()
	schemaName := filepath.Base(schemaPath)
	err = compiler.AddResource(schemaName, schemaFile)
	if err != nil {
		return nil, err
	}
	schema, err := compiler.Compile(schemaName)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
