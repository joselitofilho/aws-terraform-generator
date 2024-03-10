module "example_api_receiver_lambda" {
  source = "git@github.com:username/terraform-aws-lambda?ref=reference"

  stack_name                               = local.stack_name
  lambda_function_description              = "exampleAPIReceiver lambda"
  lambda_function_throttles_alarm_disabled = true
  lambda_function_name                     = "exampleAPIReceiver"
  lambda_function_name_prefix              = var.client
  lambda_function_vpc_config               = var.lambda_function_vpc_config
  lambda_function_kms_key_arn              = var.lambda_function_kms_key_arn
  lambda_function_sns_topic_monitoring_arn = var.alerting_sns_topic_arn
  lambda_function_source_base_path         = var.lambda_function_source_base_path
  lambda_function_existing_execute_role    = "arn:aws:iam::${var.account_id}:role/execute_lambda"

  lambda_function_env_vars = {
    REGION_AWS           = var.region
    TRACE_ENTITIES       = "Y"
    TRACE                = "1"
    SOURCE_SQS_QUEUE_URL = aws_sqs_queue.source_sqs.name

  }

  client      = var.client
  environment = var.environment
  region      = var.region
  account_id  = var.account_id
}

resource "aws_lambda_permission" "apigw_permission_example_api_receiver" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.example_api_receiver_lambda.arn
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.mystack_api.execution_arn}/*"
}

resource "aws_apigatewayv2_route" "apigw_route_example_api_receiver" {
  api_id    = aws_apigatewayv2_api.mystack_api.id
  route_key = "POST /v1/examples"
  target    = "integrations/${aws_apigatewayv2_integration.example_api_receiver.id}"
}

resource "aws_apigatewayv2_integration" "example_api_receiver" {
  api_id             = aws_apigatewayv2_api.mystack_api.id
  integration_type   = "AWS_PROXY"
  connection_type    = "INTERNET"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.example_api_receiver_lambda.invoke_arn
  lifecycle {
    ignore_changes = [
      passthrough_behavior
    ]
  }
}
