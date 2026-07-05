package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type result struct {
	status int
	ms     float64
	err    bool
}

type scenario struct {
	name   string
	method string
	path   string
	body   []byte
	auth   bool
}

func main() {
	baseURL := flag.String("url", envOr("PERF_BASE_URL", "http://localhost:8000"), "API base URL")
	requests := flag.Int("n", 200, "requests per scenario")
	concurrency := flag.Int("c", 20, "concurrent workers")
	flag.Parse()

	client := &http.Client{Timeout: 30 * time.Second}

	if !healthCheck(client, *baseURL) {
		fmt.Fprintf(os.Stderr, "API not reachable at %s\n", *baseURL)
		os.Exit(1)
	}

	token, err := setupAuth(client, *baseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "auth setup failed: %v\n", err)
		os.Exit(1)
	}

	productID, err := setupProduct(client, *baseURL, token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "product setup failed: %v\n", err)
		os.Exit(1)
	}

	scenarios := []scenario{
		{name: "List products (public)", method: http.MethodGet, path: "/users/productview"},
		{name: "Search products", method: http.MethodGet, path: "/users/search?name=Laptop"},
		{name: "Get cart (auth)", method: http.MethodGet, path: "/usercart", auth: true},
		{name: "Add to cart (auth)", method: http.MethodPost, path: "/addtocart?id=" + productID, auth: true},
	}

	fmt.Println("============================================================")
	fmt.Println("Ecommerce API Performance Test")
	fmt.Printf("Target:      %s\n", *baseURL)
	fmt.Printf("Per scenario: %d requests, %d concurrent\n", *requests, *concurrency)
	fmt.Println("============================================================")

	for _, sc := range scenarios {
		runScenario(client, *baseURL, token, sc, *requests, *concurrency)
		fmt.Println()
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func healthCheck(client *http.Client, baseURL string) bool {
	resp, err := client.Get(baseURL + "/users/productview")
	if err != nil {
		return false
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func setupAuth(client *http.Client, baseURL string) (string, error) {
	suffix := time.Now().UnixNano()
	email := fmt.Sprintf("perf-%d@example.com", suffix)
	phone := fmt.Sprintf("8%09d", suffix%1_000_000_000)
	password := "password123"

	signupBody, _ := json.Marshal(map[string]string{
		"first_name": "Perf",
		"last_name":  "Test",
		"email":      email,
		"password":   password,
		"phone":      phone,
	})
	if _, err := doRequest(client, baseURL, http.MethodPost, "/users/signup", signupBody, ""); err != nil {
		return "", err
	}

	loginBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	resp, err := doRequest(client, baseURL, http.MethodPost, "/users/login", loginBody, "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var login map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&login); err != nil {
		return "", err
	}
	token, _ := login["token"].(string)
	if token == "" {
		return "", fmt.Errorf("no token in login response")
	}
	return token, nil
}

func setupProduct(client *http.Client, baseURL, token string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"product_name":        "Perf Product",
		"product_description": "Load test product",
		"product_price":       99.99,
		"product_image":       "perf.png",
		"product_category":    "test",
		"product_stock":       1000,
		"product_rating":      4.0,
		"product_reviews":     []string{"ok"},
	})
	resp, err := doRequest(client, baseURL, http.MethodPost, "/admin/addproduct", body, token)
	if err != nil {
		return "", err
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	resp, err = doRequest(client, baseURL, http.MethodGet, "/users/productview", nil, "")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var products []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return "", err
	}
	for _, p := range products {
		if name, _ := p["product_name"].(string); name == "Perf Product" {
			if id, _ := p["product_id"].(string); id != "" {
				return id, nil
			}
		}
	}
	return "", fmt.Errorf("perf product not found")
}

func doRequest(client *http.Client, baseURL, method, path string, body []byte, token string) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, baseURL+path, reader)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return client.Do(req)
}

func runScenario(client *http.Client, baseURL, token string, sc scenario, total, workers int) {
	results := make([]result, total)
	var idx atomic.Int64

	var wg sync.WaitGroup
	start := time.Now()

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(idx.Add(1)) - 1
				if i >= total {
					return
				}

				t0 := time.Now()
				resp, err := doRequest(client, baseURL, sc.method, sc.path, sc.body, boolStr(sc.auth, token))
				elapsed := time.Since(t0).Seconds() * 1000

				r := result{ms: elapsed, err: err != nil}
				if err == nil {
					r.status = resp.StatusCode
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				}
				results[i] = r
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	printReport(sc.name, results, duration)
}

func boolStr(auth bool, token string) string {
	if auth {
		return token
	}
	return ""
}

func printReport(name string, results []result, duration time.Duration) {
	var latencies []float64
	success := 0
	errors := 0
	statusCounts := map[int]int{}

	for _, r := range results {
		if r.err {
			errors++
			continue
		}
		statusCounts[r.status]++
		if r.status >= 200 && r.status < 300 {
			success++
			latencies = append(latencies, r.ms)
		}
	}

	sort.Float64s(latencies)

	fmt.Printf("Scenario: %s\n", name)
	fmt.Printf("  Duration:     %s\n", duration.Round(time.Millisecond))
	fmt.Printf("  Total:        %d\n", len(results))
	fmt.Printf("  Success:      %d (%.1f%%)\n", success, pct(success, len(results)))
	fmt.Printf("  Errors:       %d\n", errors)
	fmt.Printf("  Throughput:   %.1f req/s\n", float64(len(results))/duration.Seconds())

	if len(statusCounts) > 0 {
		fmt.Print("  Status codes: ")
		first := true
		for code, count := range statusCounts {
			if !first {
				fmt.Print(", ")
			}
			fmt.Printf("%d=%d", code, count)
			first = false
		}
		fmt.Println()
	}

	if len(latencies) == 0 {
		fmt.Println("  Latency:      no successful responses")
		return
	}

	fmt.Printf("  Latency ms:   min=%.1f  avg=%.1f  max=%.1f\n",
		latencies[0], avg(latencies), latencies[len(latencies)-1])
	fmt.Printf("  Percentiles:  p50=%.1f  p95=%.1f  p99=%.1f\n",
		percentile(latencies, 50), percentile(latencies, 95), percentile(latencies, 99))
}

func pct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) / float64(total) * 100
}

func avg(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	sum := 0.0
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}

func percentile(sorted []float64, p int) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(p)/100*float64(len(sorted)-1) + 0.5)
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
