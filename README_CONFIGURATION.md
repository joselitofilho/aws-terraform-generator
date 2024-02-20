# Configuration

## Structure

The configuration is organized into the following sections:

- **Stacks**: Configuration for different stacks.
- **Default Templates**: Default Terraform templates for creating stacks.
- **Lambdas**: Configuration for lambda functions.
- **API Gateways**: Configuration for API Gateways.
- **SQS**: Configuration for SQS queues.
- **Buckets**: Configuration for S3 buckets.
- **RESTful APIs**: Configuration for RESTful APIs.

## Stacks

Each stack configuration includes folders for different environments (`dev`, `uat`, `prd`, etc.), default templates, and specific configurations for lambdas, API gateways, SQS, buckets, and RESTful APIs.

### Folder Structure

Each environment folder contains the following Terraform files:

- `main.tf`: Main Terraform configuration.
- `terragrunt.hcl`: Terragrunt configuration.
- `vars.tf`: Variable definitions.

### Default Templates

Default Terraform templates are provided for creating stacks. These templates include backend configuration, provider configuration, module instantiation, and variable definitions.

### Lambdas

Lambda configurations include lambda function names, descriptions, environment variables, SQS triggers, cron schedules, and code configurations.

### API Gateways

API Gateway configurations include stack names, API domain names, lambda associations, and code configurations.

### SQS

SQS configurations include queue names and maximum receive counts.

### Buckets

S3 bucket configurations include bucket names, object keys, and source paths.

### RESTful APIs

RESTful API configurations include API names.

### Full example

[template.yaml](template.yaml)

```yaml
# Structure for stacks with multiple environments
structure:
  # Stacks section
  stacks:
    # Stack for mystack service
    - stack_name: mystack
      # Folders for different environments
      folders:
        # Development environment
        - name: dev
          files:
            # Terraform configuration files for dev environment
            - name: main.tf
            - name: terragrunt.hcl
            - name: vars.tf
        # User Acceptance Testing environment
        - name: uat
          files:
            # Terraform configuration files for uat environment
            - name: main.tf
            - name: terragrunt.hcl
            - name: vars.tf
        # Production environment
        - name: prd
          files:
            # Terraform configuration files for prd environment
            - name: main.tf
            - name: terragrunt.hcl
            - name: vars.tf
        # Common module
        - name: mod
          files:
            # Terraform configuration files for module
            - name: main.tf
              # Template for generating stack_name based on environment
              tmpl: |-
                locals {
                  stack_name = "{{$.StackName}}-${var.environment}"
                }
            - name: vars.tf
        # Lambda functions
        - name: lambda

  # Default templates section
  default_templates:
    # Default Terraform template files
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

# Lambdas configuration section
lambdas:
  # Configuration for the coexampleceiver lambda function
  - source: git@github.com:username/terraform-aws-lambda?ref=reference
    name: exampleReceiver
    description: "Trigger on schedule and initiate the execution of example receiver"
    envars:
      # Environment variables
      - MYAPI_API_BASE_URL: var.myapi_api_base_url
        MYAPI_USER: var.myapi_user
        MYAPI_PASSWORD_SECRET: aws_secretsmanager_secret.myapi_password.name
        DOCDB_HOST: var.docdb_host
        DOCDB_USER: var.docdb_user
        DOCDB_PASSWORD_SECRET: var.docdb_password_secret
        SQS_QUEUE_URL: aws_sqs_queue.target_sqs.id
    # SQS triggers for lambda function
    sqs-triggers:
      - source_arn: aws_sqs_queue.source_sqs.arn
    # Cron schedule for lambda function
    crons:
      - schedule_expression: cron(0 1 * * ? *)
        is_enabled: var.trigger_enabled
    # Code configuration for lambda function
    code:
      - key: lambda
        imports: # Optional. We can specify what imports we want to add in the output GoLang file.
          - "github.com/mylogging/logging"

# API Gateways configuration section
apigateways:
  # Stack configuration for the mystack API Gateway
  - stack_name: example
    api_domain: example-api.example-${var.environment}.com
    apig: true
    # Lambdas associated with the mystack API Gateway
    lambdas:
      - 
        name: exampleAPIReceiver
        description: Trigger the example API receiver via API Gateway
        verb: POST
        path: /v1/mystacks
        # Code configuration for the lambda associated with the API Gateway
        code:
          - key: lambda
            imports: # Optional. We have the option to specify the imports to include in the resulting GoLang file.
              - "github.com/logging"
          - key: main
            tmpl:
              |- # Optional. We can specify the template for the output GoLang file.
              package main

              import(
                  "github.com/aws/aws-lambda-go/lambda"
              )

              func main() {
                  // TODO
                  lambda.Start({{$.Name}}Lambda.run)
              }

# SQS configuration section
sqs:
  # Configuration for the target SQS queue
  - name: target
    max_receive_count: 15
  # Configuration for the source SQS queue
  - name: source
    max_receive_count: 10

# Buckets configuration section
buckets:
  # Configuration for the my-bucket S3 bucket
  - name: my-bucket
    key: "new_object_key"
    source: "path/to/file"

# Restfulapis configuration section
restfulapis:
  # Configuration for the Oracle RESTful API
  - name: MyAPI
```
