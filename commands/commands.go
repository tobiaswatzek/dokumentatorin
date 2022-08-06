package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

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

	fmt.Printf("%v\n", filePaths)

	return nil
}

type arguments struct {
	DataRoot string
}

func parseArguments(rawArgs []string) (arguments, error) {
	if len(rawArgs) == 0 {
		return arguments{}, errors.New("no arguments given")
	}

	dataRoot := strings.TrimSpace(rawArgs[0])
	if dataRoot == "" {
		return arguments{}, errors.New("dataRoot is required")
	}

	return arguments{DataRoot: dataRoot}, nil
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
