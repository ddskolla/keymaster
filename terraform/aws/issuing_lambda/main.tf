/**
 * Usage:
 *
 * ```hcl
 * module "issuing_lambda" {
 *   source = "github.com/bsycorp/keymaster/terraform/aws/issuing_lambda"
 *
 *   # The environment label will be added to all named resources
 *   env_label   = "myproject-npe"
 *
 *   # Keymaster configuration file
 *   configuration = {
 *      CONFIG: "s3://km-tools-bls-01/km.yaml"
 *   }
 *
 *   # List of target roles that the lambda may issue creds for
 *   target_role_arns = [
 *    "arn:aws:iam::218296299700:role/test_env_admin"
 *   ]
 *
 *   # List of client accounts that may invoke issuing lambda
 *   client_account_arns = [
 *    "arn:aws:iam::062921715666:root",   # myproj-dev-01
 *   ]
 *
 *   # Enable auto-creation of the configuration bucket
 *   config_bucket_enable = true
 *   config_file_upload_enable = true
 *   config_file_name = "${path.module}/test_api_config.yaml"
 *
 *   resource_tags = {
 *     Name         = "baz"
 *     Created-By   = "you@your.com"
 *   }
 * }
 * ```
  */