# JWT Secret para assinatura de tokens
resource "random_password" "jwt_secret" {
  length           = 32
  special          = true
  override_special = "!@#$%^&*()_+-=[]{}|;:,.<>?"
}

# SSM Parameter para armazenar o JWT Secret de forma segura
resource "aws_ssm_parameter" "jwt_secret" {
  name        = "/morango-pay/jwt-secret"
  description = "JWT Secret Key for Morango Pay API"
  type        = "SecureString"
  value       = random_password.jwt_secret.result
  
  tags = {
    Environment = "development"
    Project     = "morango-pay"
  }
}

# IAM Role para o Lambda
resource "aws_iam_role" "lambda_role" {
  name = "morango-pay-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# Política básica de execução do Lambda
resource "aws_iam_role_policy_attachment" "lambda_basic" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Política para acessar Cognito
resource "aws_iam_role_policy" "lambda_cognito" {
  name = "lambda-cognito-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "cognito-idp:InitiateAuth",
          "cognito-idp:SignUp",
          "cognito-idp:ConfirmSignUp",
          "cognito-idp:AdminConfirmSignUp",
          "cognito-idp:GetUser",
          "cognito-idp:AdminGetUser",
          "cognito-idp:ListUsers"
        ]
        Resource = aws_cognito_user_pool.users.arn
      }
    ]
  })
}

# Política para acessar SSM Parameter
resource "aws_iam_role_policy" "lambda_ssm" {
  name = "lambda-ssm-policy"
  role = aws_iam_role.lambda_role.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters"
        ]
        Resource = aws_ssm_parameter.jwt_secret.arn
      }
    ]
  })
}

# Lambda Function
resource "aws_lambda_function" "api" {
  filename         = "${path.module}/../app/lambda.zip"
  function_name    = "MorangoPay-API"
  role            = aws_iam_role.lambda_role.arn
  handler         = "bootstrap"
  runtime         = "provided.al2"
  architectures    = ["x86_64"]
  timeout         = 30

  source_code_hash = filebase64sha256("${path.module}/../app/lambda.zip")

  environment {
    variables = {
      COGNITO_USER_POOL_ID = aws_cognito_user_pool.users.id
      COGNITO_CLIENT_ID    = aws_cognito_user_pool_client.app_client.id
      JWT_SECRET_KEY       = aws_ssm_parameter.jwt_secret.name
    }
  }
}