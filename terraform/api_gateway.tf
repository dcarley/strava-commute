resource "aws_api_gateway_rest_api" "api" {
  name = "${var.name}"
}

resource "aws_api_gateway_method" "method" {
  count         = "${length(var.http_methods)}"
  rest_api_id   = "${aws_api_gateway_rest_api.api.id}"
  resource_id   = "${aws_api_gateway_rest_api.api.root_resource_id}"
  http_method   = "${element(var.http_methods, count.index)}"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "integration" {
  count                   = "${length(var.http_methods)}"
  rest_api_id             = "${aws_api_gateway_rest_api.api.id}"
  resource_id             = "${aws_api_gateway_rest_api.api.root_resource_id}"
  http_method             = "${element(aws_api_gateway_method.method.*.http_method, count.index)}"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${data.aws_region.current.name}:lambda:path/2015-03-31/functions/${aws_lambda_function.lambda.arn}/invocations"
}

resource "aws_lambda_permission" "permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.lambda.arn}"
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.api.id}/*/*/"
}

resource "aws_api_gateway_deployment" "deployment" {
  stage_name  = "live"
  rest_api_id = "${aws_api_gateway_rest_api.api.id}"
  depends_on  = ["aws_api_gateway_integration.integration"]
}

output "url" {
  value = "${aws_api_gateway_deployment.deployment.invoke_url}"
}
