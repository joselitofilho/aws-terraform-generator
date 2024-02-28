# Templates

## Variables

The following variables can be used within the templates:

### API Gateway

| Name           | Description                                                 |
| :------------- | :---------------------------------------------------------- |
| APIDomain      | The domain associated with an API.                          |
| StackName      | The name of the stack associated with the API.              |

Default templates:

```
📦 apigateway
 ┣ 📂 tmpls
 ┗ ┗ 📜 apig.tf.tmpl
```
- [📜 apig.tf.tmpl](./internal/generators/apigateway/tmpls/apig.tf.tmpl)

### API Gateway Lambda

| Name               | Description                                             |
| :----------------- | :------------------------------------------------------ |
| StackName          | The name of the stack associated with the Lambda.       |
| Name               | The name of the Lambda function.                        |
| AsModule           | If true, the Lambda will be created as module, otherwise as resource. |
| Source             | The source of the Lambda function module.               |
| RoleName           | The role name of the Lambda execution role.             |
| Runtime            | Identifier of the Lambda runtime.                       |
| Description        | Description of the Lambda function.                     |
| Envars             | Environment variables associated with the Lambda.       |
| Verb               | HTTP verb associated with the Lambda (if applicable).   |
| Path               | Path associated with the Lambda (if applicable).        |
| Files              | Map containing files related to the Lambda. The key is the name of the file. |
| ┗ Imports          | A list of imports required for each file.               |
| ┗ Tmpl             | The template content of each file.                      |

Default temaplates:

```
📦 apigateway
 ┣ 📂 tmpls
 ┃ ┣ 📜 lambda.go.tmpl
 ┃ ┣ 📜 lambda.tf.tmpl
 ┗ ┗ 📜 main.go.tmpl
 ```
- [📜 lambda.go.tmpl](./internal/generators/apigateway/tmpls/lambda.go.tmpl)
- [📜 lambda.tf.tmpl](./internal/generators/apigateway/tmpls/lambda.tf.tmpl)
- [📜 main.go.tmpl](./internal/generators/apigateway/tmpls/main.go.tmpl)

### Kinesis

| Name            | Description                                                |
| :-------------- | :--------------------------------------------------------- |
| Name            | The name of the SQS queue.                                 |
| RetentionPeriod | The duration for which records are retained.               |
| KMSEncription   | Indicates whether server-side encryption is enabled using AWS Key Management Service (KMS). |
| KMSKeyID        | The ID of the AWS Key Management Service (KMS) key used for encryption, if enabled. |

Default temaplates:

```
📦 kinesis
 ┣ 📂 tmpls
 ┗ ┗ 📜 kinesis.tf.tmpl
```
- [📜 kinesis.tf.tmpl](./internal/generators/kinesis/tmpls/kinesis.tf.tmpl)

### Lambda

| Name                | Description                                            |
| :------------------ | :----------------------------------------------------- |
| Name                | The name of the Lambda.                                |
| AsModule            | If true, the Lambda will be created as module, otherwise as resource. |
| Source              | The source of the Lambda.                              |
| RoleName            | The role name of the Lambda execution role.            |
| Runtime             | Identifier of the Lambda runtime.                      |
| Description         | Description of the Lambda.                             |
| Envars              | Environment variables associated with the Lambda.      |
| SQSTriggers         | List of SQS triggers associated with the Lambda.       |
| ┗ SourceARN         | The Amazon Resource Name (ARN) of the SQS queue.       |
| Crons               | List of cron jobs associated with the Lambda.          |
| ┗ ScheduleExpression | The cron expression defining the schedule.            |
| ┗ IsEnabled         | Indicates whether the cron job is enabled.             |
| Files               | Map containing files related to the Lambda. The key is the name of the file. |
| ┗ Imports           | A list of imports required for each file.              |
| ┗ Tmpl              | The template content of each file.                     |

Default temaplates:

```
📦 lambda
 ┣ 📂 tmpls
 ┃ ┣ 📜 lambda.go.tmpl
 ┃ ┣ 📜 lambda.tf.tmpl
 ┗ ┗ 📜 main.go.tmpl
```
- [📜 lambda.go.tmpl](./internal/generators/lambda/tmpls/lambda.go.tmpl)
- [📜 lambda.tf.tmpl](./internal/generators/lambda/tmpls/lambda.tf.tmpl)
- [📜 main.go.tmpl](./internal/generators/lambda/tmpls/main.go.tmpl)

### S3 Buckets

| Name           | Description                                                 |
| :------------- | :---------------------------------------------------------- |
| Name           | The name of the S3 bucket.                                  |
| ExpirationDays | The number of days after which objects will expire.         |

Default temaplates:

```
📦 s3
 ┣ 📂 tmpls
 ┗ ┗ 📜 s3.tf.tmpl
```
- [📜 s3.tf.tmpl](./internal/generators/s3/tmpls/s3.tf.tmpl)

### SNS

| Name           | Description                                                 |
| :------------- | :---------------------------------------------------------- |
| Name           | The name of the SNS topic.                                  |
| BucketName     | The name of the S3 bucket for S3 notifications.             |
| Lambdas        | List of Lambda functions subscribed to the SNS topic.       |
| SQSs           | List of SQS queues subscribed to the SNS topic.             |

The `Lambdas` and `SQSs` are both of the `SNSResource` type, representing data associated with resources subscribed to an SNS topic.

| Name           | Description                                                 |
| :------------- | :---------------------------------------------------------- |
| Name           | The name of the subscribed resource.                        |
| Events         | The events for which notifications are triggered.           |
| FilterPrefix   | Prefix-based filtering for messages.                        |
| FilterSuffix   | Suffix-based filtering for messages.                        |

Default temaplates:

```
📦 sns
 ┣ 📂 tmpls
 ┗ ┗ 📜 sns.tf.tmpl
```
- [📜 sns.tf.tmpl](./internal/generators/sns/tmpls/sns.tf.tmpl)

### SQS

| Name            | Description                                                |
| :-------------- | :--------------------------------------------------------- |
| Name            | The name of the SQS queue.                                 |
| MaxReceiveCount | The maximum number of times a message can be received (int32). |

Default temaplates:

```
📦 sqs
 ┣ 📂 tmpls
 ┗ ┗ 📜 sqs.tf.tmpl
```
- [📜 sqs.tf.tmpl](./internal/generators/sqs/tmpls/sqs.tf.tmpl)

### Structure

| Name           | Description                                                 |
| :------------- | :---------------------------------------------------------- |
| StackName      | The name of the stack associated with the project structure. |

## Custom Functions

The following custom functions are available:

| Name           | Description                                                 |
| :------------- | :---------------------------------------------------------- |
| getFileByName  | Retrieves a file from a map of files by its name.           |
| getFileImports | Retrieves the imports of a file by its name.                |
| ToCamel        | Converts a string to CamelCase format.                      |
| ToKebab        | Converts a string to kebab-case format.                     |
| ToLower        | Converts a string to lowercase.                             |
| ToPascal       | Converts a string to PascalCase format.                     |
| ToSpace        | Converts a string to kebab-case and replaces hyphens with spaces. |
| ToSnake        | Converts a string to snake_case format.                     |
| ToUpper        | Converts a string to uppercase.                             |
