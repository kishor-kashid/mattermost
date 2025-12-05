// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"fmt"
	"time"
)

// ParseTimeRange converts a time range string to start and end timestamps
// Supports formats like "24h", "7d", "1w", "1m"
func ParseTimeRange(timeRange string) (int64, int64, error) {
	now := time.Now()
	var duration time.Duration

	switch timeRange {
	case "1h":
		duration = 1 * time.Hour
	case "6h":
		duration = 6 * time.Hour
	case "12h":
		duration = 12 * time.Hour
	case "24h", "1d":
		duration = 24 * time.Hour
	case "7d", "1w":
		duration = 7 * 24 * time.Hour
	case "30d", "1m":
		duration = 30 * 24 * time.Hour
	case "90d", "3m":
		duration = 90 * 24 * time.Hour
	default:
		return 0, 0, fmt.Errorf("unsupported time range: %s", timeRange)
	}

	startTime := now.Add(-duration)
	endTime := now

	return GetMillis(startTime), GetMillis(endTime), nil
}

// GetTimeRangeDescription returns a human-readable description of a time range
func GetTimeRangeDescription(startMs, endMs int64) string {
	start := time.Unix(0, startMs*int64(time.Millisecond))
	end := time.Unix(0, endMs*int64(time.Millisecond))

	duration := end.Sub(start)

	switch {
	case duration < time.Hour:
		return fmt.Sprintf("last %d minutes", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("last %d hours", int(duration.Hours()))
	case duration < 7*24*time.Hour:
		return fmt.Sprintf("last %d days", int(duration.Hours()/24))
	case duration < 30*24*time.Hour:
		return fmt.Sprintf("last %d weeks", int(duration.Hours()/(24*7)))
	default:
		return fmt.Sprintf("last %d days", int(duration.Hours()/24))
	}
}

// IsWithinTimeRange checks if a timestamp is within a given time range
func IsWithinTimeRange(timestamp, startMs, endMs int64) bool {
	return timestamp >= startMs && timestamp <= endMs
}

// GetTimestampBuckets returns time bucket labels for analytics
// Returns daily buckets for the given time range
func GetTimestampBuckets(startMs, endMs int64, bucketSize time.Duration) []time.Time {
	start := time.Unix(0, startMs*int64(time.Millisecond))
	end := time.Unix(0, endMs*int64(time.Millisecond))

	buckets := []time.Time{}
	current := start

	for current.Before(end) || current.Equal(end) {
		buckets = append(buckets, current)
		current = current.Add(bucketSize)
	}

	return buckets
}

// GetDailyBuckets returns daily time buckets for a date range
func GetDailyBuckets(startMs, endMs int64) []time.Time {
	return GetTimestampBuckets(startMs, endMs, 24*time.Hour)
}

// GetHourlyBuckets returns hourly time buckets for a date range
func GetHourlyBuckets(startMs, endMs int64) []time.Time {
	return GetTimestampBuckets(startMs, endMs, time.Hour)
}

// FormatTimestamp formats a timestamp in a human-readable way
func FormatTimestamp(timestampMs int64) string {
	t := time.Unix(0, timestampMs*int64(time.Millisecond))
	return t.Format("Jan 2, 2006 3:04 PM")
}

// FormatDuration formats a duration in milliseconds to a human-readable string
func FormatDuration(durationMs int64) string {
	duration := time.Duration(durationMs) * time.Millisecond

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%d seconds", int(duration.Seconds()))
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes", int(duration.Minutes()))
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		if minutes > 0 {
			return fmt.Sprintf("%dh %dm", hours, minutes)
		}
		return fmt.Sprintf("%d hours", hours)
	default:
		days := int(duration.Hours() / 24)
		hours := int(duration.Hours()) % 24
		if hours > 0 {
			return fmt.Sprintf("%dd %dh", days, hours)
		}
		return fmt.Sprintf("%d days", days)
	}
}

// GetRelativeTimeString returns a relative time string like "2 hours ago"
func GetRelativeTimeString(timestampMs int64) string {
	t := time.Unix(0, timestampMs*int64(time.Millisecond))
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "just now"
	case duration < time.Hour:
		mins := int(duration.Minutes())
		if mins == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", mins)
	case duration < 24*time.Hour:
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case duration < 7*24*time.Hour:
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2, 2006")
	}
}

// GetMillis converts time.Time to Mattermost timestamp (milliseconds)
func GetMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

