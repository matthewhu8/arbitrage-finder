-- Create tables for storing arbitrage history and analytics

CREATE TABLE IF NOT EXISTS arbitrage_history (
    id UUID PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL,
    sport VARCHAR(50) NOT NULL,
    home_team VARCHAR(255) NOT NULL,
    away_team VARCHAR(255) NOT NULL,
    bookmaker_home VARCHAR(100) NOT NULL,
    bookmaker_away VARCHAR(100) NOT NULL,
    home_odds DECIMAL(10, 3) NOT NULL,
    away_odds DECIMAL(10, 3) NOT NULL,
    profit_percent DECIMAL(10, 3) NOT NULL,
    home_stake DECIMAL(10, 2) NOT NULL,
    away_stake DECIMAL(10, 2) NOT NULL,
    total_stake DECIMAL(10, 2) NOT NULL,
    expected_return DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    status VARCHAR(20) DEFAULT 'active'
);

CREATE INDEX idx_arbitrage_created_at ON arbitrage_history(created_at DESC);
CREATE INDEX idx_arbitrage_event_id ON arbitrage_history(event_id);
CREATE INDEX idx_arbitrage_profit ON arbitrage_history(profit_percent DESC);

CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(255) PRIMARY KEY,
    sport VARCHAR(50) NOT NULL,
    league VARCHAR(100),
    home_team VARCHAR(255) NOT NULL,
    away_team VARCHAR(255) NOT NULL,
    start_time TIMESTAMP,
    status VARCHAR(20) DEFAULT 'upcoming',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS odds_history (
    id SERIAL PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL,
    bookmaker VARCHAR(100) NOT NULL,
    home_odds DECIMAL(10, 3) NOT NULL,
    away_odds DECIMAL(10, 3) NOT NULL,
    draw_odds DECIMAL(10, 3),
    market_type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

CREATE INDEX idx_odds_event_timestamp ON odds_history(event_id, timestamp DESC);
CREATE INDEX idx_odds_bookmaker ON odds_history(bookmaker);