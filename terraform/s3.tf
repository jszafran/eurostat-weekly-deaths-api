resource "aws_s3_bucket" "weekly_deaths_snapshots_bucket" {
  bucket = "eurostat-weekly-deaths-snapshots"

  tags = {
    Project = var.project_name
  }

}
