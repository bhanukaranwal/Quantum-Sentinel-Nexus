# infra/aws/main.tf
#
# This Terraform configuration sets up the foundational AWS infrastructure
# for the Quantum Sentinel Nexus platform.

# --- Provider Configuration ---
# Specifies the AWS provider and the region to deploy resources in.
provider "aws" {
  region = "us-east-1"
}

# --- Terraform Backend ---
# Configures where Terraform will store its state file. Using an S3 backend
# is a best practice for team collaboration, as it provides a central,
# versioned, and secure location for the state.
terraform {
  backend "s3" {
    bucket         = "qsn-terraform-state-bucket-unique-name" # <-- IMPORTANT: Change to a unique S3 bucket name
    key            = "global/aws/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "qsn-terraform-locks" # For state locking
  }
}

# --- Data Sources ---
# Fetches information about the current AWS environment.
data "aws_availability_zones" "available" {}
data "aws_caller_identity" "current" {}

# --- Locals ---
# Defines local variables for reuse throughout the configuration.
locals {
  cluster_name = "qsn-production-cluster"
  vpc_name     = "qsn-vpc"
}

# --- Networking (VPC) ---
# Creates a dedicated Virtual Private Cloud (VPC) for the QSN platform.
# This provides network isolation and security.
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.8.1"

  name = local.vpc_name
  cidr = "10.0.0.0/16"

  azs             = data.aws_availability_zones.available.names
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway = true
  enable_dns_hostnames = true

  tags = {
    Terraform   = "true"
    Environment = "production"
  }
}

# --- Kubernetes Cluster (EKS) ---
# Creates a managed Kubernetes cluster using AWS EKS.
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.8.4"

  cluster_name    = local.cluster_name
  cluster_version = "1.29"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  # EKS Managed Node Group for general purpose workloads.
  eks_managed_node_groups = {
    general_purpose = {
      min_size     = 2
      max_size     = 5
      instance_types = ["t3.large"]
    }
  }

  # EKS Managed Node Group for GPU-accelerated workloads (e.g., scoring-svc).
    gpu_workloads = {
      min_size     = 1
      max_size     = 3
      instance_types = ["g4dn.xlarge"] # GPU-enabled instance type
      # Taints would be used to ensure only GPU workloads run on these nodes.
      # taints = [{
      #   key    = "nvidia.com/gpu"
      #   value  = "true"
      #   effect = "NO_SCHEDULE"
      # }]
    }

  tags = {
    Terraform   = "true"
    Environment = "production"
  }
}

# --- Outputs ---
# Exports important information about the created infrastructure.
output "eks_cluster_endpoint" {
  description = "The endpoint for the EKS cluster's Kubernetes API server."
  value       = module.eks.cluster_endpoint
}

output "eks_cluster_name" {
  description = "The name of the EKS cluster."
  value       = module.eks.cluster_name
}
