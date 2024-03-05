// location update trigger SQS queue
resource "aws_sqs_queue" "location_update_trigger_sqs" {
  name                       = "${var.client}-${var.environment}-location-update-trigger"
  visibility_timeout_seconds = 720

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.location_update_trigger_sqs_dlq.arn
    maxReceiveCount     = 10
  })

  depends_on = [aws_sqs_queue.location_update_trigger_sqs_dlq, aws_sqs_queue.location_retrieving_initiated_sqs_dlq]
}

// location update trigger DLQ queue
resource "aws_sqs_queue" "location_update_trigger_sqs_dlq" {
  name                       = "${var.client}-${var.environment}-location-update-trigger-dlq"
  visibility_timeout_seconds = 720
}

// location retrieving initiated SQS queue
resource "aws_sqs_queue" "location_retrieving_initiated_sqs" {
  name                       = "${var.client}-${var.environment}-location-retrieving-initiated"
  visibility_timeout_seconds = 720

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.location_retrieving_initiated_sqs_dlq.arn
    maxReceiveCount     = 10
  })

  depends_on = [aws_sqs_queue.location_retrieving_initiated_sqs_dlq]
}

// location retrieving initiated DLQ queue
resource "aws_sqs_queue" "location_retrieving_initiated_sqs_dlq" {
  name                       = "${var.client}-${var.environment}-location-retrieving-initiated-dlq"
  visibility_timeout_seconds = 720
}

// location update SQS queue
resource "aws_sqs_queue" "location_update_sqs" {
  name                       = "${var.client}-${var.environment}-location-update"
  visibility_timeout_seconds = 720

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.location_update_sqs_dlq.arn
    maxReceiveCount     = 10
  })

  depends_on = [aws_sqs_queue.location_update_sqs_dlq]
}

// location update DLQ queue
resource "aws_sqs_queue" "location_update_sqs_dlq" {
  name                       = "${var.client}-${var.environment}-location-update-dlq"
  visibility_timeout_seconds = 720
}

// location upload initiated SQS queue
resource "aws_sqs_queue" "location_upload_initiated_sqs" {
  name                       = "${var.client}-${var.environment}-location-upload-initiated"
  visibility_timeout_seconds = 720

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.location_upload_initiated_sqs_dlq.arn
    maxReceiveCount     = 10
  })

  depends_on = [aws_sqs_queue.location_upload_initiated_sqs_dlq]
}

// location upload initiated DLQ queue
resource "aws_sqs_queue" "location_upload_initiated_sqs_dlq" {
  name                       = "${var.client}-${var.environment}-location-upload-initiated-dlq"
  visibility_timeout_seconds = 720
}
