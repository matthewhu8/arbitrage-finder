export interface ArbitrageOpportunity {
  id: string;
  event_id: string;
  sport: string;
  home_team: string;
  away_team: string;
  bookmaker_home: string;
  bookmaker_away: string;
  home_odds: number;
  away_odds: number;
  profit_percent: number;
  home_stake: number;
  away_stake: number;
  total_stake: number;
  expected_return: number;
  created_at: string;
  expires_at: string;
  status: 'active' | 'expired' | 'executed';
}

export interface OddsUpdate {
  id: string;
  event_id: string;
  sport: string;
  home_team: string;
  away_team: string;
  bookmaker: string;
  home_odds: number;
  away_odds: number;
  timestamp: string;
  market_type: string;
}

export interface WebSocketMessage {
  type: 'arbitrage' | 'odds_update' | 'status';
  data: any;
  timestamp: string;
}