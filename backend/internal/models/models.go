package models

import (
	"time"
)

// OddsUpdate represents odds from a sportsbook
type OddsUpdate struct {
	ID          string    `json:"id"`
	EventID     string    `json:"event_id"`
	Sport       string    `json:"sport"`
	HomeTeam    string    `json:"home_team"`
	AwayTeam    string    `json:"away_team"`
	Bookmaker   string    `json:"bookmaker"`
	HomeOdds    float64   `json:"home_odds"`
	AwayOdds    float64   `json:"away_odds"`
	DrawOdds    float64   `json:"draw_odds,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	MarketType  string    `json:"market_type"` // moneyline, spread, total
}

// ArbitrageOpportunity represents a profitable betting opportunity
type ArbitrageOpportunity struct {
	ID            string    `json:"id"`
	EventID       string    `json:"event_id"`
	Sport         string    `json:"sport"`
	HomeTeam      string    `json:"home_team"`
	AwayTeam      string    `json:"away_team"`
	BookmakerHome string    `json:"bookmaker_home"`
	BookmakerAway string    `json:"bookmaker_away"`
	HomeOdds      float64   `json:"home_odds"`
	AwayOdds      float64   `json:"away_odds"`
	ProfitPercent float64   `json:"profit_percent"`
	HomeStake     float64   `json:"home_stake"`
	AwayStake     float64   `json:"away_stake"`
	TotalStake    float64   `json:"total_stake"`
	ExpectedReturn float64  `json:"expected_return"`
	CreatedAt     time.Time `json:"created_at"`
	ExpiresAt     time.Time `json:"expires_at"`
	Status        string    `json:"status"` // active, expired, executed
}

// Event represents a sporting event
type Event struct {
	ID        string    `json:"id"`
	Sport     string    `json:"sport"`
	League    string    `json:"league"`
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	StartTime time.Time `json:"start_time"`
	Status    string    `json:"status"`
}

// WebSocketMessage for real-time updates
type WebSocketMessage struct {
	Type    string      `json:"type"` // arbitrage, odds_update, status
	Data    interface{} `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}