// processed location events Kinesis
resource "aws_kinesis_stream" "processed_location_events_kinesis" {
  name             = "ProcessedLocationEvents"
  shard_count      = 1
  retention_period = 24

}
