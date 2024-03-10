resource "aws_lambda_function" "loction_event_processor_lambda" {
  filename      = "./artifacts/loction_event_processor.zip"
  function_name = "loction_event_processor"
  description   = "loctionEventProcessor lambda"
  role          = aws_iam_role.execute_lambda.arn
  handler       = "loction_event_processor"

  source_code_hash = filebase64sha256("./artifacts/loction_event_processor.zip")

  runtime = "go1.x"

  environment {
    variables = {
      TRACE                                        = "1"
      TRACE_ENTITIES                               = "Y"
      TIME_LOCATION                                = "UTC"
      DOCDB_HOST                                   = var.doc_db_host
      DOCDB_PASSWORD_SECRET                        = var.doc_db_password_secret
      DOCDB_USER                                   = var.doc_db_user
      PROCESSED_LOCATION_EVENTS_KINESIS_STREAM_URL = aws_kinesis_stream.processed_location_events_kinesis.name
    }
  }
}

// loctionEventProcessor SQS trigger rule for lambda
resource "aws_lambda_event_source_mapping" "loction_event_processor_lambda_sqs_trigger" {
  event_source_arn = aws_sqs_queue.location_update_sqs.arn
  function_name    = aws_lambda_function.loction_event_processor_lambda.arn
  batch_size       = 1
  enabled          = true
}
