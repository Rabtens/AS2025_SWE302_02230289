# Snyk Quick Reference Guide

This quick reference provides essential commands and configurations for using Snyk with the cicd-demo project.

## Common Snyk CLI Commands

```bash
# Install Snyk CLI
npm install -g snyk

# Authenticate with Snyk
snyk auth

# Test for vulnerabilities
snyk test

# Test with specific severity threshold
snyk test --severity-threshold=high

# Test and output JSON format
snyk test --json

# Test and output SARIF format
snyk test --sarif

# Monitor project (send results to Snyk dashboard)
snyk monitor

# Test specific file
snyk test --file=pom.xml

# Test and exclude dev dependencies
snyk test --dev=false

# Ignore specific vulnerability
snyk ignore --id=SNYK-JAVA-ORGAPACHECOMMONS-1234567

# Test with policy file
snyk test --policy-path=.snyk
```

## GitHub Actions Snyk Integration

### Basic Configuration

```yaml
- name: Run Snyk to check for vulnerabilities
  uses: snyk/actions/maven@master
  env:
    SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
```

### Advanced Configuration

```yaml
- name: Run Snyk with options
  uses: snyk/actions/maven@master
  env:
    SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
  with:
    args: --severity-threshold=high --fail-on=upgradable
    command: test
```

### Multiple Scan Types

```yaml
# Dependency scanning
- uses: snyk/actions/maven@master
  with:
    command: test

# Code scanning (SAST)
- uses: snyk/actions/maven@master
  with:
    command: code test

# Container scanning
- uses: snyk/actions/docker@master
  with:
    image: your-image:tag

# Infrastructure as Code scanning
- uses: snyk/actions/iac@master
```

## Snyk CLI Options

| Option                 | Description                      | Example                     |
| ---------------------- | -------------------------------- | --------------------------- |
| `--severity-threshold` | Set minimum severity level       | `--severity-threshold=high` |
| `--fail-on`            | Fail conditions                  | `--fail-on=upgradable`      |
| `--dev`                | Include/exclude dev dependencies | `--dev=false`               |
| `--file`               | Specify manifest file            | `--file=pom.xml`            |
| `--json`               | JSON output format               | `--json`                    |
| `--sarif`              | SARIF output format              | `--sarif`                   |
| `--policy-path`        | Custom policy file path          | `--policy-path=.snyk`       |
| `--project-name`       | Custom project name              | `--project-name=my-app`     |
| `--target-reference`   | Git reference                    | `--target-reference=main`   |

## Severity Levels

| Level      | Description                       |
| ---------- | --------------------------------- |
| `low`      | Low severity vulnerabilities      |
| `medium`   | Medium severity vulnerabilities   |
| `high`     | High severity vulnerabilities     |
| `critical` | Critical severity vulnerabilities |

## .snyk File Structure

```yaml
version: v1.0.0

ignore:
  "SNYK-JAVA-ORGAPACHECOMMONS-1234567":
    - "*":
        reason: "Reason for ignoring"
        expires: "2024-12-31T23:59:59.999Z"

patches: {}

language-settings:
  java:
    # Java-specific settings

exclude:
  - "target/**"
  - "**/*.class"
```

## Common Vulnerability Actions

### Update Dependencies

```xml
<!-- Update to latest secure version -->
<dependency>
    <groupId>com.fasterxml.jackson.core</groupId>
    <artifactId>jackson-databind</artifactId>
    <version>2.15.2</version>
</dependency>
```

### Ignore Vulnerability (Temporary)

```yaml
# Add to .snyk file
ignore:
  "SNYK-JAVA-VULNERABILITY-ID":
    - "*":
        reason: "Explanation for ignoring"
        expires: "2024-12-31T23:59:59.999Z"
```

### Force Update Parent Dependency

```xml
<!-- Force specific version -->
<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-core</artifactId>
            <version>6.0.12</version>
        </dependency>
    </dependencies>
</dependencyManagement>
```

## Troubleshooting

### Common Issues

1. **Authentication Error**

   ```bash
   # Re-authenticate
   snyk auth
   ```

2. **No Manifest Found**

   ```bash
   # Ensure pom.xml exists and is valid
   mvn validate
   ```

3. **Rate Limiting**

   ```bash
   # Use organization token instead of personal
   snyk test --org=your-org-id
   ```

4. **Build Timeout**
   ```bash
   # Increase timeout in GitHub Actions
   timeout-minutes: 15
   ```

### Debug Commands

```bash
# Enable debug mode
export DEBUG=snyk*
snyk test

# Verbose output
snyk test -d

# Check Snyk version
snyk --version

# Test connectivity
snyk test --print-deps
```

## Security Best Practices

1. **Regular Scanning**: Run scans on every commit and PR
2. **Severity Thresholds**: Set appropriate thresholds for different environments
3. **Monitoring**: Use `snyk monitor` for continuous monitoring
4. **Policy Files**: Use `.snyk` files for consistent vulnerability management
5. **Secrets Management**: Store Snyk tokens securely in GitHub Secrets
6. **Documentation**: Document security decisions and ignore reasons

## Integration with Other Tools

### SonarQube Integration

```yaml
- name: SonarQube Scan
  uses: sonarqube-quality-gate-action@master

- name: Snyk Scan
  uses: snyk/actions/maven@master
```

### Slack Notifications

```yaml
- name: Notify Slack on Security Issues
  if: failure()
  uses: 8398a7/action-slack@v3
  with:
    status: failure
    text: "Security vulnerabilities detected"
  env:
    SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK }}
```

### JIRA Integration

```yaml
- name: Create JIRA Issue
  if: failure()
  uses: atlassian/gajira-create@v2
  with:
    project: SEC
    issuetype: Bug
    summary: "Security vulnerabilities in ${{ github.repository }}"
```

## Useful Resources

- [Snyk CLI Documentation](https://docs.snyk.io/snyk-cli)
- [GitHub Actions Documentation](https://docs.snyk.io/integrations/ci-cd-integrations/github-actions-integration)
- [Vulnerability Database](https://snyk.io/vuln)
- [Snyk Advisor](https://snyk.io/advisor/)
- [Security Training](https://learn.snyk.io/)

---

_Keep this reference handy while working with Snyk security scanning!_
