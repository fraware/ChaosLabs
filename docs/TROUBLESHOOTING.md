# Troubleshooting ChaosLabs

This document provides guidance for troubleshooting common issues encountered while setting up and running ChaosLabs. If you experience problems not covered here, please consider opening an issue on our [GitHub Issues](https://github.com/fraware/ChaosLabs/issues) page.

---

## Table of Contents

- [Kubernetes Deployment Issues](#kubernetes-deployment-issues)
- [Image Pull Errors](#image-pull-errors)
- [Experiment Failures](#experiment-failures)
- [Fault Injection Not Working](#fault-injection-not-working)
- [Prometheus & Grafana Metrics Issues](#prometheus--grafana-metrics-issues)
- [Distributed Tracing Problems](#distributed-tracing-problems)
- [Dashboard Not Displaying Correctly](#dashboard-not-displaying-correctly)
- [General Debugging Tips](#general-debugging-tips)

---

## Kubernetes Deployment Issues

### Symptoms:
- Deployments showing `0/1` READY pods.
- Pods in `Pending` or `CrashLoopBackOff` status.

### Steps to Resolve:
1. **Check Pod Status:**
   ```bash
   kubectl get pods -n chaoslabs

Use the following command to see detailed information:
   ```bash
   kubectl describe pod <pod-name> -n chaoslabs
  ```

2. **Review Events:**
Look at the events section in `kubectl describe` output to identify issues such as insufficient resources, scheduling conflicts, or missing secrets.

3. **Inspect Logs:**
   ```bash
kubectl logs <pod-name> -n chaoslabs
    ```
Look for error messages that may indicate what’s causing the pod to crash.

4. **Verify Configuration:**

- Ensure that environment variables (e.g., `AGENT_ENDPOINTS`) are set correctly.
- Confirm that necessary privileges (e.g., for running `tc` or `stress-ng`) are configured in your security context.

## Image Pull Errors

### Symptoms:
Pods in ErrImagePull or ImagePullBackOff status.

### Steps to Resolve:
1. **Verify Image Names and Tags:**
Ensure that the image names in your Kubernetes manifests match the images in your container registry.

2. **Push Local Images:**
If working with a local cluster (e.g., Minikube or Docker Desktop), either load your images into the cluster or push them to a registry that your cluster can access.

3. **Check Docker Hub/Registry Credentials:**
If using a private registry, confirm that your cluster is configured with the proper image pull secrets:
   ```bash
kubectl get secrets -n chaoslabs
    ```

Update your manifests with `imagePullSecrets` if necessary.

4. **Re-pull the Image Manually:**
You can test pulling the image manually on a node (if accessible) to ensure it’s available:
   ```bash
   docker pull your-dockerhub-username/chaos-controller:latest
    ```
## Experiment Failures

### Symptoms:
- Experiments do not start or complete.
- Controller logs show errors when dispatching experiments.

### Steps to Resolve:
1. **Check Controller Logs:**
   ```bash
   kubectl logs deployment/chaos-controller -n chaoslabs
    ```
   
Look for errors related to JSON parsing, scheduling, or communication with agents.

2. **Verify Experiment Request Format:**
Ensure that your JSON payloads conform to the expected schema. For example:
   ```json
    {
      "name": "CPU Stress Test",
      "description": "Runs CPU stress with 4 workers for 15s",
      "experiment_type": "cpu-stress",
      "duration": 15,
      "cpu_workers": 4
    }
    ```
   
3. **Examine Network Connectivity:**
If dispatching to multiple agents, verify that the endpoints specified in `AGENT_ENDPOINTS` are reachable from the controller pod.

## Fault Injection Not Working
### Symptoms:
- Agent logs do not show expected fault injection actions.
- No visible effect from experiments (e.g., network latency or CPU stress).

### Steps to Resolve:
1. **Run Commands Manually:**
Inside the agent container, try running a command like:
   ```bash
   tc qdisc show dev eth0
    ```
This helps verify if `tc` commands are being executed.

2. **Check Privileges:**
Many fault injection commands require elevated privileges. Confirm that the agent container is running in privileged mode.

3. **Review Command Output:**
Verify that commands like `stress-ng` produce output and do not immediately error out. Logs should indicate if a command failed.

4. **Time Duration:**
Ensure that the experiment duration is sufficient to observe the effects.

## Prometheus & Grafana Metrics Issues
### Symptoms:
- Prometheus is not scraping metrics.
- Grafana dashboards show no data.

### Steps to Resolve:
1. **Check Metrics Endpoints:**
Access the `/metrics` endpoint of your controller or agent directly:
   ```bash
   curl http://<pod-ip>:8080/metrics
    ```
Ensure that custom metrics (e.g., `controller_experiment_total`) are present.

2. **Verify Prometheus Annotations:**
Confirm that your Kubernetes manifests include annotations:
   ```yaml
annotations:
  prometheus.io/scrape: "true"
  prometheus.io/port: "8080"
    ```
3. **Prometheus Configuration:**
Check your Prometheus configuration to ensure it is scraping the correct namespaces and endpoints.

4. **Grafana Data Source:**
In Grafana, verify that the Prometheus data source is correctly configured and querying the expected metrics.

## Distributed Tracing Problems
### Symptoms:
- Jaeger UI does not show traces from ChaosLabs.
- Tracing spans are missing or incomplete.

### Steps to Resolve:
1. **Check Tracer Initialization:**
Ensure that both the controller and agent initialize the tracer properly. Review logs for any tracer-related errors.

2. **Verify Jaeger Collector Endpoint:**
Confirm that the Jaeger collector endpoint (e.g., `http://jaeger-collector:14268/api/traces`) is reachable from your pods.

3. **Sample Traces:**
Generate a few experiments and check if spans appear in the Jaeger UI. Adjust the sampling rate if needed (default is always sample).

4. **Environment Variables & Configuration:**
Make sure your tracing configuration matches your deployment (e.g., endpoints, sampling policies).

## Dashboard Not Displaying Correctly
### Symptoms:
- The dashboard web interface does not load or display data.
- Missing or incorrect experiment information.

### Steps to Resolve:
1. **Check Dashboard Logs:**

   ```bash
    kubectl logs deployment/chaos-dashboard -n chaoslabs
    ```
Look for errors in the Flask (or Node.js) logs.

2. **API Connectivity:**
Verify that the dashboard can successfully query the controller's `/experiments` endpoint.

3. **Browser Console:**
Open the browser’s developer console to check for JavaScript errors or failed API calls.

4. **Review Configuration:**
Ensure that the dashboard is configured to connect to the correct endpoints for metrics and experiments.

## General Debugging Tips
- **Use kubectl port-forward:**
Temporarily expose pod ports to your local machine to test endpoints:
   ```bash
   kubectl port-forward deployment/chaos-controller 8080:8080 -n chaoslabs
    ```
- **Incremental Testing:**
Test individual components (controller, agent, dashboard) locally before integrating them.

- **Log Verbosity:**
Increase log verbosity in your Go code if needed. Adding more detailed logs can help pinpoint where issues occur.

- **Review GitHub Issues/Community Discussions:**
Check our GitHub Issues and Discussions for similar issues reported by other users.

- **Documentation Updates:**
Refer to our User Guides & Tutorials and the FAQ sections for additional context.
