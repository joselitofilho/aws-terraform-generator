module "example_processor_lambda" {
  source = "git@github.com:username/terraform-aws-lambda?ref=reference"

  stack_name                               = local.stack_name
  lambda_function_description              = "exampleProcessor lambda"
  lambda_function_throttles_alarm_disabled = true
  lambda_function_name                     = "exampleProcessor"
  lambda_function_kms_key_arn              = var.lambda_function_kms_key_arn
  lambda_function_sns_topic_monitoring_arn = var.alerting_sns_topic_arn
  lambda_function_source_base_path         = var.lambda_function_source_base_path
  lambda_function_existing_execute_role    = "arn:aws:iam::${var.account_id}:role/execute_lambda"
  lambda_function_vpc_config               = var.lambda_function_vpc_config

  lambda_function_env_vars = {
    TRACE                   = "1"
    TRACE_ENTITIES          = "Y"
    TIME_LOCATION           = "UTC"
    DOCDBDB_HOST            = var.docdbdb_host
    DOCDBDB_PASSWORD_SECRET = var.docdbdb_password_secret
    DOCDBDB_USER            = var.docdbdb_user

  }

  client      = var.client
  environment = var.environment
  region      = var.region
  account_id  = var.account_id
}

// exampleProcessor SQS trigger rule for lambda
resource "aws_lambda_event_source_mapping" "example_processor_lambda_sqs_trigger" {
  event_source_arn = aws_sqs_queue.processor_sqs.arn
  function_name    = aws_lambda_function.example_processor_lambda.arn
  batch_size       = 1
  enabled          = true
}
