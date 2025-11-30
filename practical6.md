# Practical 6 — Infrastructure as Code with Terraform and LocalStack

Overview
--------
This practical demonstrates how to write, test and secure cloud infrastructure using Infrastructure as Code (IaC). The exercise provisions AWS-compatible resources locally with LocalStack, deploys a Next.js static site to an S3 bucket using Terraform, and validates infrastructure security using Trivy. The goal is to build repeatable, testable and secure infrastructure workflows that can be validated locally before moving to a real cloud environment.

Learning outcomes
-----------------
- Model cloud infrastructure using Terraform and run it against LocalStack for local validation.
- Deploy a Next.js static site (statically exported) to an S3 website bucket via IaC.
- Scan Terraform code (and generated plan) with Trivy to find IaC misconfigurations and security issues.
- Understand secure defaults for S3 hosting (encryption, logging, least-privilege policies).
- Automate common lifecycle tasks with scripts to make development and testing reproducible.

Technologies used
-----------------
- Terraform — declarative IaC tool.
- LocalStack — local AWS emulator to run and test AWS services offline.
- AWS S3 — storage used for static website hosting (emulated by LocalStack).
- Next.js — framework used to generate a static website for deployment.
- Trivy — security scanner for IaC and container images.
- Docker & Docker Compose — to run LocalStack and other tooling.

Prerequisites
-------------
Install and verify these tools on your machine before following the guide:

- Docker & Docker Compose
- Terraform (1.0 or newer)
- `tflocal` (terraform wrapper for LocalStack) or configure `terraform` with LocalStack endpoints
- Node.js (>= 18) and npm/yarn
- AWS CLI and `awslocal` CLI (comes with LocalStack)
- Trivy (for scanning IaC)

Verify with:

```bash
docker --version
docker-compose --version
terraform --version
tflocal --version   # if using tflocal wrapper
node --version
npm --version
aws --version
awslocal --version
trivy --version
```

Repository layout (relevant)
----------------------------
Example structure expected for this practical (adapt to your repo):

- `nextjs-app/` — the Next.js project that produces a static `out/` folder after `npm run build` and `npm run export` (or `next export`).
- `terraform/` — Terraform configuration (providers, S3 buckets, IAM settings, outputs).
- `scripts/` — helper scripts: `setup.sh`, `status.sh`, `scan.sh`, `compare-security.sh`.
- `assets/` — (optional) screenshots and evidence files used in the report.

Quick start — manual steps
--------------------------
These commands reproduce the workflow used in the practical. They assume the repo root contains `nextjs-app`, `terraform` and `scripts` directories.

1. Start LocalStack (via Docker Compose):

```bash
./scripts/setup.sh
```

2. Build the Next.js static site:

```bash
cd nextjs-app
npm ci
npm run build
# if using next export or similar static export step:
npm run export
cd ..
```

3. Initialize and apply Terraform against LocalStack:

```bash
cd terraform
tflocal init
tflocal apply -auto-approve
cd ..
```

4. Upload built static site to the deployment bucket (using `awslocal`):

```bash
DEPLOY_BUCKET=$(cd terraform && terraform output -raw deployment_bucket_name)
awslocal s3 sync nextjs-app/out/ s3://$DEPLOY_BUCKET/ --delete
```

5. Inspect status and endpoints:

```bash
./scripts/status.sh
```

Terraform configuration highlights
---------------------------------
The Terraform code used in this practical focuses on secure defaults for S3 hosting and logging. Key design points:

- Provider configuration is adapted for LocalStack endpoints and uses dummy credentials for local emulation.
- Two buckets are created: a `deployment` bucket for website content and a `logs` bucket to receive server access logs.
- The `deployment` bucket is configured as a static website host with a public read policy that is narrowly scoped to objects used by the website.
- Server-side encryption is enabled for both buckets using AES256 (SSE-S3) to enforce data-at-rest protection even in the local emulator.
- Bucket policy attaches logging permissions and explicitly enables `s3:PutObject` on the logs bucket from the deployment bucket.

Example Terraform considerations (not full TF code):

- Use outputs for bucket names so CI and deployment scripts can read them:

```hcl
output "deployment_bucket_name" {
  value = aws_s3_bucket.deployment.id
}
```

- Provider block (LocalStack example):

```hcl
provider "aws" {
  region                      = "us-east-1"
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_requesting_account_id  = true
  endpoints {
    s3 = "http://localstack:4566"
  }
}
```

Security scanning with Trivy
----------------------------
Trivy is used to scan Terraform templates and detect misconfigurations and insecure defaults. Two configurations were compared in this practical:

1. **Secure configuration** — encryption enabled, restricted policies, logging enabled. Trivy reports no critical/high issues in this mode.
2. **Insecure configuration** — encryption disabled, overly permissive IAM or public write ACLs. Trivy highlights critical/high findings such as `S3 bucket without encryption` or `Public access to sensitive bucket`.

