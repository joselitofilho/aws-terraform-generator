module "loction_event_processor_lambda" {
  source = "git@github.com:username/terraform-aws-lambda?ref=reference"

  stack_name                               = local.stack_name
  lambda_function_description              = "loctionEventProcessor lambda"
  lambda_function_throttles_alarm_disabled = true
  lambda_function_name                     = "loctionEventProcessor"
  lambda_function_kms_key_arn              = var.lambda_function_kms_key_arn
  lambda_function_sns_topic_monitoring_arn = var.alerting_sns_topic_arn
  lambda_function_source_base_path         = var.lambda_function_source_base_path
  lambda_function_existing_execute_role    = "arn:aws:iam::${var.account_id}:role/execute_lambda"
  lambda_function_vpc_config               = var.lambda_function_vpc_config

  lambda_function_env_vars = {
    TRACE                                        = "1"
    TRACE_ENTITIES                               = "Y"
    TIME_LOCATION                                = "UTC"
    DOCDB_HOST                                   = var.doc_db_host
    DOCDB_PASSWORD_SECRET                        = var.doc_db_password_secret
    DOCDB_USER                                   = var.doc_db_user
    PROCESSED_LOCATION_EVENTS_KINESIS_STREAM_URL = aws_kinesis_stream.processed_location_events_kinesis.name
  }

  client      = var.client
  environment = var.environment
  region      = var.region
  account_id  = var.account_id
}

// loctionEventProcessor SQS trigger rule for lambda
resource "aws_lambda_event_source_mapping" "loction_event_processor_lambda_sqs_trigger" {
  event_source_arn = aws_sqs_queue.location_update_sqs.arn
  function_name    = aws_lambda_function.loction_event_processor_lambda.arn
  batch_size       = 1
  enabled          = true
}
