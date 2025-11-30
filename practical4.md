# Practical 4 — Snyk SAST Integration with GitHub Actions

Table of Contents
- **Project Overview**: Short description of the demo app and goals.
- **What This Practical Covers**: Focus and learning outcomes.
- **Approach**: How I implemented the integration.
- **Implementation**: Step-by-step configuration and sample YAML snippets.
- **Evidence & Validation**: How I verified the setup.
- **Challenges**: Problems encountered and fixes.
- **Lessons Learned & Recommendations**: Practical takeaways and next steps.
- **Conclusion**: Summary of outcomes.

**Project Overview**
- **Application**: A Spring Boot REST API used to demonstrate CI/CD and security automation.
- **Key Endpoints**: `/` (health), `/version` (app version), `/nations` (sample data).
- **Dependencies**: Spring Boot Web, JavaFaker, Jackson Databind, JaCoCo for coverage.
- **Security Tools**: Snyk (SAST + dependency & container scanning), SonarCloud (quality), OWASP ZAP (DAST).

**What This Practical Covers**
- **Goal**: Automate static security scanning (SAST) using Snyk within GitHub Actions so vulnerabilities are found early and reported to the GitHub Security tab via SARIF.
- **Objectives**: Run Snyk on every push/PR, scan dependencies/code/container, upload SARIF results, and enable scheduled monitoring.

**Approach**
- **Incremental Integration**: Start small — add a basic Snyk step to the main Maven workflow, confirm it works, then add enhanced scenarios and scheduled runs.
- **Secrets & Access**: Store the Snyk token in `SNYK_TOKEN` GitHub Actions secret.
- **Fail-safes**: Use `continue-on-error` for security jobs and sensible severity thresholds to avoid blocking development while surfacing issues.

**Implementation — Practical Steps**

- **1. Get a Snyk API token**
	- Sign up at `https://snyk.io` and connect with GitHub.
	- From Account Settings → Auth Token, copy the token.

- **2. Add GitHub secret**
	- Repository Settings → Secrets and variables → Actions → `New repository secret`.
	- Name: `SNYK_TOKEN` (case-sensitive). Paste token.

- **3. Basic Snyk job inside `maven.yml`**
	- Run after tests. Build the project so Snyk can resolve dependencies.
	- Minimal example step:

```yaml
# snippet: security job in .github/workflows/maven.yml
security:
	name: Snyk Security Scan
	needs: test
	runs-on: ubuntu-latest
	steps:
		- uses: actions/checkout@v4
		- uses: actions/setup-java@v4
			with:
				java-version: '17'
				distribution: temurin
				cache: maven
		- name: Build for Snyk
			run: mvn clean compile -DskipTests
		- name: Run Snyk (dependency & code)
			uses: snyk/actions/maven@master
			continue-on-error: true
			env:
				SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
			with:
				args: --severity-threshold=high --sarif-file-output=snyk.sarif
		- name: Upload Snyk SARIF
			uses: github/codeql-action/upload-sarif@v2
			if: always()
			continue-on-error: true
			with:
				sarif_file: snyk.sarif
```

- **4. Enhanced workflow with matrix & change detection**
	- Create `.github/workflows/enhanced-security.yml` and use `dorny/paths-filter` to skip unnecessary scans.
	- Use a matrix like `scan-type: [dependencies, code]` to run different scans in parallel and use `if:` guards so scans only run when relevant files changed.

- **5. Container scanning**
	- Build the image in CI and run `snyk/actions/docker@master` to scan Docker images.
	- Example:

```yaml
- name: Run Snyk Container Scan
	uses: snyk/actions/docker@master
	continue-on-error: true
	env:
		SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
	with:
		image: cicd-demo:latest
		args: --severity-threshold=high --file=Dockerfile
```

- **6. Monitoring for production**
	- On pushes to `main`, call `snyk monitor` with the same project coordinates to enable continuous monitoring in the Snyk dashboard.

**Evidence & Validation**
- **What to check**: GitHub Actions run logs (security job), Snyk dashboard project scans, and GitHub repository Security → Code scanning alerts for SARIF findings.
- **Validation steps**:
	- Add an intentionally vulnerable dependency (e.g. older transitive library) and push; confirm Snyk reports it.
	- Confirm `snyk.sarif` is created and uploaded (check Code scanning alerts).
	- Run scheduled workflow or manually trigger the enhanced workflow to validate scheduled scans.

**Challenges Encountered**
- **Missing token**: Workflows failed with "SNYK_TOKEN is not set" — fixed by adding repository secret named exactly `SNYK_TOKEN`.
- **SARIF upload errors**: SARIF file not generated — fixed by adding `--sarif-file-output` and making upload step run with `if: always()`.
- **Build prerequisites**: Snyk dependency analysis needed a compiled project — added `mvn clean compile -DskipTests` beforehand.
- **Long runtime**: Multiple scans increased CI time — mitigated with change-detection filters and matrix parallelization.

**Lessons Learned & Recommendations**
- **Shift-left security**: Integrate Snyk early so issues are surfaced in PRs rather than production.
- **Balance feedback vs blocking**: Use `continue-on-error` and severity thresholds to provide actionable feedback without blocking developers.
- **Monitor & Triage**: Use Snyk monitor + GitHub Security tab to track vulnerabilities over time; assign owners for triage.
- **Automated fixes**: Enable Snyk's PRs for fix pull requests (where available) to streamline remediation.
- **Policy as code**: Consider adding automated rules for fail-on-critical or fail-on-unfixable depending on release policy.
- **Alerting**: Add notifications (Slack/email) for high-severity findings so they are not missed.
- **Cost & runtime trade-offs**: Schedule full scans weekly and run incremental scans on PRs to reduce runner costs.

# Screenshot Ouptut
![alt text](<Practicals/practical4/assets/Screenshot from 2025-11-30 12-28-22.png>)

![alt text](<Practicals/practical4/assets/Screenshot from 2025-11-30 12-28-41.png>)

![alt text](<Practicals/practical4/assets/Screenshot from 2025-11-30 12-29-01.png>)

![alt text](<Practicals/practical4/assets/Screenshot from 2025-11-30 12-29-17.png>)

**Notes**
- **SARIF quality**: When uploading SARIF, include `tool` metadata so GitHub matches findings to the right tool.
- **Pull request comments**: Use Snyk GitHub integration to comment on PRs with remediation information for failing dependencies.
- **Policy thresholds**: For stricter projects, use `--fail-on=upgradable` or `--fail-on=patchable` if the release process requires.
- **Automation pipeline**: Combine Snyk with `dependabot` for automatic dependency update PRs; use Snyk to validate those updates.
- **Security ownership**: Add a `SECURITY.md` and a triage playbook describing how to handle CVEs discovered by Snyk.

**Conclusion**
- **Outcome**: Snyk SAST and dependency scanning integrated into GitHub Actions, SARIF results uploaded to GitHub Security, and scheduled monitoring enabled.
- **Value**: Early detection of vulnerabilities, continuous monitoring, and improved developer visibility into security issues.
- **Next steps**: Add automated remediation PRs, tighten policy for production branches, and integrate alerting or ticket creation for high-severity findings.

---



