output "cognito_user_pool_id" {
  value = aws_cognito_user_pool.users.id
}

output "cognito_app_client_id" {
  value = aws_cognito_user_pool_client.app_client.id
}

output "api_gateway_url" {
  value = aws_apigatewayv2_stage.default.invoke_url
}

output "lambda_function_name" {
  value = aws_lambda_function.api.function_name
}

output "jwt_secret_parameter_name" {
  value = aws_ssm_parameter.jwt_secret.name
}