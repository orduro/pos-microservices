output "s3_bucket_name" {
  description = "Name of the S3 bucket"
  value       = aws_s3_bucket.pos_microservices_bucket.id
}

output "s3_bucket_arn" {
  description = "ARN of the S3 bucket"
  value       = aws_s3_bucket.pos_microservices_bucket.arn
}

output "s3_bucket_region" {
  description = "Region of the S3 bucket"
  value       = aws_s3_bucket.pos_microservices_bucket.region
}

output "s3_bucket_domain_name" {
  description = "Domain name of the S3 bucket"
  value       = aws_s3_bucket.pos_microservices_bucket.bucket_domain_name
}

output "s3_public_url_format" {
  description = "Format for public image URLs (replace {key} with actual file path)"
  value       = "https://${aws_s3_bucket.pos_microservices_bucket.bucket_domain_name}/{key}"
}
