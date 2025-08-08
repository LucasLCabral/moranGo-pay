resource "aws_apigatewayv2_api" "http_api" {
  name          = "dev-morango-api"
  protocol_type = "HTTP"
}

resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowAPIGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.hello.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.http_api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id                 = aws_apigatewayv2_api.http_api.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.hello.invoke_arn
  integration_method     = "POST"
  payload_format_version = "2.0"
}

# Hello Route
resource "aws_apigatewayv2_route" "hello_morango_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /hello"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"

  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito_jwt.id
}

# Wallet Route
resource "aws_apigatewayv2_route" "wallet_route" {
  api_id    = aws_apigatewayv2_api.http_api.id
  route_key = "GET /wallet"

  target = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"

  authorization_type = "JWT"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito_jwt.id
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.http_api.id
  name        = "$default"
  auto_deploy = true
}

# cognito 

resource "aws_apigatewayv2_authorizer" "cognito_jwt" {
  api_id = aws_apigatewayv2_api.http_api.id
  name   = "cognito-jwt-authorizer"

  authorizer_type = "JWT"

  identity_sources = ["$request.header.Authorization"]

  jwt_configuration {
    issuer   = "https://cognito-idp.${var.region}.amazonaws.com/${aws_cognito_user_pool.users.id}"
    audience = [aws_cognito_user_pool_client.app_client.id]
  }
}
