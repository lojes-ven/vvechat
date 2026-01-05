import csv
import matplotlib.pyplot as plt
import os
import sys

def plot_lru(out_dir):
    csv_file = os.path.join(out_dir, "lru_bench_results.csv")
    if not os.path.exists(csv_file):
        print("No LRU results found.")
        return
        
    data = {'uniform': {'x': [], 'y': []}, 'zipf': {'x': [], 'y': []}}
    
    with open(csv_file, 'r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            w = row['workload']
            c = int(row['capacity'])
            hr = float(row['hit_rate'])
            
            if w in data:
                data[w]['x'].append(c)
                data[w]['y'].append(hr)
                
    plt.figure(figsize=(10, 6))
    
    if data['uniform']['x']:
        plt.plot(data['uniform']['x'], data['uniform']['y'], label='Uniform', marker='o')
        
    if data['zipf']['x']:
        plt.plot(data['zipf']['x'], data['zipf']['y'], label='Zipf (Hotspot)', marker='x')
        
    plt.xlabel('Cache Capacity')
    plt.ylabel('Hit Rate')
    plt.title('LRU Cache Hit Rate vs Capacity')
    plt.legend()
    plt.grid(True)
    plt.xscale('log')
    
    out_png = os.path.join(out_dir, "lru_bench_hitrate.png")
    plt.savefig(out_png)
    print(f"Plot saved to {out_png}")

if __name__ == "__main__":
    out_dir = "d:/mmy/vvechat/test/out"
    if len(sys.argv) > 1:
        out_dir = sys.argv[1]
    plot_lru(out_dir)
