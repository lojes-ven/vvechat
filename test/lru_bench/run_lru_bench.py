import subprocess
import csv
import json
import os

def main():
    out_dir = "d:/mmy/vvechat/test/out"
    os.makedirs(out_dir, exist_ok=True)
    
    capacities = [100, 500, 1000, 5000, 10000]
    workloads = ['uniform', 'zipf']
    keyspace = 10000
    requests = 100000
    
    results = []
    
    print("Running LRU Benchmarks...")
    
    for w in workloads:
        for c in capacities:
            print(f"Running {w} with capacity {c}...")
            cmd = [
                'node', 
                'd:/mmy/vvechat/test/lru_bench/benchmark_lru.js',
                str(c),
                w,
                str(keyspace),
                str(requests),
                '1.1', # zipf_s
                '42'   # seed
            ]
            
            try:
                result = subprocess.run(cmd, capture_output=True, text=True, check=True)
                data = json.loads(result.stdout)
                results.append(data)
            except subprocess.CalledProcessError as e:
                print(f"Error running benchmark: {e.stderr}")
            except json.JSONDecodeError:
                print(f"Invalid JSON output: {result.stdout}")

    # Write CSV
    csv_file = os.path.join(out_dir, "lru_bench_results.csv")
    with open(csv_file, 'w', newline='') as f:
        writer = csv.DictWriter(f, fieldnames=results[0].keys())
        writer.writeheader()
        writer.writerows(results)
        
    print(f"Results written to {csv_file}")

if __name__ == "__main__":
    main()
