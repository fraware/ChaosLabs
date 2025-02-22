# ChaosLab Tutorial

Welcome to the ChaosLab Tutorial! This guide will walk you through setting up and running various chaos experiments using ChaosLab.

## Prerequisites

Before you begin, ensure you have:
- ChaosLab installed and running (see [Setup & Installation](../README.md#setup--installation))
- Access to a Kubernetes cluster or local environment via Docker Compose
- Basic familiarity with chaos engineering concepts

## Tutorial 1: Running a CPU Stress Test

### Step 1: Prepare the Experiment Request

Create a file named `cpu_stress.json` with the following content:

```json
{
  "name": "CPU Stress Test",
  "description": "Stress test with 4 CPU workers for 15 seconds",
  "experiment_type": "cpu-stress",
  "duration": 15,
  "cpu_workers": 4
}
```

### Step 2: Start the Experiment

Send the request to the controller:
```bash
curl -X POST -H "Content-Type: application/json" -d @cpu_stress.json http://localhost:8080/start
```

### Step 3: Monitor the Experiment

- Check controller and agent logs for status.
- Open the dashboard at http://localhost:5500 to view the experiment status.
- Visit the Prometheus metrics endpoint at `http://<controller-ip>:8080/metrics`.

## Tutorial 2: Simulating Network Latency

### Step 1: Prepare the Request

Create a file named `network_latency.json`:

```json
{
  "name": "Network Latency Test",
  "description": "Simulate 100ms network latency for 30 seconds",
  "experiment_type": "network-latency",
  "duration": 30,
  "delay_ms": 100
}
```

### Step 2: Start the Experiment
```bash
curl -X POST -H "Content-Type: application/json" -d @network_latency.json http://localhost:8080/start
```

### Step 3: Verify Execution

- Confirm in the agent logs that latency was applied and later removed.
- Check Grafana for network fault metrics.

## Tutorial 3: Memory Stress and Process Kill

### Memory Stress

Create `mem_stress.json`:
```json
{
  "name": "Memory Stress Test",
  "description": "Allocate 200 MB for 30 seconds",
  "experiment_type": "mem-stress",
  "duration": 30,
  "mem_size_mb": 200
}
```

Start the experiment:
```bash
curl -X POST -H "Content-Type: application/json" -d @mem_stress.json http://localhost:8080/start
```

### Process Kill
Create `process_kill.json`:
```json
{
  "name": "Process Kill Test",
  "description": "Kill a process matching 'go'",
  "experiment_type": "process-kill",
  "kill_process": "go"
}
```

Start the experiment:
```bash
curl -X POST -H "Content-Type: application/json" -d @process_kill.json http://localhost:8080/start
```

## Additional Scenarios
- **Scheduled Experiments:**
Include a `start_time` (RFC3339 format) in your JSON to schedule experiments.
- **Parallel Experiments:**
Set `"parallel": true` and specify `"agent_count"` to run experiments on multiple agents concurrently.

## Next Steps
- Experiment with different parameters and fault types.
- Monitor your experiments with the integrated Prometheus and Grafana dashboards.
- Share your findings or suggest improvements via GitHub Issues or Discussions.

For further assistance, consult the Troubleshooting Guide or contact the community.
