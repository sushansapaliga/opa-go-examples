# Latency Benchmark

This project benchmarks the performance of OPA policy evaluation in Go.

## Overview
It loads a large dataset from `policy/data.json` and evaluates the `data.authz.abilities` rule against various input scenarios (test cases).

## Files
- `main.go`: The benchmark runner. It loads the policy, runs test cases, and calculates average latency.
- `policy/`: Contains the Rego policy and data.
    - `data.json`: The data used for evaluation.

## Usage
Run the benchmark using:
```bash
go run main.go
```
It will output the results for each test case, including pass/fail status and average latency.
