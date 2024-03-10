module "example_receiver_lambda" {
  source = "git@github.com:username/terraform-aws-lambda?ref=reference"

  stack_name                               = local.stack_name
  lambda_function_description              = "exampleReceiver lambda"
  lambda_function_throttles_alarm_disabled = true
  lambda_function_name                     = "exampleReceiver"
  lambda_function_kms_key_arn              = var.lambda_function_kms_key_arn
  lambda_function_sns_topic_monitoring_arn = var.alerting_sns_topic_arn
  lambda_function_source_base_path         = var.lambda_function_source_base_path
  lambda_function_existing_execute_role    = "arn:aws:iam::${var.account_id}:role/execute_lambda"
  lambda_function_vpc_config               = var.lambda_function_vpc_config

  lambda_function_env_vars = {
    TRACE                = "1"
    TRACE_ENTITIES       = "Y"
    TIME_LOCATION        = "UTC"
    MY_APIAPI_BASE_URL   = var.my_apiapi_base_url
    MY_APIHOST           = var.my_apihost
    MY_APIUSER           = var.my_apiuser
    TARGET_SQS_QUEUE_URL = aws_sqs_queue.target_sqs.name

  }

  client      = var.client
  environment = var.environment
  region      = var.region
  account_id  = var.account_id
}

// exampleReceiver SQS trigger rule for lambda
resource "aws_lambda_event_source_mapping" "example_receiver_lambda_sqs_trigger" {
  event_source_arn = aws_sqs_queue.source_sqs.arn
  function_name    = aws_lambda_function.example_receiver_lambda.arn
  batch_size       = 1
  enabled          = true
}

// Trigger alarm for starting the exampleReceiver lambda
resource "aws_cloudwatch_event_rule" "example_receiver_cron" {
  name                = "runExampleReceiver"
  description         = "Trigger alarm for starting the exampleReceiver lambda"
  schedule_expression = "cron(0 2 * * ? *)"
  is_enabled          = true
}

resource "aws_cloudwatch_event_target" "example_receiver_cron_target" {
  rule = aws_cloudwatch_event_rule.example_receiver_cron.name
  arn  = aws_lambda_function.example_receiver_lambda.arn
}

resource "aws_lambda_permission" "example_receiver_allow_cron" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.example_receiver_lambda.arn
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.example_receiver_cron.arn
}