package arbitrage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/matthewhu/sportarbitrage/internal/models"
)

// Calculator handles arbitrage detection
type Calculator struct {
	minProfit float64 // Minimum profit percentage to consider
}

// NewCalculator creates a new arbitrage calculator
func NewCalculator(minProfit float64) *Calculator {
	return &Calculator{
		minProfit: minProfit,
	}
}

// DetectArbitrage checks if arbitrage opportunity exists
func (c *Calculator) DetectArbitrage(odds1, odds2 *models.OddsUpdate) *models.ArbitrageOpportunity {
	// Must be same event
	if odds1.EventID != odds2.EventID {
		return nil
	}

	// Find best odds combination
	var bestHomeOdds, bestAwayOdds float64
	var bestHomeBookmaker, bestAwayBookmaker string

	// Check all combinations
	combinations := []struct {
		homeOdds      float64
		awayOdds      float64
		homeBookmaker string
		awayBookmaker string
	}{
		{odds1.HomeOdds, odds2.AwayOdds, odds1.Bookmaker, odds2.Bookmaker},
		{odds2.HomeOdds, odds1.AwayOdds, odds2.Bookmaker, odds1.Bookmaker},
		{odds1.HomeOdds, odds1.AwayOdds, odds1.Bookmaker, odds1.Bookmaker},
		{odds2.HomeOdds, odds2.AwayOdds, odds2.Bookmaker, odds2.Bookmaker},
	}

	for _, combo := range combinations {
		if combo.homeOdds > bestHomeOdds {
			bestHomeOdds = combo.homeOdds
			bestHomeBookmaker = combo.homeBookmaker
		}
		if combo.awayOdds > bestAwayOdds {
			bestAwayOdds = combo.awayOdds
			bestAwayBookmaker = combo.awayBookmaker
		}
	}

	// Calculate implied probabilities
	impliedProbHome := 1.0 / bestHomeOdds
	impliedProbAway := 1.0 / bestAwayOdds
	totalImpliedProb := impliedProbHome + impliedProbAway

	// Check for arbitrage (total probability < 100%)
	if totalImpliedProb >= 1.0 {
		return nil
	}

	// Calculate profit
	profitPercent := (1.0/totalImpliedProb - 1.0) * 100

	// Check minimum profit threshold
	if profitPercent < c.minProfit {
		return nil
	}

	// Calculate optimal stakes for $1000 total bet
	totalStake := 1000.0
	homeStake := (totalStake * impliedProbHome) / totalImpliedProb
	awayStake := (totalStake * impliedProbAway) / totalImpliedProb
	expectedReturn := totalStake * (1.0 + profitPercent/100)

	return &models.ArbitrageOpportunity{
		ID:            uuid.New().String(),
		EventID:       odds1.EventID,
		Sport:         odds1.Sport,
		HomeTeam:      odds1.HomeTeam,
		AwayTeam:      odds1.AwayTeam,
		BookmakerHome: bestHomeBookmaker,
		BookmakerAway: bestAwayBookmaker,
		HomeOdds:      bestHomeOdds,
		AwayOdds:      bestAwayOdds,
		ProfitPercent: profitPercent,
		HomeStake:     homeStake,
		AwayStake:     awayStake,
		TotalStake:    totalStake,
		ExpectedReturn: expectedReturn,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(5 * time.Minute),
		Status:        "active",
	}
}

// CalculateThreeWay handles arbitrage for events with draw option
func (c *Calculator) CalculateThreeWay(homeOdds, drawOdds, awayOdds float64) (bool, float64) {
	impliedProbHome := 1.0 / homeOdds
	impliedProbDraw := 1.0 / drawOdds
	impliedProbAway := 1.0 / awayOdds
	
	total := impliedProbHome + impliedProbDraw + impliedProbAway
	
	if total < 1.0 {
		profit := (1.0/total - 1.0) * 100
		return true, profit
	}
	
	return false, 0
}

// FormatOdds converts decimal odds to American format
func FormatOdds(decimal float64) string {
	if decimal >= 2.0 {
		american := (decimal - 1) * 100
		return fmt.Sprintf("+%.0f", american)
	} else {
		american := -100 / (decimal - 1)
		return fmt.Sprintf("%.0f", american)
	}
}