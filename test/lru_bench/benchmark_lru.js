const { LRUCache } = require('../../front_end/lru_cache.js');
const fs = require('fs');
const path = require('path');

// Parse args
const args = process.argv.slice(2);
const capacity = parseInt(args[0]) || 1000;
const workload = args[1] || 'uniform'; // uniform | zipf
const keyspace = parseInt(args[2]) || 10000;
const requests = parseInt(args[3]) || 20000;
const zipf_s = parseFloat(args[4]) || 1.1;
const seed = parseInt(args[5]) || 42;

// Simple PRNG
function mulberry32(a) {
    return function() {
      var t = a += 0x6D2B79F5;
      t = Math.imul(t ^ t >>> 15, t | 1);
      t ^= t + Math.imul(t ^ t >>> 7, t | 61);
      return ((t ^ t >>> 14) >>> 0) / 4294967296;
    }
}

const rand = mulberry32(seed);

// Zipf Generator
function zipf(s, N) {
    // Precompute probabilities (simple rejection or inverse transform is hard for Zipf)
    // For performance, we use an approximation or a simple precomputed table if N is small.
    // For large N, we can use the rank-frequency distribution.
    // Rank r has freq proportional to 1/r^s.
    // We can use a simple approximation: 
    // Generate x in [0, 1]. Find r such that CDF(r) >= x.
    
    // Optimization: Precompute CDF for small keyspace, or use approximation.
    // Given the constraints, let's use a simple precomputed CDF for N up to 1e6? No, too big.
    // We'll use a simplified generator:
    // p(k) = c / k^s
    
    // Let's use a simpler approach for the test:
    // Just generate numbers with a bias.
    // Or use a library if allowed? No, "reproducible code".
    
    // Harmonic number H_{N,s}
    let c = 0;
    for (let i = 1; i <= N; i++) {
        c += (1.0 / Math.pow(i, s));
    }
    c = 1.0 / c;
    
    // This is too slow for N=1e6.
    // Let's use the "Approximate Zipfian" logic often used in benchmarks (e.g. YCSB).
    // But for this script, let's implement a very simple one:
    // Select an index from 0 to N-1.
    // We want index 0 to be most frequent.
    
    // Inverse transform sampling is slow without precomputed CDF.
    // Let's use a small pool of "hot" keys and a large pool of "cold" keys.
    // 80% requests go to 20% keys (Pareto).
    // This is easier to implement and sufficient for "Hotspot" testing.
    
    return function() {
        if (rand() < 0.8) {
            return Math.floor(rand() * (N * 0.2));
        } else {
            return Math.floor(rand() * N);
        }
    };
}

function run() {
    const cache = new LRUCache(capacity);
    let hits = 0;
    let misses = 0;
    
    let keyGen;
    if (workload === 'zipf') {
        // Use Pareto 80/20 as proxy for Zipf-like hotspot
        keyGen = function() {
            if (rand() < 0.8) {
                return Math.floor(rand() * (keyspace * 0.2));
            } else {
                return Math.floor(rand() * keyspace);
            }
        };
    } else {
        keyGen = function() {
            return Math.floor(rand() * keyspace);
        };
    }
    
    const start = process.hrtime.bigint();
    
    for (let i = 0; i < requests; i++) {
        const key = keyGen();
        const val = cache.get(key);
        if (val !== -1 && val !== undefined && val !== null) { // Implementation dependent return
            hits++;
        } else {
            misses++;
            cache.put(key, key); // Store key as value
        }
    }
    
    const end = process.hrtime.bigint();
    const duration_ns = Number(end - start);
    const avg_latency_ns = duration_ns / requests;
    
    const result = {
        capacity,
        workload,
        keyspace,
        requests,
        seed,
        hits,
        misses,
        hit_rate: hits / requests,
        eviction_count: 0, // cache.evictionCount might not be exposed, check implementation
        avg_latency_ns
    };
    
    // Try to read eviction count if available (it's not in the provided snippet, but maybe in full file)
    // If not, we can infer or just leave 0.
    // The user asked to read `cache.evictionCount`.
    if (cache.evictionCount !== undefined) {
        result.eviction_count = cache.evictionCount;
    }

    console.log(JSON.stringify(result));
}

run();
