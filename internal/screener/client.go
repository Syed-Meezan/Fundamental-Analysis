// Package screener fetches and parses public company pages from
// screener.in for use in fundamental analysis.
package screener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	"sampleproject/internal/model"
)

const (
	baseURL   = "https://www.screener.in"
	userAgent = "Mozilla/5.0 (compatible; FundamentalAnalysisTool/1.0; personal use)"
	cacheTTL  = 15 * time.Minute
	// minRequestGap throttles outbound requests so this tool stays a light,
	// respectful client of screener.in rather than hammering it.
	minRequestGap = 400 * time.Millisecond
)

// companyURLPattern restricts which relative paths we will ever fetch, so a
// caller-supplied URL can never be used to make this server request an
// arbitrary host (SSRF) - only screener.in company pages are allowed.
var companyURLPattern = regexp.MustCompile(`^/company/(id/[0-9]+|[A-Za-z0-9.\-]+)/(consolidated/)?$`)

// Client talks to screener.in over HTTP with a small in-memory cache and
// basic rate limiting.
type Client struct {
	http *http.Client

	mu      sync.Mutex
	cache   map[string]cacheEntry
	lastReq time.Time
}

type cacheEntry struct {
	body    []byte
	expires time.Time
}

// NewClient builds a ready-to-use screener.in client.
func NewClient() *Client {
	return &Client{
		http:  &http.Client{Timeout: 15 * time.Second},
		cache: make(map[string]cacheEntry),
	}
}

// Search looks up companies by (partial) name via screener's own search API.
func (c *Client) Search(query string) ([]model.SearchResult, error) {
	u := baseURL + "/api/company/search/?q=" + url.QueryEscape(query)
	body, err := c.get(u)
	if err != nil {
		return nil, err
	}
	var results []model.SearchResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("decoding search response: %w", err)
	}
	return results, nil
}

// FetchCompanyPage fetches the raw HTML for a company page. relURL must be a
// path previously returned by Search (e.g. "/company/TCS/consolidated/").
func (c *Client) FetchCompanyPage(relURL string) ([]byte, error) {
	if !companyURLPattern.MatchString(relURL) {
		return nil, fmt.Errorf("invalid company url: %q", relURL)
	}
	return c.get(baseURL + relURL)
}

// get performs a rate-limited, cached GET request.
func (c *Client) get(fullURL string) ([]byte, error) {
	c.mu.Lock()
	if entry, ok := c.cache[fullURL]; ok && time.Now().Before(entry.expires) {
		c.mu.Unlock()
		return entry.body, nil
	}
	wait := minRequestGap - time.Since(c.lastReq)
	c.mu.Unlock()

	if wait > 0 {
		time.Sleep(wait)
	}

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html,application/json")
	req.Header.Set("Accept-Language", "en-IN,en;q=0.9")

	c.mu.Lock()
	c.lastReq = time.Now()
	c.mu.Unlock()

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20)) // 10MB cap
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("screener.in returned %s for %s", resp.Status, fullURL)
	}

	c.mu.Lock()
	c.cache[fullURL] = cacheEntry{body: body, expires: time.Now().Add(cacheTTL)}
	c.mu.Unlock()

	return body, nil
}
