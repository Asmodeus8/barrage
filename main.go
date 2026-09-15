package main

import (
 "flag"
 "fmt"
 "net/http"
 "sort"
 "sync"
 "sync/atomic"
 "time"
)

func main() {
 target := flag.String("url", "http://localhost:8080", "target URL")
 requests := flag.Int("n", 100, "number of requests")
 concurrency := flag.Int("c", 10, "concurrent workers")
 flag.Parse()
 if *requests < 1 || *concurrency < 1 { panic("n and c must be positive") }
 jobs := make(chan struct{})
 var ok, failed int64
 var mu sync.Mutex
 latencies := make([]time.Duration, 0, *requests)
 client := &http.Client{Timeout: 10 * time.Second}
 start := time.Now()
 var wg sync.WaitGroup
 for i := 0; i < *concurrency; i++ {
  wg.Add(1)
  go func() {
   defer wg.Done()
   for range jobs {
    began := time.Now()
    resp, err := client.Get(*target)
    elapsed := time.Since(began)
    mu.Lock(); latencies = append(latencies, elapsed); mu.Unlock()
    if err != nil { atomic.AddInt64(&failed, 1); continue }
    resp.Body.Close()
    if resp.StatusCode >= 200 && resp.StatusCode < 400 { atomic.AddInt64(&ok, 1) } else { atomic.AddInt64(&failed, 1) }
   }
  }()
 }
 for i := 0; i < *requests; i++ { jobs <- struct{}{} }
 close(jobs); wg.Wait()
 total := time.Since(start)
 sort.Slice(latencies, func(i,j int) bool { return latencies[i] < latencies[j] })
 percentile := func(p float64) time.Duration { if len(latencies)==0{return 0}; idx:=int(float64(len(latencies)-1)*p); return latencies[idx] }
 fmt.Printf("Barrage results\nrequests=%d success=%d failed=%d duration=%s throughput=%.2f req/s\np50=%s p95=%s p99=%s\n", *requests, ok, failed, total, float64(*requests)/total.Seconds(), percentile(.50), percentile(.95), percentile(.99))
}
