# IAM Role para GitHub Actions
resource "aws_iam_role" "github_actions" {
  name = "github-actions-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRoleWithWebIdentity"
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Condition = {
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:${var.github_repository}:*"
          }
        }
      }
    ]
  })

  depends_on = [aws_iam_openid_connect_provider.github]
}

resource "aws_iam_role_policy" "github_actions" {
  name = "github-actions-policy"
  role = aws_iam_role.github_actions.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "apigateway:*",
          "cognito-idp:*",
          "dynamodb:*",
          "iam:*",
          "lambda:*",
          "ssm:*",
          "s3:*",
          "logs:*",
          "cloudformation:*",
          "ec2:*",
          "sts:*"
        ]
        Resource = "*"
      },
      {
        Effect = "Allow"
        Action = [
          "s3:ListBucket",
          "dynamodb:DescribeTable"
        ]
        Resource = [
          "arn:aws:s3:::morango-pay-terraform-state",
          "arn:aws:dynamodb:sa-east-1:*:table/terraform-lock"
        ]
      }
    ]
  })

  depends_on = [aws_iam_role.github_actions]
}

resource "aws_iam_openid_connect_provider" "github" {
  url = "https://token.actions.githubusercontent.com"

  client_id_list = ["sts.amazonaws.com"]

  thumbprint_list = [
    "6938fd4d98bab03faadb97b34396831e3780aea1",
    "1c58a3a8518e8759bf075b76b750d4f2df264fcd"
  ]
}

data "aws_caller_identity" "current" {}

variable "github_repository" {
  description = "GitHub repository (format: owner/repo)"
  type        = string
  default     = "LucasLCabral/moranGo-pay"
}

output "github_actions_role_arn" {
  description = "ARN da role para GitHub Actions"
  value       = aws_iam_role.github_actions.arn
}