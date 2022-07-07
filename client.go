package moonphase

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	BaseURL string
	APIKey  string

	mu    sync.Mutex
	cache map[string]MoonPhase
}

type MoonPhase struct {
	Name  string    `json:"text"`
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

func unknownMoonPhase() MoonPhase {
	return MoonPhase{
		Name: "Unknown",
	}
}

type AstronomyForecast struct {
	Moon struct {
		CurrentPhase MoonPhase `json:"current"`
	} `json:"moonPhase"`
	Time time.Time `json:"time"`
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		cache:   make(map[string]MoonPhase),
	}
}

func (m *Client) Current(lat, lng float64) MoonPhase {
	if cached, ok := m.cachedForecast(lat, lng); ok {
		return cached
	}

	q := url.Values{}
	q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("lng", strconv.FormatFloat(lng, 'f', -1, 64))

	resp, err := m.queryAPI(http.MethodGet, m.BaseURL+"/v2/astronomy/point?"+q.Encode(), nil)
	if err != nil {
		log.Fatalf("failed to query Stormglass API: %s", err)
	}
	defer resp.Body.Close()

	var data struct {
		Forecasts []AstronomyForecast `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatalf("failed to decode response data: %s", err)
	}

	if len(data.Forecasts) == 0 {
		return unknownMoonPhase()
	}

	sort.Slice(data.Forecasts, func(i, j int) bool {
		return data.Forecasts[i].Time.Before(data.Forecasts[j].Time)
	})

	return data.Forecasts[0].Moon.CurrentPhase
}

func (m *Client) cachedForecast(lat, lng float64) (MoonPhase, bool) {
	k := m.cacheKey(lat, lng)

	m.mu.Lock()
	defer m.mu.Unlock()

	cached, ok := m.cache[k]
	if !ok {
		return unknownMoonPhase(), false
	}

	if cached.Time.Sub(time.Now()) > 12*time.Hour {
		delete(m.cache, k)

		return unknownMoonPhase(), false
	}

	return cached, true
}

func (m *Client) queryAPI(method, u string, data io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, u, data)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", m.APIKey)

	return http.DefaultClient.Do(req)
}

func (*Client) cacheKey(lat, lng float64) string {
	return fmt.Sprintf("%f,%f", lat, lng)
}
