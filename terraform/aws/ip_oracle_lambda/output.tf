output "ip_oracle_sign_invoke_arn" {
  value = aws_lambda_function.ip_oracle_sign.invoke_arn
}

output "ip_oracle_verify_invoke_arn" {
  value = aws_lambda_function.ip_oracle_verify.invoke_arn
}