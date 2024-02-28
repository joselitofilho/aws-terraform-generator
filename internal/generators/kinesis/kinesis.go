package kinesis

import (
	_ "embed"
	"fmt"
	"os"
	"path"

	"github.com/joselitofilho/aws-terraform-generator/internal/generators"
	"github.com/joselitofilho/aws-terraform-generator/internal/generators/config"
	generatorserrs "github.com/joselitofilho/aws-terraform-generator/internal/generators/errors"
	"github.com/joselitofilho/aws-terraform-generator/internal/utils"
)

type Data struct {
	Name            string
	KMSEncription   bool
	RetentionPeriod string
	KMSKeyID        string
}

type Kinesis struct {
	configFileName string
	output         string
}

func NewKinesis(configFileName, output string) *Kinesis {
	return &Kinesis{configFileName: configFileName, output: output}
}

func (k *Kinesis) Build() error {
	yamlParser := config.NewYAML(k.configFileName)

	yamlConfig, err := yamlParser.Parse()
	if err != nil {
		return fmt.Errorf("%w: %s", generatorserrs.ErrYAMLParse, err)
	}

	modPath := path.Join(k.output, "mod")
	_ = os.MkdirAll(modPath, os.ModePerm)

	tmplName := "kinesis-tf-template"
	result := ""

	for i := range yamlConfig.Kinesis {
		conf := yamlConfig.Kinesis[i]

		data := Data{
			Name:            conf.Name,
			KMSEncription:   conf.KMSKeyID != "",
			RetentionPeriod: conf.RetentionPeriod,
			KMSKeyID:        conf.KMSKeyID,
		}

		if len(conf.Files) > 0 {
			filesConf := generators.CreateFilesMap(conf.Files)

			err = generators.GenerateFiles(defaultTfTemplateFiles, filesConf, data, modPath)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			fmt.Printf("Kinesis '%s' has been generated successfully\n", conf.Name)

			continue
		}

		output, err := generators.Build(data, tmplName, string(kinesisTFTmpl))
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = fmt.Sprintf("%s\n%s", result, output)
	}

	if result != "" {
		outputFile := path.Join(modPath, "kinesis.tf")

		if err := generators.BuildFile(Data{}, tmplName, result, outputFile); err != nil {
			return fmt.Errorf("%w", err)
		}

		if err := utils.TerraformFormat(outputFile); err != nil {
			fmt.Println(err)
		}

		fmt.Println("Kinesis has been generated successfully")
	}

	return nil
}
