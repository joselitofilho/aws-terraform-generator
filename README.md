# aws-terraform-generator

The AWS Terraform Generator is a powerful tool designed to simplify and streamline the process of creating Terraform configurations for AWS infrastructure. With this tool, you can quickly generate Terraform code to provision AWS resources such as EC2 instances, S3 buckets, RDS databases, Lambda functions, API Gateways, and much more.

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

## Configuration file

[Spec](README_CONFIGURATION.md)

## How it works

// TODO: 

## Usage

To use these configurations:

1. Navigate to the desired stack/environment folder.
2. Customize the Terraform files (`main.tf`, `vars.tf`, etc.) according to your requirements.
3. Run commands to manage the infrastructure.
```bash
$ aws-terraform-generator diagram -s mystack -i examples/diagram.drawio.xml -o examples/mystack.yaml
$ aws-terraform-generator structure -i examples/structure.yaml -o ./output
$ aws-terraform-generator apigateway -i examples/mystack.yaml -o ./output
$ aws-terraform-generator lambda -i examples/mystack.yaml -o ./output/mystack
$ aws-terraform-generator sqs -i examples/mystack.yaml -o output/mystack/mod/sqs.tf
$ aws-terraform-generator s3 -i examples/mystack.yaml -o output/mystack/mod/s3.tf
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
