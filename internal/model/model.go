package model

import "time"

type CachedURL struct {
	ID  int    `json:"url_id"`
	URL string `json:"full_url"`
}
type Link struct {
	FullURL   string `json:"full_url"`
	ShortCode string `json:"short_code"`
}
type GetAllResponse struct {
	Links []Link `json:"links"`
}
type Visit struct {
	UserAgent string    `json:"user-agent"`
	VisitedAt time.Time `json:"visited_at"`
}

type RawAnalyticsResponse struct {
	FullURL   string  `json:"full_url"`
	ShortCode string  `json:"short_code"`
	Clicks    int     `json:"clicks"`
	Visits    []Visit `json:"visits"`
}

type TimeStat struct {
	Period string `json:"period"`
	Clicks int    `json:"clicks"`
}

type TimeAnalyticsResponse struct {
	FullURL   string     `json:"full_url"`
	ShortCode string     `json:"short_code"`
	Data      []TimeStat `json:"data"`
}

type UserAgentStat struct {
	UserAgent string `json:"user_agent"`
	Clicks    int    `json:"clicks"`
}

type UserAgentAnalyticsResponse struct {
	FullURL   string          `json:"full_url"`
	ShortCode string          `json:"short_code"`
	Data      []UserAgentStat `json:"data"`
}
