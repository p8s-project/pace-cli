# {{ .Team }} Account Local Variables
# This file contains local variables used across the configuration
# 
# CUSTOMIZATION INSTRUCTIONS:
# - This file is NOT managed by automation and can be safely customized
# - Add custom values in the "Custom Configuration" section below
# - All other .tf files are managed by automation - do not edit them directly

locals {
  # Basic account information
  account_id        = ""
  environment       = "sbx"
  domain_name       = "{{ .Team }}"
  short_domain_name = "{{ .Team }}"

  # IAM roles and permissions
  root_dns                     = "${local.environment}.${module.metadata.global_domain}"
  dns                          = "${local.short_domain_name}.${local.environment}.${module.metadata.global_domain}"
  fixtures_deployment_role_arn = "arn:aws:iam::${local.account_id}:role/fixtures-deployment-role"

  # Common tags applied to resources
  common_tags = {
    Environment = local.environment
    Domain      = local.domain_name
    ManagedBy   = "terraform"
    Repository  = "otp-fixtures"
  }

  # ===================================================================
  # Custom Configuration - Safe to modify
  # ===================================================================

  # Service Repositories Configuration
  # List of GitHub repositories for service infrastructure management
  # Used for both OIDC role assumptions and webhook creation
  # Format: "organization/repository"
  service_repositories = []

  # EKS Configuration
  # EKS cluster version (change as needed)
  eks_cluster_version = "1.32"

  # Custom EKS node pools (leave empty for default node group)
  eks_custom_node_pools = {}

  # EKS cluster logging (leave empty for no additional logging)
  eks_cluster_log_types = []

  # EKS addons configuration (leave empty for defaults)
  eks_addons = {}

  # VPC Configuration  
  # Availability Zones (change as needed)
  vpc_availability_zones = ["us-west-2a", "us-west-2b"]

  # VPC Flow Logs (change as needed)
  vpc_enable_flow_log = true

  # Transit Gateway attachment (change as needed)
  vpc_enable_tgw_attachment = true

  # Route53 hosted zone associations (change as needed)
  vpc_trueaccord_hosted_zones = []
}
