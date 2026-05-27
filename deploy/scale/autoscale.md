# Auto-Scaling Recipes — k3s and Nomad

## k3s

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: joblantern-web
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: joblantern-web
  minReplicas: 2
  maxReplicas: 40
  metrics:
    - type: Resource
      resource:
        name: cpu
        target: { type: Utilization, averageUtilization: 60 }
    - type: Pods
      pods:
        metric: { name: http_inflight_requests }
        target: { type: AverageValue, averageValue: "50" }
```

## Nomad

```hcl
job "joblantern-web" {
  group "web" {
    count = 4
    scaling {
      enabled = true
      min     = 2
      max     = 40
      policy {
        cooldown = "1m"
        evaluation_interval = "10s"
        check "cpu" {
          source = "prometheus"
          query  = "avg(rate(node_cpu_seconds_total{mode='user'}[1m]))"
          strategy "target-value" { target = 0.6 }
        }
      }
    }
  }
}
```

## Notes

- Verification jobs are queue-driven, so worker pools scale on queue length, not CPU.
- The agent layer's LLM calls are cost-bounded — cap worker concurrency per region.
