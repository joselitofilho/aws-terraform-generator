package kinesis

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joselitofilho/aws-terraform-generator/internal/fmtcolor"
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
		return fmt.Errorf("%w: %w", generatorserrs.ErrYAMLParse, err)
	}

	modPath := path.Join(k.output, "mod")
	_ = os.MkdirAll(modPath, os.ModePerm)

	result := make([]string, 0, len(yamlConfig.Kinesis))

	templates := utils.MergeStringMap(defaultTfTemplateFiles,
		generators.CreateTemplatesMap(yamlConfig.OverrideDefaultTemplates.Kinesis))

	tg := generators.NewGenerator()

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

			generators.MustGenerateFiles(tg, nil, filesConf, data, modPath)

			fmtcolor.White.Printf("Kinesis '%s' has been generated successfully\n", conf.Name)

			continue
		}

		output, err := tg.Build(data, "kinesis-tf-template", templates[filenameKinesisTf])
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		result = append(result, output)
	}

	if len(result) > 0 {
		outputFile := path.Join(modPath, filenameKinesisTf)

		generators.MustGenerateFile(tg, nil, filenameKinesisTf, strings.Join(result, "\n"), outputFile, Data{})

		fmtcolor.White.Println("Kinesis has been generated successfully")
	}

	return nil
}
