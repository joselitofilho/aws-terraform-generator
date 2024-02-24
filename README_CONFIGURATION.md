# Configuration

## Structure

The configuration is organized into the following sections:

- **Diagram**: Configuration for diagram.
- **Structure**:
  - **Stacks**: Configuration for different stacks.
  - **Default Templates**: Default Terraform templates for creating stacks.
- **Lambdas**: Configuration for lambda functions.
- **API Gateways**: Configuration for API Gateways.
- **SNS**: Configuration for SNS.
- **SQS**: Configuration for SQS.
- **Buckets**: Configuration for S3 buckets.
- **RESTful APIs**: Configuration for RESTful APIs.

## Full example with comments

[template.yaml](template.yaml)

# Templates

## Variables

The following variables can be used within the templates:

### API Gateway

| Name           | Description                                                 |
|----------------|-------------------------------------------------------------|
| APIDomain      | The domain associated with an API.                          |
| StackName      | The name of the stack associated with the API.              |

### API Gateway Lambda

| Name               | Description                                             |
|--------------------|---------------------------------------------------------|
| ModuleLambdaSource | The source of the Lambda function module.               |
| StackName          | The name of the stack associated with the Lambda.       |
| Name               | The name of the Lambda function.                        |
| Description        | Description of the Lambda function.                     |
| Envars             | Environment variables associated with the Lambda.       |
| Verb               | HTTP verb associated with the Lambda (if applicable).   |
| Path               | Path associated with the Lambda (if applicable).        |
| Files              | Map containing files related to the Lambda. The key is the name of the file. |
| ↳ Imports          | A list of imports required for each file.               |
| ↳ Tmpl             | The template content of each file.                      |

### Lambda

| Name                | Description                                            |
|---------------------|--------------------------------------------------------|
| ModuleLambdaSource  | The source of the Lambda function module.              |
| Name                | The name of the Lambda.                                |
| Description         | Description of the Lambda.                             |
| Envars              | Environment variables associated with the Lambda.      |
| SQSTriggers         | List of SQS triggers associated with the Lambda.       |
| ↳ SourceARN         | The Amazon Resource Name (ARN) of the SQS queue.       |
| Crons               | List of cron jobs associated with the Lambda.          |
| ↳ ScheduleExpression | The cron expression defining the schedule.            |
| ↳ IsEnabled         | Indicates whether the cron job is enabled.             |
| Files               | Map containing files related to the Lambda. The key is the name of the file. |
| ↳ Imports           | A list of imports required for each file.              |
| ↳ Tmpl              | The template content of each file.                     |

### S3 Buckets

| Name           | Description                                                 |
|----------------|-------------------------------------------------------------|
| Name           | The name of the S3 bucket.                                  |
| ExpirationDays | The number of days after which objects will expire.         |

### SNS

| Name           | Description                                                 |
|----------------|-------------------------------------------------------------|
| Name           | The name of the SNS topic.                                  |
| BucketName     | The name of the S3 bucket for S3 notifications.             |
| Lambdas        | List of Lambda functions subscribed to the SNS topic.       |
| SQSs           | List of SQS queues subscribed to the SNS topic.             |

The `Lambdas` and `SQSs` are both of the `SNSResource` type, representing data associated with resources subscribed to an SNS topic.

| Name           | Description                                                 |
|----------------|-------------------------------------------------------------|
| Name           | The name of the subscribed resource.                        |
| Events         | The events for which notifications are triggered.           |
| FilterPrefix   | Prefix-based filtering for messages.                        |
| FilterSuffix   | Suffix-based filtering for messages.                        |

### SQS

| Name            | Description                                                 |
|-----------------|-------------------------------------------------------------|
| Name            | The name of the SQS queue.                                  |
| MaxReceiveCount | The maximum number of times a message can be received (int32). |

### Structure

| Name           | Description                                                 |
|----------------|-------------------------------------------------------------|
| StackName      | The name of the stack associated with the project structure. |

## Custom Functions

The following custom functions are available:

| Name           | Description                                                 |
|----------------|-------------------------------------------------------------|
| getFileByName  | Retrieves a file from a map of files by its name.           |
| getFileImports | Retrieves the imports of a file by its name.                |
| ToCamel        | Converts a string to CamelCase format.                      |
| ToKebab        | Converts a string to kebab-case format.                     |
| ToLower        | Converts a string to lowercase.                             |
| ToPascal       | Converts a string to PascalCase format.                     |
| ToSpace        | Converts a string to kebab-case and replaces hyphens with spaces. |
| ToSnake        | Converts a string to snake_case format.                     |
| ToUpper        | Converts a string to uppercase.                             |