import subprocess
import csv
import json
import os
import argparse
import math


def _mean(values):
    return sum(values) / len(values) if values else 0.0


def _std_sample(values):
    # Sample standard deviation (N-1). If N<2, return 0.
    n = len(values)
    if n < 2:
        return 0.0
    m = _mean(values)
    var = sum((x - m) ** 2 for x in values) / (n - 1)
    return math.sqrt(var)


def _repo_paths():
    here = os.path.dirname(os.path.abspath(__file__))
    repo_root = os.path.abspath(os.path.join(here, "..", ".."))
    out_dir = os.path.join(repo_root, "test", "out")
    bench_js = os.path.join(here, "benchmark_lru.js")
    return repo_root, out_dir, bench_js

def main():
    parser = argparse.ArgumentParser(description="Run LRU cache micro-benchmarks")
    parser.add_argument("--repeats", type=int, default=5, help="repeat each config N times")
    parser.add_argument("--warmup", type=int, default=20000, help="warm-up requests per run (excluded from metrics)")
    parser.add_argument("--requests", type=int, default=100000, help="measured requests per run")
    parser.add_argument("--keyspace", type=int, default=10000, help="unique keyspace size")
    parser.add_argument("--seed", type=int, default=42, help="PRNG seed")
    args = parser.parse_args()

    _, out_dir, bench_js = _repo_paths()
    os.makedirs(out_dir, exist_ok=True)
    
    capacities = [100, 500, 1000, 5000, 10000]
    workloads = ['uniform', 'zipf']
    keyspace = args.keyspace
    requests = args.requests
    zipf_s_values = [0.8, 1.0, 1.2, 1.5]  # 不同的 Zipf 参数
    repeats = max(1, args.repeats)
    warmup = max(0, args.warmup)
    seed = args.seed
    
    results = []
    runs_detail = []
    
    print("=" * 60)
    print("Running LRU Benchmarks...")
    print("=" * 60)
    
    for w in workloads:
        s_list = zipf_s_values if w == 'zipf' else [1.1]
        
        for s in s_list:
            for c in capacities:
                print(f"Running {w} (s={s}) with capacity {c}...")
                per_run = []

                for _ in range(repeats):
                    cmd = [
                        'node',
                        bench_js,
                        str(c),
                        w,
                        str(keyspace),
                        str(requests),
                        str(s),
                        str(seed),
                        str(warmup),
                    ]
                
                    try:
                        result = subprocess.run(cmd, capture_output=True, text=True, check=True)
                        data = json.loads(result.stdout)
                        per_run.append(data)
                    except subprocess.CalledProcessError as e:
                        print(f"Error running benchmark: {e.stderr}")
                    except json.JSONDecodeError:
                        print(f"Invalid JSON output: {result.stdout}")

                if not per_run:
                    continue

                hit_rates = [r.get('hit_rate', 0.0) for r in per_run]
                throughputs = [r.get('throughput_ops_per_sec', 0.0) for r in per_run]
                latencies = [r.get('avg_latency_ns', 0.0) for r in per_run]
                evictions = [r.get('eviction_count', 0) for r in per_run]

                mean_hit_rate = _mean(hit_rates)
                std_hit_rate = _std_sample(hit_rates)
                mean_throughput = _mean(throughputs)
                std_throughput = _std_sample(throughputs)
                mean_latency = _mean(latencies)
                std_latency = _std_sample(latencies)
                mean_eviction = _mean(evictions)
                std_eviction = _std_sample(evictions)

                agg = {
                    'capacity': c,
                    'workload': w,
                    'keyspace': keyspace,
                    'requests': requests,
                    'warmup_requests': warmup,
                    'zipf_s': s if w == 'zipf' else None,
                    'seed': seed,
                    'repeats': repeats,
                    # Backward-compatible columns
                    'hits': None,
                    'misses': None,
                    'hit_rate': mean_hit_rate,
                    'eviction_count': mean_eviction,
                    'avg_latency_ns': mean_latency,
                    'throughput_ops_per_sec': mean_throughput,
                    # New variability columns
                    'hit_rate_std': std_hit_rate,
                    'eviction_count_std': std_eviction,
                    'avg_latency_ns_std': std_latency,
                    'throughput_ops_per_sec_std': std_throughput,
                }

                results.append(agg)
                runs_detail.append({
                    'workload': w,
                    'capacity': c,
                    'zipf_s': s if w == 'zipf' else None,
                    'repeats': repeats,
                    'runs': per_run,
                })

                print(f"  -> Hit Rate: {mean_hit_rate:.2%} ± {std_hit_rate:.2%}, Throughput: {mean_throughput:.0f} ± {std_throughput:.0f} ops/s")

    # Write CSV
    csv_file = os.path.join(out_dir, "lru_bench_results.csv")
    if results:
        fieldnames = [
            'capacity', 'workload', 'keyspace', 'requests', 'warmup_requests', 'zipf_s', 'seed', 'repeats',
            'hits', 'misses', 'hit_rate', 'hit_rate_std',
            'eviction_count', 'eviction_count_std',
            'avg_latency_ns', 'avg_latency_ns_std',
            'throughput_ops_per_sec', 'throughput_ops_per_sec_std',
        ]
        with open(csv_file, 'w', newline='') as f:
            writer = csv.DictWriter(f, fieldnames=fieldnames)
            writer.writeheader()
            writer.writerows(results)
        
    # Write JSON for detailed analysis
    json_file = os.path.join(out_dir, "lru_bench_results.json")
    with open(json_file, 'w') as f:
        json.dump(results, f, indent=2)

    runs_file = os.path.join(out_dir, "lru_bench_results_runs.json")
    with open(runs_file, 'w') as f:
        json.dump(runs_detail, f, indent=2)

    # Print Summary
    print("\n" + "=" * 60)
    print("SUMMARY")
    print("=" * 60)
    
    # Group by workload
    for w in workloads:
        print(f"\n{w.upper()} Workload:")
        w_results = [r for r in results if r['workload'] == w]
        for r in w_results:
            s_info = f" (s={r.get('zipf_s', '')})" if w == 'zipf' else ""
            print(f"  Capacity {r['capacity']:>5}{s_info}: Hit Rate = {r['hit_rate']:.2%} ± {r.get('hit_rate_std', 0.0):.2%}")
        
    print(f"\nResults written to {csv_file}")

if __name__ == "__main__":
    main()
