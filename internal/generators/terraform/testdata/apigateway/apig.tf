locals {
  api_domain     = "stack-api.domain-${var.environment}.com"
  gateway_format = "{\"requestId\":\"$context.requestId\", \"ip\":$context.identity.sourceIp\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\", \"routeKey\":\"$context.routeKey\", \"path\":\"$context.path\", \"status\":\"$context.status\", \"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\", \"ErrMessage\":\"$context.error.message\"}"
}

resource "aws_apigatewayv2_api" "location_api" {
  name          = local.api_domain
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_stage" "location_api" {
  api_id      = aws_apigatewayv2_api.location_api.id
  name        = "$default"
  auto_deploy = true
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.location_api_logs.arn
    format          = local.gateway_format
  }
  lifecycle {
    ignore_changes = [
      deployment_id
    ]
  }
}

resource "aws_cloudwatch_log_group" "location_api_logs" {
  name = local.api_domain
}

resource "aws_apigatewayv2_domain_name" "location_api" {
  domain_name = local.api_domain

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.location_api_validation.certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}
resource "aws_route53_record" "location_api" {
  name    = aws_apigatewayv2_domain_name.location_api.domain_name
  type    = "A"
  zone_id = var.zone_id
  alias {
    name                   = aws_apigatewayv2_domain_name.location_api.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.location_api.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_apigatewayv2_api_mapping" "location_api" {
  api_id      = aws_apigatewayv2_api.location_api.id
  domain_name = aws_apigatewayv2_domain_name.location_api.id
  stage       = aws_apigatewayv2_stage.location_api.id
}

// if adding multiple domains here (SANs), then aws_route53_record will have to be able to recognise the correct zoneID
resource "aws_acm_certificate" "location_api" {
  domain_name       = local.api_domain
  validation_method = "DNS"
}

resource "aws_route53_record" "location_api_validation" {
  name    = tolist(aws_acm_certificate.location_api.domain_validation_options)[0].resource_record_name
  type    = tolist(aws_acm_certificate.location_api.domain_validation_options)[0].resource_record_type
  zone_id = var.zone_id
  records = [tolist(aws_acm_certificate.location_api.domain_validation_options)[0].resource_record_value]
  ttl     = 60
}

resource "aws_acm_certificate_validation" "location_api_validation" {
  certificate_arn         = aws_acm_certificate.location_api.arn
  validation_record_fqdns = [aws_route53_record.location_api_validation.fqdn]
}

// 5XXError: alarm for failed api invocations alarm
resource "aws_cloudwatch_metric_alarm" "api_5XXError_alarm" {
  alarm_name        = "${local.api_domain}_5XXError_alarm"
  alarm_description = "API 5XXError Alarm: ${local.api_domain}"

  namespace           = "AWS/ApiGateway"
  metric_name         = "5xx"
  statistic           = "Sum"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  threshold           = 1
  evaluation_periods  = 1
  period              = var.api_http_error_alarm_period
  treat_missing_data  = "notBreaching"

  alarm_actions = [var.alerting_sns_topic_arn]
  ok_actions    = [var.alerting_sns_topic_arn]

  dimensions = {
    ApiId = aws_apigatewayv2_api.location_api.id
  }
}

// 5XXError: alarm for failed api invocations alarm
resource "aws_cloudwatch_metric_alarm" "api_latency_alarm" {
  alarm_name        = "${local.api_domain}_LatencyError_alarm"
  alarm_description = "API Latency Alarm: ${local.api_domain}"

  namespace           = "AWS/ApiGateway"
  metric_name         = "Latency"
  statistic           = "Average"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  threshold           = var.api_latency_threshold_millis
  evaluation_periods  = 1
  period              = var.api_http_error_alarm_period
  treat_missing_data  = "notBreaching"

  alarm_actions = [var.alerting_sns_topic_arn]
  ok_actions    = [var.alerting_sns_topic_arn]

  dimensions = {
    ApiId = aws_apigatewayv2_api.location_api.id
  }
}
