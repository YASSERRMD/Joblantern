// k6 load test for Joblantern.
// Target: 10k submissions/minute sustained on a documented hardware footprint.
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    sustained_submit: {
      executor: 'constant-arrival-rate',
      rate: 10000,
      timeUnit: '60s',
      duration: '10m',
      preAllocatedVUs: 200,
      maxVUs: 1000,
    },
  },
  thresholds: {
    'http_req_duration{name:submit}': ['p(99)<8000'],
  },
};

const submission = JSON.stringify({
  country: 'AE',
  industry: 'transport',
  body: 'URGENT driver wanted, AED 18000/mo. WhatsApp now.',
});

export default function () {
  const res = http.post(`${__ENV.TARGET}/api/v1/verifications`, submission, {
    headers: { 'content-type': 'application/json' },
    tags: { name: 'submit' },
  });
  check(res, { 'status is 2xx': r => r.status >= 200 && r.status < 300 });
  sleep(0.1);
}
