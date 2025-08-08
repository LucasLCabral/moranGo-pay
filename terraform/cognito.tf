variable "region" {
  type    = string
  default = "sa-east-1"
}


resource "aws_cognito_user_pool" "users" {
  name = "morangopay-userpool"

  # verifica e-mail automaticamente
  auto_verified_attributes = ["email"]
  username_attributes      = ["email"]

  password_policy {
    minimum_length    = 8
    require_lowercase = true
    require_numbers   = true
    require_symbols   = false
    require_uppercase = true
  }

  admin_create_user_config {
    allow_admin_create_user_only = false
  }
}

resource "aws_cognito_user_pool_client" "app_client" {
  name         = "morangopay-client"
  user_pool_id = aws_cognito_user_pool.users.id

  # não gera secret (útil para clients públicos / testes)
  generate_secret = false

  # flows necessários para login via USER_PASSWORD_AUTH e refresh token
  explicit_auth_flows = [
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_SRP_AUTH"
  ]

  # token validity (hour)
  access_token_validity  = 1
  id_token_validity      = 1
  refresh_token_validity = 30
}

output "cognito_user_pool_id" {
  value = aws_cognito_user_pool.users.id
}

output "cognito_app_client_id" {
  value = aws_cognito_user_pool_client.app_client.id
}