Commands used:

```bash
./scripts/scan.sh terraform       # scans secure config
./scripts/scan.sh insecure        # scans intentionally-insecure config
./scripts/compare-security.sh     # compares results side-by-side
```

What Trivy checks (examples):

- S3 buckets without server-side encryption
- Publicly readable or writable buckets
- Overly permissive IAM policies (e.g., `*` actions or `*` principals)
- Missing logging configurations
- Deprecated or insecure resource attributes

Evidence and screenshots
------------------------
Include the following evidence files in `assets/` (or link them in your report):

![alt text](<Practicals/practical_06/assets/Screenshot from 2025-11-30 12-51-50.png>)

![alt text](<Practicals/practical_06/assets/Screenshot from 2025-11-30 12-52-09.png>)

![alt text](<Practicals/practical_06/assets/Screenshot from 2025-11-30 12-52-20.png>)

![alt text](<Practicals/practical_06/assets/Screenshot from 2025-11-30 12-52-42.png>)
-----------------------
The practical demonstrates how small differences in Terraform configuration create large security differences. For example:

- Enabling `server_side_encryption_configuration` on S3 prevents cleartext storage of sensitive objects.
- Enabling `logging` to a separate encrypted bucket provides an audit trail for requests.
- Tightening bucket policy to only allow `s3:GetObject` for `arn:aws:iam::...:role/website` or the AWS principal that serves content reduces exposure compared to an open `Principal = *` policy.

Sample remediation checklist (when Trivy reports issues):

1. Enable SSE (server-side encryption) on S3 buckets.
2. Turn on S3 server access logging and send logs to an encrypted logs bucket.
3. Remove `public` write or put policies; only allow `GetObject` for public hosting where necessary.
4. Replace wildcard IAM principals/actions with least-privilege statements and resource ARNs.
5. Add explicit `prevent_destroy` lifecycle rules for production-critical buckets to avoid accidental deletion.

Reflection — why scan IaC?
-------------------------
Scanning IaC is essential because insecure infrastructure definitions are a primary cause of cloud incidents. IaC lets you introduce misconfigurations at scale quickly; scanning finds those mistakes early, reducing blast radius and remediation cost.

How LocalStack helps the workflow
---------------------------------
- LocalStack provides a fast feedback loop for infrastructure code without incurring cloud costs or waiting for remote deployments.
- It enables testing of Terraform plans, provider behavior and post-deploy tasks (like `awslocal s3 sync`) in CI pipelines or local dev environments.
- Note: LocalStack is an emulator — for production readiness always run final validation against the real cloud provider.

Key takeaways
-------------
- IaC makes infrastructure reproducible and version controlled; combine it with security scans to make deployments safe.
- LocalStack is excellent for local testing and CI-level validation, but it is not a perfect substitute for real cloud testing.
- Trivy (and similar scanners) catch misconfigurations such as missing encryption, open access policies, or insecure resource attributes.
- Automating the pipeline (scripts and CI) reduces human error and improves developer productivity.

Next steps and recommended improvements
--------------------------------------
- Integrate Trivy and Terraform plan scanning into a CI workflow (GitHub Actions) to automatically block merges when critical issues are found.
- Add a `terraform fmt` and `terraform validate` check in CI to enforce style and syntactic correctness.
- Use Terragrunt or module patterns for larger infra to reduce duplication and centralize secure defaults.
- Add automated tests (localstack + test harness) that validate the website is reachable after the `awslocal s3 sync` step.
- Protect critical buckets with Terraform `lifecycle` rules and use remote state with encryption in real deployments.

Appendix — helpful scripts and snippets
--------------------------------------
Example `scripts/setup.sh` (simplified):

```bash
#!/usr/bin/env bash
set -euo pipefail

# Start LocalStack (assumes docker-compose.yml present)
docker-compose up -d
echo "Waiting for LocalStack to become healthy..."
sleep 10
awslocal s3 mb s3://local-logs || true
```

Example `scripts/scan.sh` (simplified):

```bash
#!/usr/bin/env bash
target=${1:-terraform}
echo "Scanning $target with Trivy..."
trivy config --format table --severity CRITICAL,HIGH $target
```

Example `scripts/compare-security.sh` (simplified):

```bash
#!/usr/bin/env bash
./scripts/scan.sh terraform > secure.txt
./scripts/scan.sh insecure > insecure.txt || true
echo "--- Secure vs Insecure ---"
diff -u secure.txt insecure.txt || true
```

Acknowledgements & References
----------------------------
- LocalStack: https://localstack.cloud/
- Terraform: https://www.terraform.io/
- Trivy: https://aquasecurity.github.io/trivy/
- Next.js: https://nextjs.org/




