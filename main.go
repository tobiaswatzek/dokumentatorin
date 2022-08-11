package main

import (
	"flag"

	"watzek.dev/apps/dokumentatorin/commands"
)

func main() {
	args, err := parseArguments()
	if err != nil {
		panic(err)
	}

	err = commands.Execute(args)
	if err != nil {
		panic(err)
	}
}

func parseArguments() (commands.Arguments, error) {
	dataRoot := flag.String("dataRoot", "", "Directory that contains data files.")
	schemaPath := flag.String("schema", "", "Optional JSON schema used to validate data.")
	templatePath := flag.String("template", "", "Template that is used to render data.")

	flag.Parse()

	return commands.NewArguments(*dataRoot, *schemaPath, *templatePath)
}
