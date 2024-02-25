package utils

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
)

func TerraformFormat(folder string) error {
	cmd := exec.Command("terraform", "fmt", folder)

	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("please consider to install terraform. Terraform format fails: %w", err)
	}

	return nil
}

func GoFormat(fileName string) error {
	inputFileContent, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}

	formattedContent, err := format.Source(inputFileContent)
	if err != nil {
		return fmt.Errorf("error formatting source code: %w", err)
	}

	err = os.WriteFile(fileName, formattedContent, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}

	return nil
}
