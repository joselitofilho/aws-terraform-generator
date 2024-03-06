module "location_big_query_emitter_lambda" {
  source = "git@github.com:username/terraform-aws-lambda?ref=reference"

  stack_name                               = local.stack_name
  lambda_function_description              = "locationBigQueryEmitter lambda"
  lambda_function_throttles_alarm_disabled = true
  lambda_function_name                     = "locationBigQueryEmitter"
  lambda_function_kms_key_arn              = var.lambda_function_kms_key_arn
  lambda_function_sns_topic_monitoring_arn = var.alerting_sns_topic_arn
  lambda_function_source_base_path         = var.lambda_function_source_base_path
  lambda_function_existing_execute_role    = "arn:aws:iam::${var.account_id}:role/execute_lambda"
  lambda_function_vpc_config               = var.lambda_function_vpc_config

  lambda_function_env_vars = {
    TRACE          = "1"
    TRACE_ENTITIES = "Y"
    TIME_LOCATION  = "UTC"

  }

  client      = var.client
  environment = var.environment
  region      = var.region
  account_id  = var.account_id
}

resource "aws_lambda_permission" "location_big_query_emitter_allow_kinesis" {
  statement_id  = "AllowExecutionFromKinesis"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.location_big_query_emitter_lambda.function_name
  principal     = "kinesis.amazonaws.com"
}

resource "aws_lambda_event_source_mapping" "location_big_query_emitter_kinesis_mapping" {
  event_source_arn  = aws_kinesis_stream.processed_location_events_kinesis.arn
  function_name     = aws_lambda_function.location_big_query_emitter_lambda.function_name
  batch_size        = 10
  starting_position = "TRIM_HORIZON"
}
