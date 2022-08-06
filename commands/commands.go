package commands

import (
	"errors"
	"fmt"
	"strings"
)

func Execute(args []string) error {

	parsedArgs, err := parseArguments(args)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", parsedArgs)

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
