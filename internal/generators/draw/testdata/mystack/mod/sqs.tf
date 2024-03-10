// source SQS queue
resource "aws_sqs_queue" "source_sqs" {
  name                       = "${var.client}-${var.environment}-source"
  visibility_timeout_seconds = 720

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.source_sqs_dlq.arn
    maxReceiveCount     = 10
  })

  depends_on = [aws_sqs_queue.source_sqs_dlq]
}

// source DLQ queue
resource "aws_sqs_queue" "source_sqs_dlq" {
  name                       = "${var.client}-${var.environment}-source-dlq"
  visibility_timeout_seconds = 720
}

// target SQS queue
resource "aws_sqs_queue" "target_sqs" {
  name                       = "${var.client}-${var.environment}-target"
  visibility_timeout_seconds = 720

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.target_sqs_dlq.arn
    maxReceiveCount     = 10
  })

  depends_on = [aws_sqs_queue.target_sqs_dlq]
}

// target DLQ queue
resource "aws_sqs_queue" "target_sqs_dlq" {
  name                       = "${var.client}-${var.environment}-target-dlq"
  visibility_timeout_seconds = 720
}

// processor SQS queue
resource "aws_sqs_queue" "processor_sqs" {
  name                       = "${var.client}-${var.environment}-processor"
  visibility_timeout_seconds = 720

  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.processor_sqs_dlq.arn
    maxReceiveCount     = 10
  })

  depends_on = [aws_sqs_queue.processor_sqs_dlq]
}

// processor DLQ queue
resource "aws_sqs_queue" "processor_sqs_dlq" {
  name                       = "${var.client}-${var.environment}-processor-dlq"
  visibility_timeout_seconds = 720
}
