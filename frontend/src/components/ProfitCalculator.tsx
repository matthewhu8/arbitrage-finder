import React, { useState } from 'react';
import { ArbitrageOpportunity } from '../types';

interface Props {
  opportunity: ArbitrageOpportunity;
  onClose: () => void;
}

export function ProfitCalculator({ opportunity, onClose }: Props) {
  const [bankroll, setBankroll] = useState(1000);
  
  // Calculate stakes based on custom bankroll
  const impliedProbHome = 1 / opportunity.home_odds;
  const impliedProbAway = 1 / opportunity.away_odds;
  const totalImpliedProb = impliedProbHome + impliedProbAway;
  
  const homeStake = (bankroll * impliedProbHome) / totalImpliedProb;
  const awayStake = (bankroll * impliedProbAway) / totalImpliedProb;
  const expectedReturn = bankroll * (1 + opportunity.profit_percent / 100);
  const profit = expectedReturn - bankroll;

  return (
    <div className="fixed inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-gray-900 rounded-2xl p-8 max-w-2xl w-full mx-4 border border-gray-700">
        <div className="flex justify-between items-start mb-6">
          <h2 className="text-2xl font-bold text-white">Profit Calculator</h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-white transition"
          >
            ✕
          </button>
        </div>

        {/* Game Info */}
        <div className="mb-6 p-4 bg-gray-800 rounded-lg">
          <div className="text-lg font-semibold text-white mb-2">
            {opportunity.home_team} vs {opportunity.away_team}
          </div>
          <div className="text-sm text-gray-400">
            Sport: {opportunity.sport} • Profit: {opportunity.profit_percent.toFixed(2)}%
          </div>
        </div>

        {/* Bankroll Input */}
        <div className="mb-6">
          <label className="block text-sm text-gray-400 mb-2">
            Total Bankroll to Use
          </label>
          <input
            type="number"
            value={bankroll}
            onChange={(e) => setBankroll(Number(e.target.value))}
            className="w-full px-4 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-blue-500"
            min="1"
            step="100"
          />
        </div>

        {/* Betting Instructions */}
        <div className="space-y-4 mb-6">
          <div className="p-4 bg-blue-900/30 border border-blue-700 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <div className="text-white font-semibold">
                Bet 1: {opportunity.home_team}
              </div>
              <div className="text-blue-400 font-mono text-lg">
                ${homeStake.toFixed(2)}
              </div>
            </div>
            <div className="text-sm text-gray-400">
              Place at {opportunity.bookmaker_home} • Odds: {opportunity.home_odds.toFixed(2)}
            </div>
          </div>

          <div className="p-4 bg-purple-900/30 border border-purple-700 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <div className="text-white font-semibold">
                Bet 2: {opportunity.away_team}
              </div>
              <div className="text-purple-400 font-mono text-lg">
                ${awayStake.toFixed(2)}
              </div>
            </div>
            <div className="text-sm text-gray-400">
              Place at {opportunity.bookmaker_away} • Odds: {opportunity.away_odds.toFixed(2)}
            </div>
          </div>
        </div>

        {/* Results */}
        <div className="p-4 bg-green-900/30 border border-green-700 rounded-lg">
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <div className="text-gray-400 text-sm">Total Stake</div>
              <div className="text-white font-mono text-xl">${bankroll.toFixed(2)}</div>
            </div>
            <div>
              <div className="text-gray-400 text-sm">Return</div>
              <div className="text-green-400 font-mono text-xl">${expectedReturn.toFixed(2)}</div>
            </div>
            <div>
              <div className="text-gray-400 text-sm">Net Profit</div>
              <div className="text-yellow-400 font-mono text-xl font-bold">
                +${profit.toFixed(2)}
              </div>
            </div>
          </div>
        </div>

        {/* Warning */}
        <div className="mt-6 p-3 bg-yellow-900/20 border border-yellow-800 rounded-lg">
          <div className="text-xs text-yellow-400">
            ⚠️ Place both bets quickly - odds change rapidly. This opportunity expires in real-time.
          </div>
        </div>
      </div>
    </div>
  );
}