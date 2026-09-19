package checkorder

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func httpPostForm(targetURL string, data map[string]string, cookie string) ([]byte, error) {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func httpPostJSON(targetURL string, data map[string]string, headers map[string]string) ([]byte, error) {
	body, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func parseJSON(data []byte) map[string]interface{} {
	var result map[string]interface{}
	_ = json.Unmarshal(data, &result)
	return result
}

func ensureScheme(rawURL, def string) string {
	if len(rawURL) >= 7 && (rawURL[:7] == "http://" || rawURL[:8] == "https://") {
		return rawURL
	}
	return def + "://" + rawURL
}

func httpPostFormHeaders(targetURL string, data map[string]string, cookie string, headers map[string]string) ([]byte, error) {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func httpGet(targetURL string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func newCookieJar() http.CookieJar {
	jar, _ := cookiejar.New(nil)
	return jar
}

func readAll(r io.Reader) []byte {
	data, _ := io.ReadAll(r)
	return data
}
