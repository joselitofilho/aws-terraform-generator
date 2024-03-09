# Configuration

The configuration is organized into the following sections:

- [**Override default templates**](#override_default_templates): Configuration for overriding default templates.
- [**Diagram**](#diagram): Configuration for diagram.
- [**Structure**](#structure):
  - **Stacks**: Configuration for different stacks.
  - **Default Templates**: Default Terraform templates for creating stacks.
- [**API Gateways**](#apigateways): Configuration for API Gateways.
- [**Lambdas**](#lambdas): Configuration for lambda functions.
- [**Kinesis**](#kinesis): Configuration for Kinesis streams.
- [**SNS**](#sns): Configuration for SNS.
- [**SQS**](#sqs): Configuration for SQS.
- [**Buckets**](#buckets): Configuration for S3 buckets.
- [**RESTful APIs**](#restfulapis): Configuration for RESTful APIs.
- [**Draw**](#draw): Draw configurations.

### override_default_templates

Configuration for overriding default templates.

```yaml
override_default_templates:
  # Templates for API Gateway
  apigateway:
    # Terraform configuration for API Gateway
    - apig.tf: |-
        resource "aws_apigatewayv2_api" "{{$.StackName}}_api" {}
    # Lambda function code
    - lambda.go: |-
        type {{$.Name}}Lambda struct {}
    # Terraform configuration for Lambda function
    - lambda.tf: |-
        resource "aws_lambda_function" "{{ToSnake $.Name}}_lambda" {}
    # Main function code
    - main.go: |-
        func main() {}
  # Templates for Kinesis stream
  kinesis:
    # Terraform configuration for Kinesis stream
    - kinesis.tf: |-
        resource "aws_kinesis_stream" "{{ToSnake $.Name}}_kinesis" {}
  # Templates for Lambda function
  lambda:
    # Lambda function code
    - lambda.go: |-
        type {{$.Name}}Lambda struct {}
    # Terraform configuration for Lambda function
    - lambda.tf: |-
        resource "aws_lambda_function" "{{ToSnake $.Name}}_lambda" {}
    # Main function code
    - main.go: |-
        func main() {}
  # Templates for S3 bucket
  bucket:
    # Terraform configuration for S3 bucket
    - s3.tf: |-
        resource "aws_s3_bucket" "{{ToSnake $.Name}}_bucket" {}
  # Templates for SNS
  sns:
    # Terraform configuration for SNS topic
    - sns.tf: |-
        resource "aws_s3_bucket_notification" "s3_bucket_notification_{{ToSnake $.Name}}" {}
  # Templates for SQS
  sqs:
    # Terraform configuration for SQS queue
    - sqs.tf: |-
        resource "aws_sqs_queue" "{{ToSnake $.Name}}_sqs" {}
```

### diagram

Diagram configurations include modules to specify the URL pointing to the GitHub
repository for the resources module.

```yaml
diagram:
  # To specify the stack name for the diagram
  stack_name: mystack
  lambda:
    # URL pointing to the GitHub repository for the Lambda module
    # Replace "username" with the actual GitHub username
    # Replace "terraform-aws-lambda" with the actual repository name
    # Replace "reference" with the actual reference (branch, tag, or commit)
    source: git@github.com:username/terraform-aws-lambda?ref=reference
    # The name of the IAM role that will be assumed by the Lambda function
    role_name: execute_lambda
    # The runtime environment for the Lambda function (e.g., Python, Node.js, Go)
    runtime: go1.x
```

### structure

Structure for stacks with multiple environments.

```yaml
structure:
  # Stacks section. Each stack configuration includes folders for different environments (`dev`, `uat`, `prd`, etc.),
  # default templates, and specific configurations for lambdas, API gateways, SQS, and so on.
  stacks:
    - name: mystack
      # Folders for different environments Each environment folder contains the following Terraform files:
      #  - `main.tf`: Main Terraform configuration.
      #  - `terragrunt.hcl`: Terragrunt configuration.
      #  - `vars.tf`: Variable definitions.
      folders:
        # Development environment
        - name: dev
          # Terraform configuration files for dev environment
          files:
            - name: main.tf
            - name: terragrunt.hcl
            - name: vars.tf
        # User Acceptance Testing environment
        - name: uat
          # Terraform configuration files for uat environment
          files:
            - name: main.tf
            - name: terragrunt.hcl
            - name: vars.tf
        # Production environment
        - name: prd
          # Terraform configuration files for prd environment
          files:
            - name: main.tf
            - name: terragrunt.hcl
            - name: vars.tf
        # Common module
        - name: mod
          # Terraform configuration files for module
          files:
            - name: main.tf
              # Template for generating stack_name based on environment
              tmpl: |-
                locals {
                  stack_name = "{{$.StackName}}-${var.environment}"
                }
            - name: vars.tf
        # Lambda functions
        - name: lambda
  # Default templates are provided for creating stacks. These templates include backend configuration, provider
  # configuration, module instantiation, and variable definitions.
  default_templates:
    - main.tf: |-
        # Terraform backend and required providers configuration
        terraform {
          backend "s3" {
          }
          required_providers {
            aws = {
              source  = "hashicorp/aws"
              version = "~> 3.71"
            }
          }
        }

        # AWS provider configuration
        provider "aws" {
          region  = var.region
          profile = "${var.client}-sdv-${var.environment}"

          allowed_account_ids = [var.account_id]
        }

        # Module instantiation
        module "{{$.StackName}}" {
          source = "../mod"

          client      = var.client
          environment = var.environment
          region      = var.region
          account_id  = var.account_id

          // Variables from global

          dns_zone_id                      = var.zone_id
          alerting_sns_topic_arn           = var.alerting_sns_topic_arn
          lambda_function_source_base_path = var.lambda_function_source_base_path
          lambda_function_vpc_config       = var.lambda_function_vpc_config
          lambda_function_kms_key_arn      = var.lambda_function_kms_key_arn
        }

      terragrunt.hcl: |-
        # Terragrunt configuration
        include {
          path = find_in_parent_folders()
        }

      vars.tf: |-
        # Variables definition
        variable "client" {
          type = string
        }

        variable "environment" {
          type = string
        }

        variable "region" {
          type = string
        }

        variable "account_id" {
          type = string
        }

        variable "zone_id" {
          type = string
        }

        variable "alerting_sns_topic_arn" {
          type = string
        }

        variable "lambda_function_source_base_path" {
          type = string
        }

        variable "lambda_function_vpc_config" {
          type = map(list(string))
        }

        variable "lambda_function_kms_key_arn" {
          type = string
        }
```

### apigateways

API Gateway configurations include stack names, API domain names, lambda 
associations, and code configurations.

```yaml
apigateways:
  # To specify the stack name for the API Gateway
  - stack_name: mystack
    # The domain for the API Gateway
    api_domain: mystack-api.domain-${var.environment}.com
    # Indicates whether an API Gateway should be provisioned or not
    apig: true
    # Lambdas associated with the mystack API Gateway
    lambdas:
      - name: exampleAPIReceiver
        # URL pointing to the GitHub repository for the Lambda module
        # Replace "username" with the actual GitHub username
        # Replace "terraform-aws-lambda" with the actual repository name
        # Replace "reference" with the actual reference (branch, tag, or commit)
        source: git@github.com:username/terraform-aws-lambda?ref=reference
        # The name of the IAM role that will be assumed by the Lambda function
        role_name: execute_lambda
        # The runtime environment for the Lambda function (e.g., Python, Node.js, Go)
        runtime: go1.x
        # Description of the Lambda function
        description: Trigger the example API receiver via API Gateway
        # HTTP verb for the API Gateway endpoint
        verb: POST
        # The path for the API Gateway endpoint
        path: /v1/examples
        # Environment variables for the Lambda function
        envars:
          MYVAR: MYVAR_VALUE
        # File configuration for the lambda associated with the API Gateway
        files:
          - name: lambda.go
            # Optional. We have the option to specify the imports to include in the resulting GoLang file.
            imports:
              - context
              - github.com/logging
            # Optional. We can specify the template for the output GoLang file.
            tmpl: |-
              package main

              import (
                {{ range getFileImports $.Files "lambda.go" }}"{{ . }}"
                {{end}}
              )
          - name: main.go
            # Optional. We can specify the template for the output GoLang file.
            tmpl: |-
              package main

              import(
                  "github.com/aws/aws-lambda-go/lambda"
              )

              func main() {
                  // TODO
                  lambda.Start({{$.Name}}Lambda.run)
              }
```

### lambdas

Lambda configurations include lambda function names, descriptions, environment 
variables, SQS triggers, cron schedules, and code configurations.

```yaml
lambdas:
  # Name of the Lambda function
  - name: exampleReceiver
    # URL pointing to the GitHub repository for the Lambda module
    # Replace "username" with the actual GitHub username
    # Replace "terraform-aws-lambda" with the actual repository name
    # Replace "reference" with the actual reference (branch, tag, or commit)
    source: git@github.com:username/terraform-aws-lambda?ref=reference
    # The name of the IAM role that will be assumed by the Lambda function
    role_name: execute_lambda
    # The runtime environment for the Lambda function (e.g., Python, Node.js, Go)
    runtime: go1.x
    # Description of the Lambda function
    description: "Trigger on schedule and initiate the execution of example receiver"
    # Environment variables for the Lambda function
    envars:
      MYAPI_API_BASE_URL: var.myapi_api_base_url
      MYAPI_USER: var.myapi_user
      MYAPI_PASSWORD_SECRET: aws_secretsmanager_secret.myapi_password.name
      DOCDB_HOST: var.docdb_host
      DOCDB_USER: var.docdb_user
      DOCDB_PASSWORD_SECRET: var.docdb_password_secret
      SQS_QUEUE_URL: aws_sqs_queue.target_sqs.name
    # Kinesis triggers for the Lambda function
    kinesis-triggers:
      - source_arn: aws_kinesis_stream.mykinesis_kinesis.arn
    # SQS triggers for the Lambda function
    sqs-triggers:
      - source_arn: aws_sqs_queue.source_sqs.arn
    # Cron schedule for the Lambda function
    crons:
      - schedule_expression: cron(0 1 * * ? *)
        # Whether the trigger is enabled or not
        is_enabled: var.trigger_enabled
    # Optional. List of files that we can customize
    files:
      - name: lambda.go
        # Optional. We can specify what imports we want to add in the output GoLang file.
        imports:
          - "github.com/mylogging/logging"
        # Optional. We can specify the template for the output GoLang file.
        tmpl: |-
          package main
```

### kinesis

Kinesis configurations include stream names, retention period and KMS.

```yaml
kinesis:
  # Name of the Kinesis stream
  - name: myKinesis
    # Retention period for the Kinesis stream in hours
    retention_period: 24
    # KMS key ID for encryption
    kms_key_id: var.lambda_function_kms_key_arn
    # Custom Terraform file for defining the Kinesis stream resource
    files:
      - name: "custom.tf"
        # Template for the custom Terraform file
        tmpl: |-
          resource "aws_kinesis_stream" "{{ToSnake $.Name}}_kinesis" {
            # Add your custom configuration for the Kinesis stream here
          }
```

### sqs

SQS configurations include queue names and maximum receive counts.

```yaml
sqs:
  # Name of the SQS queue
  - name: target
    # Maximum number of times a message can be received from the queue before it's moved to the dead-letter queue
    max_receive_count: 15
    # Optional. List of files that we can customize
    files:
      - name: "target-sqs.tf"
        # Template for the Terraform file defining the target SQS queue resource
        tmpl: |-
          resource "aws_sqs_queue" "{{ToSnake $.Name}}_sqs" {}
  # Configuration for the source SQS queue
  - name: source
    # Maximum number of times a message can be received from the queue before it's moved to the dead-letter queue
    max_receive_count: 10
```

### sns

SNS configuration section.

```yaml
sns:
  # Name of the SNS notification
  - name: example
    # Name of the S3 bucket
    bucket_name: my-bucket
    # List of Lambda functions triggered by S3 events
    lambdas:
      - name: exampleReceiver
        # Events triggering Lambda
        events:
          - "s3:ObjectCreated:*" # Event indicating an object creation in S3
        # Optional. Prefix filter for S3 objects
        filter_prefix: "my_prefix"
        # Optional. Suffix filter for S3 objects
        filter_suffix: ".txt"
    # List of SQS to receive notification from an S3 bucket
    sqs:
      - name: target
        # SQS receiving notification from an S3 bucket
        events:
          - "s3:ObjectCreated:*" # Event indicating an object creation in S3
        # Optional. Prefix filter for S3 objects
        filter_prefix: "my_prefix"
        # Optional. Suffix filter for S3 objects
        filter_suffix: ".txt"
    # Optional. List of files that we can customize
    files:
      - name: "example-sns.tf"
        # Template for the Terraform file defining S3 bucket notification configuration
        tmpl: |-
          resource "aws_s3_bucket_notification" "s3_bucket_notification_{{ToSnake $.Name}}" {}
```

### buckets

S3 bucket configurations include bucket names, object keys, and source paths.

```yaml
buckets:
  # Name of the S3 bucket
  - name: my-bucket
    # Expiration period for objects in the bucket (in days)
    expiration-days: 90
    # Optional. List of files that we can customize
    files:
      - name: "my-bucket-s3.tf"
        # Template for the Terraform file defining the S3 bucket resource
        tmpl: |-
          resource "aws_s3_bucket" "{{ToSnake $.Name}}_bucket" {}
```

### restfulapis

RESTful API configurations include API names.

```yaml
restfulapis:
  # Name of the RESTful API
  - name: MyAPI
```

## Full example with comments

[fullexample.config.yaml](fullexample.config.yaml)

### draw

Draw configurations includes graph orientation, images and filters.

```yaml
draw:
  # Defines the direction of graph layout. See: https://graphviz.org/docs/attrs/rankdir/
  orientation: LR
  # Definitions of images for the available resources
  images:
    apigateway: "assets/diagram/api_gateway.svg"
    cron: "assets/diagram/cron.svg"
    database: "assets/diagram/database_dynamo_db.svg"
    endpoint: "assets/diagram/endpoint.svg"
    googlebq: "assets/diagram/google_bigquery.svg"
    kinesis: "assets/diagram/kinesis_data_stream.svg"
    lambda: "assets/diagram/lambda.svg"
    restfulapi: "assets/diagram/restful_api.svg"
    s3: "assets/diagram/s3_bucket.svg"
    sns: "assets/diagram/sns.svg"
    sqs: "assets/diagram/sqs.svg"
  # Filters for matching and excluding resources by name.
  filters:
    apigateway:
      match:
      not_match:
    cron:
      match:
      not_match:
    database:
      match:
      not_match:
    endpoint:
      match:
      not_match:
    googlebq:
      match:
      not_match:
    kinesis:
      # Patterns to match
      match:
        - "^ProcessedName" # Match regex pattern for ProcessedLocation
      # Patterns to exclude
      not_match:
        - "^ProcessedA" # Exclude regex pattern for ProcessedA
        - "^ProcessedB" # Exclude regex pattern for ProcessedB
    lambda:
      match:
      not_match:
    restfulapi:
      match:
      not_match:
    s3:
      match:
      not_match:
    sns:
      match:
      not_match:
    sqs:
      match:
      not_match:
```

Default images:

| Image                                       | Resource   | Path              |
| :-----------------------------------------: | :--------- | :---------------- |
| ![](assets/diagram/api_gateway.svg)         | apigateway | assets/diagram/api_gateway.svg |
| ![](assets/diagram/cron.svg)                | cron       | assets/diagram/cron.svg |
| ![](assets/diagram/database_dynamo_db.svg)  | database   | assets/diagram/database_dynamo_db.svg |
| ![](assets/diagram/endpoint.svg)            | endpoint   | assets/diagram/endpoint.svg |
| ![](assets/diagram/google_bigquery.svg)     | googlebq   | assets/diagram/google_bigquery.svg |
| ![](assets/diagram/kinesis_data_stream.svg) | kinesis    | assets/diagram/kinesis_data_stream.svg |
| ![](assets/diagram/lambda.svg)              | lambda     | assets/diagram/lambda.svg |
| ![](assets/diagram/restful_api.svg)         | restfulapi | assets/diagram/restful_api.svg |
| ![](assets/diagram/s3_bucket.svg)           | s3         | assets/diagram/s3_bucket.svg |
| ![](assets/diagram/sns.svg)                 | sns        | assets/diagram/sns.svg |
| ![](assets/diagram/sqs.svg)                 | sqs        | assets/diagram/sqs.svg |

- Available resources: [internal/resources/resource_type_enum.go](internal/resources/resource_type_enum.go)
- Recommend image size: 48px x 48px