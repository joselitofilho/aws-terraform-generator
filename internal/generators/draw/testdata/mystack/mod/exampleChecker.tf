module "example_checker_lambda" {
  source = "git@github.com:username/terraform-aws-lambda?ref=reference"

  stack_name                               = local.stack_name
  lambda_function_description              = "exampleChecker lambda"
  lambda_function_throttles_alarm_disabled = true
  lambda_function_name                     = "exampleChecker"
  lambda_function_kms_key_arn              = var.lambda_function_kms_key_arn
  lambda_function_sns_topic_monitoring_arn = var.alerting_sns_topic_arn
  lambda_function_source_base_path         = var.lambda_function_source_base_path
  lambda_function_existing_execute_role    = "arn:aws:iam::${var.account_id}:role/execute_lambda"
  lambda_function_vpc_config               = var.lambda_function_vpc_config

  lambda_function_env_vars = {
    TRACE                   = "1"
    TRACE_ENTITIES          = "Y"
    TIME_LOCATION           = "UTC"
    PROCESSOR_SQS_QUEUE_URL = aws_sqs_queue.processor_sqs.name
    STORAGE_S3_BUCKET       = aws_s3_bucket.storage_bucket.bucket
    STORAGE_S3_DIRECTORY    = "example_checker_files"

  }

  client      = var.client
  environment = var.environment
  region      = var.region
  account_id  = var.account_id
}

// exampleChecker SQS trigger rule for lambda
resource "aws_lambda_event_source_mapping" "example_checker_lambda_sqs_trigger" {
  event_source_arn = aws_sqs_queue.target_sqs.arn
  function_name    = aws_lambda_function.example_checker_lambda.arn
  batch_size       = 1
  enabled          = true
}
