module "oracle_location_data_receiver_lambda" {
  source = "git@github.com:username/terraform-aws-lambda?ref=reference"

  stack_name                               = local.stack_name
  lambda_function_description              = "oracleLocationDataReceiver lambda"
  lambda_function_throttles_alarm_disabled = true
  lambda_function_name                     = "oracleLocationDataReceiver"
  lambda_function_kms_key_arn              = var.lambda_function_kms_key_arn
  lambda_function_sns_topic_monitoring_arn = var.alerting_sns_topic_arn
  lambda_function_source_base_path         = var.lambda_function_source_base_path
  lambda_function_existing_execute_role    = "arn:aws:iam::${var.account_id}:role/execute_lambda"
  lambda_function_vpc_config               = var.lambda_function_vpc_config

  lambda_function_env_vars = {
    TRACE                                       = "1"
    TRACE_ENTITIES                              = "Y"
    TIME_LOCATION                               = "UTC"
    LOCATION_RETRIEVING_INITIATED_SQS_QUEUE_URL = aws_sqs_queue.location_retrieving_initiated_sqs.name
    ORACLE_API_BASE_URL                         = var.oracle_api_base_url
    ORACLE_HOST                                 = var.oracle_host
    ORACLE_USER                                 = var.oracle_user

  }

  client      = var.client
  environment = var.environment
  region      = var.region
  account_id  = var.account_id
}

// oracleLocationDataReceiver SQS trigger rule for lambda
resource "aws_lambda_event_source_mapping" "oracle_location_data_receiver_lambda_sqs_trigger" {
  event_source_arn = aws_sqs_queue.location_update_trigger_sqs.arn
  function_name    = aws_lambda_function.oracle_location_data_receiver_lambda.arn
  batch_size       = 1
  enabled          = true
}

// Trigger alarm for starting the oracleLocationDataReceiver lambda
resource "aws_cloudwatch_event_rule" "oracle_location_data_receiver_cron" {
  name                = "runOracleLocationDataReceiver"
  description         = "Trigger alarm for starting the oracleLocationDataReceiver lambda"
  schedule_expression = "cron(0 1 * * ? *)"
  is_enabled          = true
}

resource "aws_cloudwatch_event_target" "oracle_location_data_receiver_cron_target" {
  rule = aws_cloudwatch_event_rule.oracle_location_data_receiver_cron.name
  arn  = aws_lambda_function.oracle_location_data_receiver_lambda.arn
}

resource "aws_lambda_permission" "oracle_location_data_receiver_allow_cron" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.oracle_location_data_receiver_lambda.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.oracle_location_data_receiver_cron.arn
}