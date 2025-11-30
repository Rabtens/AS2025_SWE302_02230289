# Practical 4b — Dynamic Application Security Testing (DAST) with OWASP ZAP

Overview
--------
This practical demonstrates how to add Dynamic Application Security Testing (DAST) to a CI/CD workflow using OWASP ZAP (Zed Attack Proxy) and GitHub Actions. DAST inspects a running application from the outside — simulating attacker behavior — to uncover runtime problems such as missing security headers, unsafe server configurations, weak session handling, and exploitable inputs. This work complements the static analysis done in Practical 4a and helps build a layered (defense-in-depth) security posture.

Contents
- Goals and scope
- High-level approach
- Implementation details (rules, workflows, scans)
- Evidence and reporting
- Challenges and mitigations
- Lessons learned and recommendations

Goals and scope
----------------
- Add automated DAST to the repo that runs in CI and on a schedule.
- Provide fast, non-blocking checks for PRs and deeper scheduled scans for comprehensive coverage.
- Produce readable artifacts (HTML/JSON/Markdown) for developers and machine-readable outputs for automation.
- Tune rules to reduce noise and fail the pipeline only on meaningful risks.

Why DAST?
----------
DAST inspects the application at runtime. It finds configuration mistakes and runtime vulnerabilities that static analysis cannot see: missing headers, weak TLS configuration, unsafe HTTP methods, session and cookie issues, and endpoints that can be exploited. OWASP ZAP is a popular, free tool that supports baseline (quick) scans, full active scans, and API-focused scans.

Approach
--------
I followed a staged, practical strategy:

1. Understand the application: identify endpoints (e.g., `/nations`, `/currencies`) and confirm the app runs locally on port 5000.
2. Define three scan profiles:
   - Baseline: quick passive checks for PRs (1–2 minutes)
   - Full: active exploitation tests scheduled weekly (15–30 minutes)
   - API: targeted tests for REST interfaces when OpenAPI specs are available
3. Use Docker containers to run both the app and ZAP in a repeatable environment.
4. Add a `.zap/rules.tsv` to control which findings block a build and which only warn.
5. Expose reports as CI artifacts and tune scan options to fit GitHub Actions limits.

Implementation details
----------------------
Files I added or adapted:

- `.zap/rules.tsv` — rule table that maps ZAP check IDs to thresholds and actions (FAIL/WARN/IGNORE).
- `.github/workflows/zap-scan.yml` — baseline PR workflow.
- `.github/workflows/zap-full-scan.yml` — scheduled full-scan workflow.
- small helper scripts to wait for app readiness and collect reports.

Rules and policy
----------------
I used a TSV file with three columns: `rule_id<TAB>threshold<TAB>action`. The idea is to fail builds only for high-confidence, high-risk findings and warn for medium/low items while suppressing noisy informational checks.

Example (excerpt of `.zap/rules.tsv`):

```
# rule_id	threshold	action
40018	HIGH	FAIL    # SQL injection
40012	HIGH	FAIL    # Reflected XSS
10020	MEDIUM	FAIL   # Missing X-Frame-Options
10038	MEDIUM	WARN   # Missing Content-Security-Policy
10096	LOW	IGNORE    # Timestamp disclosure (dev-only noise)
```

Workflow: baseline scan (PRs)
----------------------------
The baseline workflow runs on pull requests and provides quick feedback. Key steps:

1. Build the application and package it into a container image.
2. Launch the container and wait until the health endpoint responds (use `curl -f` with a timeout loop, not a blind sleep).
3. Run the ZAP baseline action (`zaproxy/action-baseline`) with `.zap/rules.tsv` and AJAX spidering enabled for dynamic pages.
4. Always upload HTML/JSON/Markdown reports as artifacts and remove containers in an `if: always()` cleanup step.

Important bits (health check example):

```bash
timeout 60 bash -c 'until curl -f http://localhost:5000/; do sleep 2; done'
```

Workflow: full scan (scheduled)
------------------------------
The full scan uses `zaproxy/action-full-scan` and performs active attacks. It runs on a schedule (e.g., weekly) or manually. Because it is resource-intensive, I limited crawl depth, excluded static assets, and set a maximum runtime to fit within CI limits.

Optimization flags used:

```
-a -j -m 15 -config spider.maxDepth=5 -config spider.excludeUrl=".*\.(jpg|png|css|js)$"
```

API scanning
-----------
If an OpenAPI/Swagger spec is present, ZAP’s API scan mode (`zap-api-scan.py`) targets API endpoints directly and validates authorization, input handling and typical API-specific issues.

Report handling
---------------
- HTML: human-friendly report for developers
- JSON: machine-readable output for automation and custom checks
- Markdown: attachable to PR comments or documentation

Artifacts are uploaded via `actions/upload-artifact` and retained per repository policy.

Decision logic for failing a workflow
------------------------------------
Rather than failing on every alert, the pipeline examines the ZAP results and enforces rules from `.zap/rules.tsv`. During rollout the policy is permissive (WARN-only) and gradually moves to FAIL for high-risk checks once the team has time to remediate.

Example enforcement snippet (pseudo):

```bash
HIGH_ALERTS=$(jq '.site[0].alerts[] | select(.riskcode=="3")' zap_report.json | jq -s 'length')
if [ "$HIGH_ALERTS" -gt 0 ]; then
  echo "High-risk alerts found: $HIGH_ALERTS"; exit 1
fi
```

Evidence and reporting
----------------------

![alt text](<Practicals/practical4b/assets/Screenshot from 2025-11-30 13-07-46.png>)

![alt text](<Practicals/practical4b/assets/Screenshot from 2025-11-30 13-12-34.png>)

![alt text](<Practicals/practical4b/assets/Screenshot from 2025-11-30 13-12-53.png>)

Suggested screenshots to include in the practical:

- ZAP dashboard summary (high/medium/low counts)
- GitHub Actions run showing baseline and full scans
- Examples of PR decoration or artifact downloads

Challenges and mitigations
--------------------------
1) Application readiness
   - Problem: Scans started before the app was fully ready, producing false failures.
   - Fix: Use an HTTP health check loop with `curl -f` and a timeout instead of static sleeps.

2) Docker networking
   - Problem: Containers cannot access `localhost` of the runner when using separate namespaces.
   - Fix: In GitHub Actions the runner can access mapped ports on `localhost`. For multi-container setups use a shared Docker network or host networking for simplicity during CI runs.

3) Scan runtime and resource use
   - Problem: Full active scans exceed CI time limits.
   - Fix: Restrict depth, exclude static files, set a max duration (`-m`), and schedule full scans at off-hours.

4) Noisy results and false positives
   - Problem: Default rules produce many low-value findings.
   - Fix: Tune `.zap/rules.tsv` to IGNORE benign checks and WARN on medium items until they are addressed.

Lessons learned and recommendations
----------------------------------
- Use multiple scan profiles: fast baseline for PRs, deeper scheduled scans for thorough testing, and API scans when necessary.
- Start with permissive enforcement and tighten gates incrementally to avoid blocking development.
- Keep a documented exceptions list for known false positives and revisit it regularly.
- Combine DAST with SAST and dependency scanning for comprehensive coverage.
- Make readiness checks robust — test the actual HTTP response code and content when possible.
- Limit scan scope in CI to avoid timeouts and cost overruns; run exhaustive scans on a schedule.

Conclusion
----------
Adding OWASP ZAP to GitHub Actions gives realistic, attacker-style validation that complements static analysis. The recommended pipeline uses quick baseline scans during the PR review process and scheduled full scans for deeper coverage. Tuning and incremental enforcement are key to maintaining development velocity while improving runtime security.

