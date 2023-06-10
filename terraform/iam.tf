resource "aws_iam_user" "snapshot_s3_user" {
  name = "snapshot_user"

  tags = {
    Project = var.project_name
  }
}

resource "aws_iam_access_key" "snapshot_user_access_key" {
  user = aws_iam_user.snapshot_s3_user.name
}

output "snapshot_user_access_key_secret" {
  value = aws_iam_access_key.snapshot_user_access_key

  sensitive = true
}

resource "aws_iam_policy" "snapshots_bucket_read_write_policy" {
  name   = "eurostat_snapshots_s3_bucket_read_write_policy"
  policy = data.aws_iam_policy_document.allow_read_and_write_to_bucket.json
}

data "aws_iam_policy_document" "allow_read_and_write_to_bucket" {
  version = "2012-10-17"
  statement {
    effect = "Allow"
    actions = [
      "s3:PutObject",
      "s3:GetObject",
      "s3:DeleteObject"
    ]
    resources = [
      "arn:aws:s3:::eurostat-weekly-deaths-snapshots/*",
    ]
  }

  statement {
    effect = "Allow"
    actions = [
      "s3:ListBucket"
    ]
    resources = [
      "arn:aws:s3:::eurostat-weekly-deaths-snapshots"
    ]
  }
}

resource "aws_iam_user_policy_attachment" "s3_access_for_snapshot_user" {
  user       = aws_iam_user.snapshot_s3_user.name
  policy_arn = aws_iam_policy.snapshots_bucket_read_write_policy.arn
}
