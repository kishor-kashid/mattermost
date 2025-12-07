// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	defaultFetchTimeout        = 10 * time.Second
	maxFetchBytes        int64 = 5 * 1024 * 1024 // 5MB
	maxRedirects               = 5
	maxConcurrentFetches       = 8
)

var (
	errTooLarge         = errors.New("content exceeds size limit")
	errTooManyRedirects = errors.New("too many redirects")
	fetchLimiterOnce    sync.Once
	fetchLimiter        chan struct{}
)

func initFetchLimiter() {
	fetchLimiterOnce.Do(func() {
		fetchLimiter = make(chan struct{}, maxConcurrentFetches)
	})
}

// fetchURL retrieves the raw body for a URL with sane defaults.
func (a *App) fetchURL(ctx context.Context, url string) ([]byte, string, error) {
	initFetchLimiter()
	select {
	case fetchLimiter <- struct{}{}:
		defer func() { <-fetchLimiter }()
	case <-ctx.Done():
		return nil, "", ctx.Err()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	// Use a realistic browser User-Agent to avoid bot detection
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// Note: Don't set Accept-Encoding manually - Go's default transport handles gzip automatically
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{
		Timeout: defaultFetchTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return errTooManyRedirects
			}
			return nil
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", errors.New("failed to fetch url: " + resp.Status)
	}

	limited := io.LimitReader(resp.Body, maxFetchBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", err
	}

	if int64(len(body)) > maxFetchBytes {
		return nil, "", errTooLarge
	}

	return body, resp.Header.Get("Content-Type"), nil
}
