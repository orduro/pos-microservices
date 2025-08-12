terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-southeast-2"
}

resource "aws_s3_bucket" "pos_microservices_bucket" {
  bucket = "pos-microservices-bucket-${random_id.bucket_suffix.hex}"

  tags = {
    Name        = "POS Microservices Storage"
    Environment = "production"
    Project     = "pos-microservices"
  }
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# configure bucket versioning
resource "aws_s3_bucket_versioning" "pos_bucket_versioning" {
  bucket = aws_s3_bucket.pos_microservices_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

# configure server-side encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "pos_bucket_encryption" {
  bucket = aws_s3_bucket.pos_microservices_bucket.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# allow public read access for images
resource "aws_s3_bucket_public_access_block" "pos_bucket_pab" {
  bucket = aws_s3_bucket.pos_microservices_bucket.id

  block_public_acls       = true
  block_public_policy     = false
  ignore_public_acls      = true
  restrict_public_buckets = false
}

# bucket policy for public read access to images
resource "aws_s3_bucket_policy" "pos_bucket_policy" {
  bucket = aws_s3_bucket.pos_microservices_bucket.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource  = "${aws_s3_bucket.pos_microservices_bucket.arn}/*"
      }
    ]
  })

  depends_on = [aws_s3_bucket_public_access_block.pos_bucket_pab]
}

# configure lifecycle management
# - delete old versions after 30 days
# - clean up incomplete uploads after 7 days
resource "aws_s3_bucket_lifecycle_configuration" "pos_bucket_lifecycle" {
  bucket = aws_s3_bucket.pos_microservices_bucket.id

  rule {
    id     = "cleanup_old_versions"
    status = "Enabled"

    filter {}

    noncurrent_version_expiration {
      noncurrent_days = 30
    }

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }
}
