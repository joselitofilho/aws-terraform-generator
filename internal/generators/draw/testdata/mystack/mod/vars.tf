variable "client" {
  type = string
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "account_id" {
  type = string
}

variable "zone_id" {
  type = string
}

variable "alerting_sns_topic_arn" {
  type = string
}

variable "lambda_function_source_base_path" {
  type = string
}

variable "lambda_function_vpc_config" {
  type = map(list(string))
}

variable "lambda_function_kms_key_arn" {
  type = string
}