package utils

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
)

var (
	terraformCommand func(string) error = func(folder string) error {
		_, err := exec.Command("terraform", "fmt", folder).Output()
		return err
	}

	osReadFile   = os.ReadFile
	formatSource = format.Source
	osWriteFile  = os.WriteFile
)

func TerraformFormat(folder string) error {
	err := terraformCommand(folder)
	if err != nil {
		return fmt.Errorf("please consider to install terraform. Terraform format fails: %w", err)
	}

	return nil
}

func GoFormat(filename string) error {
	inputFileContent, err := osReadFile(filename)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}

	formattedContent, err := formatSource(inputFileContent)
	if err != nil {
		return fmt.Errorf("error formatting source code: %w", err)
	}

	err = osWriteFile(filename, formattedContent, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}

	return nil
}
