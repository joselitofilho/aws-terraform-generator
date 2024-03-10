resource "aws_s3_bucket" "storage_bucket" {
  bucket = "${var.client}-${var.environment}-storage"
}

resource "aws_s3_bucket_acl" "storage_acl" {
  bucket = aws_s3_bucket.storage_bucket.id
  acl    = "private"
}

resource "aws_s3_bucket_lifecycle_configuration" "storage_bucket_config" {
  bucket = aws_s3_bucket.storage_bucket.id

  rule {
    id = "expiration"

    expiration {
      days = 90
    }

    status = "Enabled"
  }
}
