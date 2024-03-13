<div align="center">

# AWS Terraform Generator

[![GitHub tag](https://img.shields.io/github/release/joselitofilho/aws-terraform-generator?include_prereleases=&sort=semver&color=2ea44f&style=for-the-badge)](https://github.com/joselitofilho/aws-terraform-generator/releases/)
[![Go Report Card](https://goreportcard.com/badge/github.com/joselitofilho/aws-terraform-generator?style=for-the-badge)](https://goreportcard.com/report/github.com/joselitofilho/aws-terraform-generator)
[![Code coverage](https://img.shields.io/badge/Coverage-91.9%25-2ea44f?style=for-the-badge)](#)

[![Made with Golang](https://img.shields.io/badge/Golang-1.21.6-blue?logo=go&logoColor=white&style=for-the-badge)](https://go.dev "Go to Golang homepage")
[![Using Terraform](https://img.shields.io/badge/Terraform-3.76.1-blueviolet?logo=terraform&logoColor=white&style=for-the-badge)](https://registry.terraform.io/providers/hashicorp/aws/3.76.1/docs "Go to Terraform docs")
[![Using Diagrams](https://img.shields.io/badge/diagrams.net-orange?logo=&logoColor=white&style=for-the-badge)](https://app.diagrams.net/ "Go to Diagrams homepage")

[![BuyMeACoffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black)](https://www.buymeacoffee.com/joselitofilho)

</div>

# Overview

The AWS Terraform Generator is a powerful tool designed to simplify and streamline the process of creating Terraform configurations for AWS infrastructure. With this tool, you can quickly generate Terraform code to provision AWS resources such as EC2 instances, S3 buckets, RDS databases, Lambda functions, API Gateways, and much more.

[![Start Here](https://img.shields.io/badge/start%20here-blue?style=for-the-badge)](#recommended-step-by-step)

**Table of contents**

- [Install](#install)
- [Features](#features)
- [How it works](#how-it-works)
- [Recommended step by step](#recommended-step-by-step)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Install

 ```bash
 $ go install github.com/joselitofilho/aws-terraform-generator/cmd/aws-terraform-generator@latest
 ```

## Features:
- Generate initial stack infrastructure folders.
- Generate GoLang code and Terraform files.
- [Diagrams](https://app.diagrams.net/) integration: Generate everything based on the exported XML diagram.
- Customization Options: Tailor generated code to your specific requirements using customizable templates and configuration parameters.
- Best Practices: Adhere to AWS and Terraform best practices with automatically generated code that follows industry standards.
- [Supported resources](https://drive.google.com/file/d/1Lrh6SikW1bvGXrfJLRDFBB4BChQdAPqz/view?usp=sharing):
    - APIGateway
    - Cron
    - Database
    - Google BigQuery
    - Kinesis streams
    - Lambda
    - Restful API
    - SNS
    - SQS with DLQ
    - S3

## How it works

The code generator already comes with some pre-configured templates for generating Terraform and GoLang files. All generator 
configuration is based on YAML files, making it easy to customize the available resources and templates.

The first step is to write the configuration file specified [here](CONFIGURATION.md). You can also use this [example](example/) as a reference.

There you go! Now you can generate the structure of your project or the files based on the configured resources. You can execute the commands in any order.

If you're using [diagrams](https://app.diagrams.net/), you can also generate the initial configuration file based on the XML generated by the tool.

If you have any questions or suggestions, feel free to create an [issue](https://github.com/joselitofilho/aws-terraform-generator/issues). Your contribution is much appreciated.

<div style="text-align:center"><img src="assets/general-overview.svg" /></div>

## Recommended step by step

**Step 1**: Create a folder to organize the diagram and configuration files, ideally named after your stack.
```bash
$ mkdir mystack
```

**Step 2**: Create your diagram using [Diagrams](https://app.diagrams.net/). If you have already created one, proceed to the next step.

**Step 3**: Export and download your diagram as an XML file (file name suggestion: `diagram.xml`).
You can find instructions on how to do that at this link: https://www.drawio.com/doc/faq/export-to-xml.

Move the file to the folder created in the Step 1.

```bash
$ mv ~/Downloads/diagram.xml mystack/diagram.xml
```

**Step 4**: Create the diagram configuration file.

Suggestion [diagram.config.yaml](./example/diagram.config.yaml):
```bash
$ cp ./example/diagram.config.yaml mystack/diagram.config.yaml
```

Change the values according to your project.

**Step 5**: Create the structure configuration file.

Suggestion [structure.config.yaml](./example/structure.config.yaml):
```bash
$ cp ./example/structure.config.yaml mystack/structure.config.yaml
```

**Step 6**: Run the generator guide to assist you.

```bash
$ aws-terraform-generator --workdir mystack
```

Output:
```bash


                 ██████╗ ██████╗ ██████╗ ███████╗     ██████╗ ███████╗███╗   ██╗
                ██╔════╝██╔═══██╗██╔══██╗██╔════╝    ██╔════╝ ██╔════╝████╗  ██║
                ██║     ██║   ██║██║  ██║█████╗      ██║  ███╗█████╗  ██╔██╗ ██║
                ██║     ██║   ██║██║  ██║██╔══╝      ██║   ██║██╔══╝  ██║╚██╗██║
                ╚██████╗╚██████╔╝██████╔╝███████╗    ╚██████╔╝███████╗██║ ╚████║
                 ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝


? What would you like to do?  [Use arrows to move, type to filter]
> Generate a diagram config file
  Generate the initial structure
  Generate code
  Exit
```

## Usage

To use these configurations:

1. Navigate to the desired stack/environment folder.
2. Customize the Terraform files (`main.tf`, `vars.tf`, etc.) according to your requirements.
3. Run commands to manage the infrastructure.
```bash
$ aws-terraform-generator diagram -s mystack -c examples/diagram.config.yaml -d examples/diagram.drawio.xml -o examples/mystack.yaml
$ aws-terraform-generator structure -c examples/structure.yaml -o ./output
$ aws-terraform-generator apigateway -c examples/mystack.yaml -o ./output
$ aws-terraform-generator lambda -c examples/mystack.yaml -o ./output/mystack
$ aws-terraform-generator kinesis -c examples/mystack.yaml -o output/mystack
$ aws-terraform-generator s3 -c examples/mystack.yaml -o output/mystack
$ aws-terraform-generator sns -c examples/mystack.yaml -o output/mystack
$ aws-terraform-generator sqs -c examples/mystack.yaml -o output/mystack
$ aws-terraform-generator --workdir ./examples
```

## Configuration

All you need know regarding configuration you can find in the [configuration](CONFIGURATION.md) section.

[![open - Configuration](https://img.shields.io/badge/open-configuration-blue?style=for-the-badge)](CONFIGURATION.md "Go to configuration")

## Template

For code generation, we are using the standard Golang library [text/template](https://pkg.go.dev/text/template). Further details about the available variables and the definition of some added utility functions can be found in the [template](TEMPLATE.md) section.

[![open - Template](https://img.shields.io/badge/open-template-blue?style=for-the-badge)](TEMPLATE.md "Go to configuration")

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to create an [issue](https://github.com/joselitofilho/aws-terraform-generator/issues) or submit a pull request. Your contribution is much appreciated.

## License

This project is licensed under the [MIT License](LICENSE).
