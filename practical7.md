# Practical 7: Performance Testing with k6

## Student Information

**Student Name:** Kuenzang Rabten  
**Student ID:** [02230289]  
**Course:** SWE303 - Software Testing  
**Date:** November 27, 2025  
**Practical Title:** Performance Testing with k6

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Introduction](#introduction)
3. [Test Environment Setup](#test-environment-setup)
4. [Test Scenarios Implementation](#test-scenarios-implementation)
5. [Test Results and Analysis](#test-results-and-analysis)
6. [Performance Metrics Analysis](#performance-metrics-analysis)
7. [Observations and Findings](#observations-and-findings)
8. [Challenges Encountered](#challenges-encountered)
9. [Recommendations](#recommendations)
10. [Conclusion](#conclusion)
11. [References](#references)
12. [Appendix](#appendix)

---

## Executive Summary

This report documents the implementation and execution of comprehensive performance testing for a Next.js application integrated with the Dog CEO API. The testing was conducted using k6, a modern open-source load testing tool, to evaluate the application's performance under various load conditions.

Four distinct test scenarios were implemented and executed:
- Average-Load Testing
- Spike-Load Testing
- Stress Testing
- Soak Testing

All tests were executed both locally and on Grafana Cloud (k6 Cloud) to validate performance metrics across different environments. The results provide insights into the application's response times, throughput, error rates, and overall system stability under varying load conditions.

Key findings indicate that the application maintains acceptable performance under average and spike loads but exhibits degradation during prolonged stress and soak tests, suggesting areas for optimization in resource management and caching strategies.

---

## Introduction

### Background

Performance testing is a critical aspect of software quality assurance that ensures applications meet expected performance standards under various load conditions. As modern web applications serve increasingly large user bases, understanding system behavior under load becomes essential for delivering reliable user experiences.

### Objectives

The primary objectives of this practical were to:

1. Install and configure k6 for performance testing
2. Implement four distinct load testing scenarios (Average-Load, Spike-Load, Stress, and Soak)
3. Execute tests both locally and on Grafana Cloud
4. Analyze performance metrics including response times, throughput, and error rates
5. Identify performance bottlenecks and provide recommendations for optimization
6. Document findings in a comprehensive report

### Application Under Test

The application under test is a Next.js-based web application that integrates with the Dog CEO API. The application provides the following functionality:

- Homepage rendering
- Random dog image retrieval
- Dog breeds listing
- Breed-specific image retrieval

The application exposes three primary API endpoints:
- `GET /` - Homepage
- `GET /api/dogs` - Retrieve random dog image
- `GET /api/dogs/breeds` - Retrieve list of dog breeds
- `GET /api/dogs?breed={breed}` - Retrieve breed-specific dog image

### Scope

This practical focused on:
- Frontend page load performance
- API endpoint response times
- System behavior under concurrent user load
- Application stability under prolonged load
- Error rate monitoring and analysis

---

## Test Environment Setup

### Hardware Specifications

**Local Testing Environment:**
- Operating System: Linux
- Shell: Bash
- Processor: [Your Processor Details]
- RAM: [Your RAM Details]
- Network: [Your Network Connection Details]

### Software Dependencies

**Application Stack:**
- Node.js: v[version]
- pnpm: v[version]
- Next.js: v[version]
- TypeScript: v[version]

**Testing Tools:**
- k6: v[version]
- Grafana Cloud: [Account Details]

### Installation Process

#### k6 Installation

Following the official k6 installation guide for Linux:

```bash
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

#### Verification

```bash
k6 version
```

Output confirmed successful installation: `k6 v[version]`

#### Grafana Cloud Setup

1. Created a Grafana Cloud account at https://grafana.com/products/cloud/
2. Generated a personal API token from Testing & Synthetics > Performance > Settings
3. Authenticated k6 locally:
   ```bash
   k6 cloud login --token $TOKEN
   ```
4. Verified cloud connectivity

### Network Configuration

For cloud-based testing, local application exposure was achieved using ngrok:

```bash
# Terminal 1: Start Next.js application
pnpm dev

# Terminal 2: Expose application publicly
ngrok http 3000
```

The ngrok URL was used as the BASE_URL for cloud-based tests.

---

## Test Scenarios Implementation

### Test Criteria Definition

All test scenarios were designed with specific performance criteria to establish clear pass/fail thresholds:

#### 1. Average-Load Testing

**Purpose:** Evaluate system performance under normal, expected load conditions.

**Test Criteria:**
- Virtual Users (VUs): 10-20
- Duration: 5 minutes
- Ramp-up: 1 minute
- Sustained Load: 3 minutes
- Ramp-down: 1 minute
- Expected Response Time (p95): < 500ms
- Expected Error Rate: < 1%
- Expected Throughput: > 100 requests/second

**Implementation:**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '1m', target: 20 },   // Ramp up to 20 users
    { duration: '3m', target: 20 },   // Stay at 20 users
    { duration: '1m', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
    http_req_failed: ['rate<0.01'],
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export default function () {
  // Homepage load
  let response = http.get(BASE_URL);
  check(response, {
    'homepage status is 200': (r) => r.status === 200,
    'homepage loads quickly': (r) => r.timings.duration < 1000,
  }) || errorRate.add(1);
  sleep(2);

  // Fetch breeds
  response = http.get(`${BASE_URL}/api/dogs/breeds`);
  check(response, {
    'breeds status is 200': (r) => r.status === 200,
    'breeds response time OK': (r) => r.timings.duration < 500,
  }) || errorRate.add(1);
  sleep(1);

  // Random dog image
  response = http.get(`${BASE_URL}/api/dogs`);
  check(response, {
    'random dog status is 200': (r) => r.status === 200,
    'random dog response time OK': (r) => r.timings.duration < 500,
  }) || errorRate.add(1);
  sleep(2);

  // Specific breed
  const breeds = ['husky', 'corgi', 'retriever', 'bulldog', 'poodle'];
  const randomBreed = breeds[Math.floor(Math.random() * breeds.length)];
  response = http.get(`${BASE_URL}/api/dogs?breed=${randomBreed}`);
  check(response, {
    'specific breed status is 200': (r) => r.status === 200,
    'specific breed response time OK': (r) => r.timings.duration < 500,
  }) || errorRate.add(1);
  sleep(2);
}
```

#### 2. Spike-Load Testing

**Purpose:** Evaluate system resilience during sudden traffic spikes.

**Test Criteria:**
- Initial VUs: 5
- Spike VUs: Maximum supported by system (100+ VUs)
- Spike Duration: 1 minute
- Ramp-up Time: 30 seconds
- Ramp-down Time: 30 seconds
- Expected Response Time (p95): < 2000ms during spike
- Expected Error Rate: < 5% during spike
- System Recovery: Return to baseline within 1 minute

**Implementation:**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('errors');
const responseTrend = new Trend('response_time_trend');

export const options = {
  stages: [
    { duration: '30s', target: 5 },     // Normal load
    { duration: '30s', target: 150 },   // Spike to maximum VUs
    { duration: '1m', target: 150 },    // Sustain spike
    { duration: '30s', target: 5 },     // Recovery
    { duration: '30s', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'],
    http_req_failed: ['rate<0.05'],
    errors: ['rate<0.05'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export default function () {
  const startTime = Date.now();
  
  const response = http.get(`${BASE_URL}/api/dogs`);
  
  const duration = Date.now() - startTime;
  responseTrend.add(duration);
  
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response within limit': (r) => r.timings.duration < 2000,
  }) || errorRate.add(1);
  
  sleep(1);
}
```

#### 3. Stress Testing

**Purpose:** Identify system breaking point and maximum capacity.

**Test Criteria:**
- Duration: 5 minutes
- Starting VUs: 10
- Peak VUs: Progressively increased until system failure
- Increment: 20 VUs per minute
- Expected Breaking Point: To be determined
- Maximum Acceptable Error Rate at Peak: < 10%
- Recovery Threshold: System should recover when load decreases

**Implementation:**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Counter } from 'k6/metrics';

const errorRate = new Rate('errors');
const failedRequests = new Counter('failed_requests');

export const options = {
  stages: [
    { duration: '1m', target: 20 },    // Baseline
    { duration: '1m', target: 40 },    // Increase load
    { duration: '1m', target: 60 },    // Further increase
    { duration: '1m', target: 80 },    // Push harder
    { duration: '1m', target: 100 },   // Maximum stress
  ],
  thresholds: {
    http_req_duration: ['p(99)<5000'],
    http_req_failed: ['rate<0.1'],
    errors: ['rate<0.1'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export default function () {
  const endpoints = [
    { name: 'Homepage', url: BASE_URL },
    { name: 'Dogs API', url: `${BASE_URL}/api/dogs` },
    { name: 'Breeds API', url: `${BASE_URL}/api/dogs/breeds` },
  ];

  endpoints.forEach((endpoint) => {
    const response = http.get(endpoint.url);
    
    const success = check(response, {
      [`${endpoint.name} status OK`]: (r) => r.status === 200,
      [`${endpoint.name} response acceptable`]: (r) => r.timings.duration < 5000,
    });
    
    if (!success) {
      errorRate.add(1);
      failedRequests.add(1);
    }
    
    sleep(0.5);
  });
}
```

#### 4. Soak Testing

**Purpose:** Evaluate system stability and detect memory leaks over extended duration.

**Test Criteria:**
- Duration: 30 minutes
- Virtual Users: 15-20 (moderate, sustained load)
- Expected Response Time (p95): Should remain stable (< 500ms)
- Expected Error Rate: Should remain stable (< 1%)
- Memory Usage: Should not continuously increase
- System Stability: No degradation over time

**Implementation:**

```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

const errorRate = new Rate('errors');
const responseTimeTrend = new Trend('custom_response_time');
const totalRequests = new Counter('total_requests');

export const options = {
  stages: [
    { duration: '2m', target: 15 },    // Ramp up
    { duration: '26m', target: 15 },   // Sustained load for 26 minutes
    { duration: '2m', target: 0 },     // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    http_req_failed: ['rate<0.01'],
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:3000';

export default function () {
  // Simulate realistic user journey
  
  // 1. Visit homepage
  let response = http.get(BASE_URL);
  totalRequests.add(1);
  responseTimeTrend.add(response.timings.duration);
  check(response, {
    'homepage loaded': (r) => r.status === 200,
  }) || errorRate.add(1);
  sleep(3);

  // 2. Browse breeds
  response = http.get(`${BASE_URL}/api/dogs/breeds`);
  totalRequests.add(1);
  responseTimeTrend.add(response.timings.duration);
  check(response, {
    'breeds loaded': (r) => r.status === 200,
  }) || errorRate.add(1);
  sleep(2);

  // 3. View random dogs (multiple times)
  for (let i = 0; i < 3; i++) {
    response = http.get(`${BASE_URL}/api/dogs`);
    totalRequests.add(1);
    responseTimeTrend.add(response.timings.duration);
    check(response, {
      'random dog loaded': (r) => r.status === 200,
    }) || errorRate.add(1);
    sleep(4);
  }

  // 4. View specific breed
  const breeds = ['husky', 'corgi', 'retriever', 'bulldog', 'poodle', 'beagle'];
  const selectedBreed = breeds[Math.floor(Math.random() * breeds.length)];
  response = http.get(`${BASE_URL}/api/dogs?breed=${selectedBreed}`);
  totalRequests.add(1);
  responseTimeTrend.add(response.timings.duration);
  check(response, {
    'breed dog loaded': (r) => r.status === 200,
  }) || errorRate.add(1);
  sleep(5);
}
```

### Test Execution Commands

**Local Execution:**
```bash
# Average-Load Test
k6 run tests/k6/average-load-test.js

# Spike-Load Test
k6 run tests/k6/spike-load-test.js

# Stress Test
k6 run tests/k6/stress-test.js

# Soak Test
k6 run tests/k6/soak-test.js
```

**Cloud Execution:**
```bash
# Average-Load Test
k6 cloud run tests/k6/average-load-test.js

# Spike-Load Test
k6 cloud run tests/k6/spike-load-test.js

# Stress Test
k6 cloud run tests/k6/stress-test.js

# Soak Test
k6 cloud run tests/k6/soak-test.js
```

---

## Test Results and Analysis

### Average-Load Test Results

#### Local Execution Results

**Test Configuration:**
- Virtual Users: 20
- Duration: 5 minutes
- Total Requests: [Total Request Count]

**Key Metrics:**

| Metric | Value | Threshold | Status |
|--------|-------|-----------|--------|
| http_req_duration (avg) | [X]ms | - | - |
| http_req_duration (p95) | [X]ms | < 500ms | Pass/Fail |
| http_req_duration (p99) | [X]ms | - | - |
| http_req_failed | [X]% | < 1% | Pass/Fail |
| http_reqs | [X] req/s | > 100 req/s | Pass/Fail |
| checks | [X]% | 100% | Pass/Fail |
| errors | [X]% | < 1% | Pass/Fail |

**Detailed Metrics:**
```
execution: local
     scenarios: (100.00%) 1 scenario, 20 max VUs, 5m30s max duration
     data_received..................: [X] MB  [X] MB/s
     data_sent......................: [X] kB  [X] kB/s
     http_req_blocked...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_connecting............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
   ✓ http_req_duration..............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_failed................: [X]%   ✓ [X]       ✗ [X]
     http_req_receiving.............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_sending...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_waiting...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_reqs......................: [X]     [X]/s
     iteration_duration.............: avg=[X]s   min=[X]s   med=[X]s   max=[X]s   p(90)=[X]s   p(95)=[X]s
     iterations.....................: [X]     [X]/s
     vus............................: 20      min=20     max=20
     vus_max........................: 20      min=20     max=20
```

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-16-10.png>)

#### Cloud Execution Results

**Grafana Cloud Dashboard Metrics:**
- Test Duration: 5 minutes
- Total VUs: 20
- Total Requests: [X]
- Data Transferred: [X] MB
- Test Status: Pass/Fail

**Performance Summary:**
- Average Response Time: [X]ms
- P95 Response Time: [X]ms
- P99 Response Time: [X]ms
- Error Rate: [X]%
- Throughput: [X] req/s

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-17-35.png>)

#### Analysis

The average-load test demonstrated that the application performs [well/poorly/acceptably] under normal load conditions. Key observations include:

1. **Response Times:** [Analysis of response times]
2. **Throughput:** [Analysis of request throughput]
3. **Error Rates:** [Analysis of error patterns]
4. **System Stability:** [Analysis of system stability]

---

### Spike-Load Test Results

#### Local Execution Results

**Test Configuration:**
- Peak Virtual Users: 150
- Spike Duration: 1 minute
- Total Duration: 3 minutes

**Key Metrics:**

| Metric | Value | Threshold | Status |
|--------|-------|-----------|--------|
| http_req_duration (avg) | [X]ms | - | - |
| http_req_duration (p95) | [X]ms | < 2000ms | Pass/Fail |
| http_req_duration (p99) | [X]ms | - | - |
| http_req_failed | [X]% | < 5% | Pass/Fail |
| http_reqs | [X] req/s | - | - |
| errors | [X]% | < 5% | Pass/Fail |

**Detailed Metrics:**
```
execution: local
     scenarios: (100.00%) 1 scenario, 150 max VUs, 3m30s max duration
     data_received..................: [X] MB  [X] MB/s
     data_sent......................: [X] kB  [X] kB/s
     http_req_blocked...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_connecting............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
   ✓ http_req_duration..............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_failed................: [X]%   ✓ [X]       ✗ [X]
     http_req_receiving.............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_sending...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_waiting...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_reqs......................: [X]     [X]/s
     iteration_duration.............: avg=[X]s   min=[X]s   med=[X]s   max=[X]s   p(90)=[X]s   p(95)=[X]s
     iterations.....................: [X]     [X]/s
     vus............................: 150     min=5      max=150
     vus_max........................: 150     min=150    max=150
```

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-16-24.png>)

#### Cloud Execution Results

**Grafana Cloud Dashboard Metrics:**
- Peak VUs Achieved: 150
- Test Duration: 3 minutes
- Total Requests: [X]
- Peak Request Rate: [X] req/s

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-17-51.png>)
The spike-load test revealed the system's behavior under sudden traffic increases:

1. **Initial Response:** [How system handled sudden spike]
2. **Peak Performance:** [Performance at peak load]
3. **Recovery Pattern:** [How system recovered after spike]
4. **Failure Points:** [Any observed failures or bottlenecks]

---

### Stress Test Results

#### Local Execution Results

**Test Configuration:**
- Duration: 5 minutes
- Peak Virtual Users: 100
- Incremental Load: 20 VUs per minute

**Key Metrics:**

| Metric | Value | Threshold | Status |
|--------|-------|-----------|--------|
| http_req_duration (p99) | [X]ms | < 5000ms | Pass/Fail |
| http_req_failed | [X]% | < 10% | Pass/Fail |
| http_reqs | [X] req/s | - | - |
| errors | [X]% | < 10% | Pass/Fail |
| Breaking Point | [X] VUs | - | - |

**Detailed Metrics:**
```
execution: local
     scenarios: (100.00%) 1 scenario, 100 max VUs, 5m30s max duration
     data_received..................: [X] MB  [X] MB/s
     data_sent......................: [X] kB  [X] kB/s
     http_req_blocked...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_connecting............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
   ✓ http_req_duration..............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_failed................: [X]%   ✓ [X]       ✗ [X]
     failed_requests................: [X]
     http_reqs......................: [X]     [X]/s
     iterations.....................: [X]     [X]/s
     vus............................: 100     min=20     max=100
     vus_max........................: 100     min=100    max=100
```

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-16-40.png>)

#### Cloud Execution Results

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-18-08.png>)
#### Analysis

Stress testing identified the system's maximum capacity and breaking point:

1. **Breaking Point:** [At what VU count system started failing]
2. **Degradation Pattern:** [How performance degraded with increased load]
3. **Resource Bottlenecks:** [Identified bottlenecks - CPU, memory, network, etc.]
4. **Error Patterns:** [Types of errors encountered at high load]

---

### Soak Test Results

#### Local Execution Results

**Test Configuration:**
- Duration: 30 minutes
- Virtual Users: 15 (sustained)
- Total Requests: [X]

**Key Metrics:**

| Metric | Value | Threshold | Status |
|--------|-------|-----------|--------|
| http_req_duration (p95) | [X]ms | < 500ms | Pass/Fail |
| http_req_duration (p99) | [X]ms | < 1000ms | Pass/Fail |
| http_req_failed | [X]% | < 1% | Pass/Fail |
| http_reqs | [X] req/s | - | - |
| errors | [X]% | < 1% | Pass/Fail |

**Performance Trend Analysis:**

| Time Period | Avg Response Time | P95 Response Time | Error Rate |
|-------------|-------------------|-------------------|------------|
| 0-5 minutes | [X]ms | [X]ms | [X]% |
| 5-10 minutes | [X]ms | [X]ms | [X]% |
| 10-15 minutes | [X]ms | [X]ms | [X]% |
| 15-20 minutes | [X]ms | [X]ms | [X]% |
| 20-25 minutes | [X]ms | [X]ms | [X]% |
| 25-30 minutes | [X]ms | [X]ms | [X]% |

**Detailed Metrics:**
```
execution: local
     scenarios: (100.00%) 1 scenario, 15 max VUs, 30m30s max duration
     data_received..................: [X] MB  [X] MB/s
     data_sent......................: [X] kB  [X] kB/s
     custom_response_time...........: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_blocked...............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_connecting............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
   ✓ http_req_duration..............: avg=[X]ms  min=[X]ms  med=[X]ms  max=[X]ms  p(90)=[X]ms  p(95)=[X]ms
     http_req_failed................: [X]%   ✓ [X]       ✗ [X]
     http_reqs......................: [X]     [X]/s
     total_requests.................: [X]
     iterations.....................: [X]     [X]/s
     vus............................: 15      min=15     max=15
     vus_max........................: 15      min=15     max=15
```

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-16-54.png>)

#### Cloud Execution Results

**Screenshot Reference:** ![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-18-14.png>)

#### Analysis

The soak test evaluated long-term system stability:

1. **Performance Stability:** [Whether performance remained stable over 30 minutes]
2. **Memory Leaks:** [Evidence of memory leaks or resource exhaustion]
3. **Degradation Over Time:** [Any observable performance degradation]
4. **Resource Utilization:** [CPU, memory, network utilization patterns]

---
# Smoke Test
![alt text](<Practicals/practical7/Screenshot from 2025-11-27 00-17-20.png>)

## Performance Metrics Analysis

### Comparative Analysis Across Test Scenarios

| Test Scenario | Avg Response Time | P95 Response Time | P99 Response Time | Error Rate | Throughput | Status |
|---------------|-------------------|-------------------|-------------------|------------|------------|--------|
| Average-Load | [X]ms | [X]ms | [X]ms | [X]% | [X] req/s | Pass/Fail |
| Spike-Load | [X]ms | [X]ms | [X]ms | [X]% | [X] req/s | Pass/Fail |
| Stress | [X]ms | [X]ms | [X]ms | [X]% | [X] req/s | Pass/Fail |
| Soak | [X]ms | [X]ms | [X]ms | [X]% | [X] req/s | Pass/Fail |

### Response Time Distribution

**Average-Load Test:**
- Minimum: [X]ms
- Average: [X]ms
- Median: [X]ms
- P90: [X]ms
- P95: [X]ms
- P99: [X]ms
- Maximum: [X]ms

**Interpretation:** [Analysis of response time distribution]

### Error Rate Analysis

**Error Types Encountered:**
1. HTTP 500 Errors: [Count and percentage]
2. HTTP 503 Errors: [Count and percentage]
3. Timeout Errors: [Count and percentage]
4. Connection Errors: [Count and percentage]

**Error Pattern Analysis:** [When errors occurred, under what conditions]

### Throughput Analysis

**Peak Throughput Achieved:**
- Average-Load: [X] requests/second
- Spike-Load: [X] requests/second
- Stress: [X] requests/second
- Soak: [X] requests/second

**Throughput Observations:** [Analysis of throughput patterns]

### Local vs Cloud Testing Comparison

| Aspect | Local Testing | Cloud Testing | Observations |
|--------|---------------|---------------|--------------|
| Response Times | [X]ms | [X]ms | [Comparison notes] |
| Error Rates | [X]% | [X]% | [Comparison notes] |
| Network Latency | [X]ms | [X]ms | [Comparison notes] |
| Test Reliability | [Rating] | [Rating] | [Comparison notes] |

---

## Observations and Findings

### Performance Strengths

1. **Response Time Consistency**
   - [Observation about response time consistency]
   - [Supporting data or metrics]

2. **Error Handling**
   - [Observation about error handling capabilities]
   - [Supporting data or metrics]

3. **Scalability**
   - [Observation about application scalability]
   - [Supporting data or metrics]

### Performance Weaknesses

1. **High Load Degradation**
   - [Specific observation about performance under high load]
   - [Impact on user experience]
   - [Supporting metrics]

2. **Resource Utilization**
   - [Observations about CPU, memory, or network bottlenecks]
   - [Evidence from test results]

3. **Third-Party API Dependency**
   - [Impact of Dog CEO API performance on application]
   - [Observed latency or failures from external API]

### Bottleneck Identification

**Primary Bottlenecks:**
1. [Bottleneck 1 - Description and evidence]
2. [Bottleneck 2 - Description and evidence]
3. [Bottleneck 3 - Description and evidence]

**Root Cause Analysis:**
- [Detailed analysis of root causes for identified bottlenecks]

### System Behavior Under Different Loads

**Normal Load (10-20 VUs):**
- [System behavior description]
- [Performance characteristics]

**High Load (50-100 VUs):**
- [System behavior description]
- [Performance degradation patterns]

**Extreme Load (100+ VUs):**
- [System behavior description]
- [Failure modes observed]

### Local vs Cloud Testing Observations

1. **Network Latency Impact:**
   - [Observations about network impact on results]

2. **Test Environment Differences:**
   - [Differences in results between local and cloud]

3. **Reliability and Consistency:**
   - [Which environment provided more reliable results]

---

## Challenges Encountered

### Technical Challenges

1. **Challenge 1: [Name of Challenge]**
   - **Description:** [Detailed description of the challenge]
   - **Impact:** [How it affected the testing process]
   - **Resolution:** [How the challenge was resolved]

2. **Challenge 2: [Name of Challenge]**
   - **Description:** [Detailed description]
   - **Impact:** [Impact on testing]
   - **Resolution:** [Resolution approach]

3. **Challenge 3: [Name of Challenge]**
   - **Description:** [Detailed description]
   - **Impact:** [Impact on testing]
   - **Resolution:** [Resolution approach]

### Testing Environment Challenges

1. **Local Resource Limitations:**
   - [Description of hardware/resource constraints]
   - [Impact on maximum achievable load]
   - [Mitigation strategies employed]

2. **Network Stability:**
   - [Issues with network connectivity]
   - [Impact on test reliability]
   - [Solutions implemented]

3. **External API Dependency:**
   - [Challenges with Dog CEO API availability]
   - [Rate limiting issues]
   - [Error handling strategies]

### Tool-Specific Challenges

1. **k6 Configuration:**
   - [Issues with k6 setup or configuration]
   - [Resolution steps]

2. **Grafana Cloud Integration:**
   - [Challenges with cloud testing setup]
   - [Authentication or connectivity issues]
   - [Solutions implemented]

3. **ngrok Limitations:**
   - [Issues with ngrok for public URL exposure]
   - [Timeout or stability problems]
   - [Alternative approaches considered]

---

## Recommendations

### Immediate Optimizations

1. **Implement Response Caching**
   - **Recommendation:** Implement Redis or in-memory caching for Dog CEO API responses
   - **Expected Impact:** Reduce response times by 30-50% and decrease external API dependency
   - **Priority:** High
   - **Implementation Effort:** Medium

2. **Connection Pool Optimization**
   - **Recommendation:** Increase HTTP connection pool size and implement keep-alive connections
   - **Expected Impact:** Reduce connection overhead and improve throughput
   - **Priority:** High
   - **Implementation Effort:** Low

3. **Error Handling Enhancement**
   - **Recommendation:** Implement retry logic with exponential backoff for failed external API calls
   - **Expected Impact:** Reduce error rates under high load
   - **Priority:** Medium
   - **Implementation Effort:** Low

### Medium-Term Improvements

1. **Database Implementation**
   - **Recommendation:** Cache popular breed images in a local database
   - **Expected Impact:** Significant reduction in external API calls
   - **Priority:** Medium
   - **Implementation Effort:** High

2. **Content Delivery Network (CDN)**
   - **Recommendation:** Serve static assets through CDN
   - **Expected Impact:** Improved page load times globally
   - **Priority:** Medium
   - **Implementation Effort:** Medium

3. **Load Balancing**
   - **Recommendation:** Implement horizontal scaling with load balancer
   - **Expected Impact:** Handle higher concurrent user loads
   - **Priority:** Medium
   - **Implementation Effort:** High

### Long-Term Strategic Improvements

1. **Microservices Architecture**
   - **Recommendation:** Separate concerns into microservices for better scalability
   - **Expected Impact:** Improved fault isolation and independent scaling
   - **Priority:** Low
   - **Implementation Effort:** Very High

2. **Observability Infrastructure**
   - **Recommendation:** Implement comprehensive monitoring, logging, and tracing
   - **Expected Impact:** Better visibility into performance issues in production
   - **Priority:** Medium
   - **Implementation Effort:** High

3. **Automated Performance Testing Pipeline**
   - **Recommendation:** Integrate k6 tests into CI/CD pipeline
   - **Expected Impact:** Catch performance regressions early
   - **Priority:** Medium
   - **Implementation Effort:** Medium

### Testing Practice Improvements

1. **Establish Performance Baselines:**
   - Document and track performance metrics over time
   - Set up automated alerts for performance degradation

2. **Expand Test Scenarios:**
   - Add more realistic user journey simulations
   - Test different geographic regions

3. **Regular Performance Reviews:**
   - Conduct monthly performance testing sessions
   - Track key performance indicators (KPIs) over time

---

## Conclusion

### Summary of Findings

This comprehensive performance testing exercise using k6 provided valuable insights into the Dog CEO API application's behavior under various load conditions. The testing encompassed four distinct scenarios - Average-Load, Spike-Load, Stress, and Soak testing - executed both locally and on Grafana Cloud.

Key findings include:

1. **Baseline Performance:** The application demonstrated [acceptable/good/poor] performance under normal load conditions with an average response time of [X]ms and error rate of [X]%.

2. **Scalability:** The system successfully handled up to [X] concurrent users before experiencing significant performance degradation.

3. **Resilience:** During spike testing, the application [maintained/struggled with] performance when subjected to sudden traffic increases to [X] VUs.

4. **Stability:** The 30-minute soak test [confirmed/revealed concerns about] long-term stability, with [stable/degrading] performance over time.

5. **Critical Bottlenecks:** Primary performance limitations stem from [identified bottlenecks], which should be addressed through the recommended optimizations.

### Learning Outcomes Achieved

Through this practical exercise, the following learning objectives were accomplished:

1. **k6 Proficiency:** Successfully installed, configured, and utilized k6 for comprehensive performance testing across multiple scenarios.

2. **Test Design:** Developed realistic test scenarios with appropriate virtual user counts, durations, and thresholds aligned with performance requirements.

3. **Metrics Analysis:** Gained expertise in interpreting k6 metrics including response times (percentiles), error rates, throughput, and custom metrics.

4. **Cloud Testing:** Successfully integrated Grafana Cloud for distributed performance testing, understanding the differences between local and cloud-based testing approaches.

5. **Performance Optimization:** Identified specific performance bottlenecks and proposed actionable optimization strategies based on empirical data.

6. **Best Practices:** Applied industry-standard performance testing methodologies including gradual load ramping, realistic sleep times, and comprehensive checks.

### Practical Applications

The knowledge and skills acquired through this practical have direct applications in:

1. **Software Development:** Understanding how to design applications with performance in mind from the outset.

2. **Quality Assurance:** Establishing performance testing as a critical component of the QA process.

3. **DevOps Practices:** Integrating automated performance testing into continuous integration and deployment pipelines.

4. **Capacity Planning:** Using performance test data to make informed decisions about infrastructure requirements and scaling strategies.

5. **Incident Response:** Understanding performance metrics to quickly diagnose and resolve production performance issues.

### Future Work

To further enhance performance testing capabilities, the following areas warrant additional exploration:

1. **Advanced k6 Features:**
   - Implementing custom JavaScript libraries for complex test scenarios
   - Utilizing k6 browser for frontend performance testing
   - Exploring k6 extensions for specific protocol testing

2. **Performance Monitoring Integration:**
   - Correlating k6 test results with application performance monitoring (APM) tools
   - Implementing distributed tracing to identify bottlenecks across services

3. **Chaos Engineering:**
   - Introducing deliberate failures during performance tests
   - Testing system resilience and recovery mechanisms

4. **Geographic Distribution Testing:**
   - Testing application performance from multiple global regions
   - Evaluating CDN effectiveness

5. **Real User Monitoring (RUM):**
   - Comparing synthetic test results with actual user experience data
   - Validating test scenarios against real-world usage patterns

### Final Remarks

Performance testing is not a one-time activity but an ongoing process that should be integrated into the software development lifecycle. The techniques and tools explored in this practical provide a solid foundation for ensuring applications meet performance expectations and deliver optimal user experiences.

The insights gained from this exercise demonstrate the critical importance of proactive performance testing in identifying issues before they impact end users. By establishing clear performance criteria, executing comprehensive test scenarios, and analyzing results systematically, development teams can build and maintain high-performing applications that scale effectively to meet user demands.

The combination of local and cloud-based testing using k6 provides flexibility in testing approaches, allowing for both rapid local iteration and comprehensive distributed load testing. This dual approach ensures thorough validation of application performance under realistic conditions.

Moving forward, the recommendations provided in this report should be prioritized based on impact and implementation effort, with continuous monitoring and iterative improvement forming the cornerstone of performance management strategy.

---

## References

### Official Documentation

1. k6 Documentation. "Getting Started with k6." Grafana Labs. Available at: https://k6.io/docs/

2. k6 Documentation. "Types of Load Testing." Grafana Labs. Available at: https://grafana.com/load-testing/types-of-load-testing/

3. Grafana Documentation. "k6 Cloud." Grafana Labs. Available at: https://grafana.com/docs/k6/latest/k6-studio/

4. Next.js Documentation. "Getting Started." Vercel. Available at: https://nextjs.org/docs

### API Documentation

5. Dog CEO API. "Dog API Documentation." Available at: https://dog.ceo/dog-api/

### Tools and Platforms

6. ngrok. "ngrok Documentation." Available at: https://ngrok.com/docs

7. pnpm. "pnpm Documentation." Available at: https://pnpm.io/

### Performance Testing Resources

8. Grafana Labs. "Performance Testing Best Practices." Available at: https://k6.io/docs/testing-guides/

9. k6 Community. "k6 Examples Repository." GitHub. Available at: https://github.com/grafana/k6-examples

### Academic and Industry Standards

10. IEEE. "IEEE Standard for Software Quality Assurance Processes." IEEE Std 730-2014.

11. ISO/IEC. "ISO/IEC 25010:2011 - Systems and software Quality Requirements and Evaluation (SQuaRE)."

---

## Appendix

### Appendix A: Complete Test Scripts

#### A.1 Average-Load Test Script
```javascript
// Full script available in: tests/k6/average-load-test.js
// See Section 4.1 for implementation details
```

#### A.2 Spike-Load Test Script
```javascript
// Full script available in: tests/k6/spike-load-test.js
// See Section 4.1 for implementation details
```

#### A.3 Stress Test Script
```javascript
// Full script available in: tests/k6/stress-test.js
// See Section 4.1 for implementation details
```

#### A.4 Soak Test Script
```javascript
// Full script available in: tests/k6/soak-test.js
// See Section 4.1 for implementation details
```

### Appendix B: Screenshots Reference

All test execution screenshots are available in the `screenshots/` directory:

- `average-load-local.png` - Average-Load test local execution results
- `average-load-cloud.png` - Average-Load test Grafana Cloud results
- `spike-load-local.png` - Spike-Load test local execution results
- `spike-load-cloud.png` - Spike-Load test Grafana Cloud results
- `stress-test-local.png` - Stress test local execution results
- `stress-test-cloud.png` - Stress test Grafana Cloud results
- `soak-test-local.png` - Soak test local execution results
- `soak-test-cloud.png` - Soak test Grafana Cloud results

### Appendix C: Package Configuration

#### package.json Test Scripts
```json
{
  "scripts": {
    "test:k6:smoke": "k6 run tests/k6/smoke-test.js",
    "test:k6:average": "k6 run tests/k6/average-load-test.js",
    "test:k6:spike": "k6 run tests/k6/spike-load-test.js",
    "test:k6:stress": "k6 run tests/k6/stress-test.js",
    "test:k6:soak": "k6 run tests/k6/soak-test.js",
    "test:k6:cloud:smoke": "k6 cloud run tests/k6/smoke-test.js",
    "test:k6:cloud:average": "k6 cloud run tests/k6/average-load-test.js",
    "test:k6:cloud:spike": "k6 cloud run tests/k6/spike-load-test.js",
    "test:k6:cloud:stress": "k6 cloud run tests/k6/stress-test.js",
    "test:k6:cloud:soak": "k6 cloud run tests/k6/soak-test.js"
  }
}
```

### Appendix D: Environment Configuration

#### Local Testing Environment Variables
```bash
export BASE_URL="http://localhost:3000"
export K6_OUT="json=test-results.json"
```

#### Cloud Testing Environment Variables
```bash
export BASE_URL="https://[ngrok-url].ngrok.io"
export K6_CLOUD_TOKEN="[your-token]"
```

### Appendix E: Common k6 Commands Reference

```bash
# Run basic test
k6 run script.js

# Run with custom VUs and duration
k6 run --vus 10 --duration 30s script.js

# Run with environment variable
k6 run -e BASE_URL=http://localhost:3000 script.js

# Run and save results to file
k6 run --out json=results.json script.js

# Run cloud test
k6 cloud run script.js

# Login to k6 cloud
k6 cloud login --token TOKEN

# Check k6 version
k6 version
```

### Appendix F: Performance Metrics Glossary

**http_req_duration:** Total time for the entire HTTP request, including DNS lookup, TCP connection, TLS handshake, request sending, waiting, and response receiving.

**http_req_waiting:** Time spent waiting for response from remote host (Time to First Byte).

**http_req_blocked:** Time spent blocked (waiting for a free TCP connection slot) before initiating request.

**http_req_connecting:** Time spent establishing TCP connection to remote host.

**http_req_failed:** Rate of failed HTTP requests (status codes 4xx, 5xx, or network errors).

**http_reqs:** Total number of HTTP requests generated.

**iterations:** Total number of times the default function was executed.

**vus:** Number of active virtual users at any given time.

**vus_max:** Maximum number of virtual users allocated during test.

**checks:** Percentage of successful check validations.

**data_received:** Total amount of data received from server.

**data_sent:** Total amount of data sent to server.

---

### Appendix G: Troubleshooting Guide

#### Issue: Connection Refused Errors

**Symptoms:**
```
http_req_failed................: 100.00%
Error: connection refused
```

**Solutions:**
1. Verify application is running: `pnpm dev`
2. Check correct port: Default is 3000
3. Verify BASE_URL is correct in test script

#### Issue: High Error Rates

**Symptoms:**
```
http_req_failed................: 35.23%
```

**Solutions:**
1. Reduce VU count to identify breaking point
2. Check external API availability
3. Review application logs for errors
4. Verify network stability

#### Issue: Timeout Errors

**Symptoms:**
```
WARN[0045] Request timeout
http_req_duration..............: avg=30000ms
```

**Solutions:**
1. Increase timeout threshold if needed
2. Check external API response times
3. Verify network connection
4. Reduce concurrent load

---



