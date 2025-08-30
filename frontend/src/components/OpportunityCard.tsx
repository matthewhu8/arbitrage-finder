import { useEffect, useState } from 'react';
import { ArbitrageOpportunity } from '../types';

interface Props {
  opportunity: ArbitrageOpportunity;
  onClick: () => void;
}

export function OpportunityCard({ opportunity, onClick }: Props) {
  const [timeLeft, setTimeLeft] = useState('');
  
  useEffect(() => {
    const updateTimer = () => {
      const expiresAt = new Date(opportunity.expires_at);
      const now = new Date();
      const diff = expiresAt.getTime() - now.getTime();
      
      if (diff <= 0) {
        setTimeLeft('Expired');
      } else {
        const minutes = Math.floor(diff / 60000);
        const seconds = Math.floor((diff % 60000) / 1000);
        setTimeLeft(`${minutes}:${seconds.toString().padStart(2, '0')}`);
      }
    };

    updateTimer();
    const interval = setInterval(updateTimer, 1000);
    
    return () => clearInterval(interval);
  }, [opportunity.expires_at]);

  const isHighProfit = opportunity.profit_percent >= 2.0;
  const isExpiring = timeLeft.includes(':') && parseInt(timeLeft.split(':')[0]) < 1;

  return (
    <div
      onClick={onClick}
      className={`
        bg-gray-800/80 backdrop-blur rounded-xl p-6 border cursor-pointer
        transition-all duration-300 hover:scale-105 hover:shadow-2xl
        animate-slide-in
        ${isHighProfit ? 'border-yellow-500 shadow-yellow-500/20' : 'border-gray-700'}
        ${isExpiring ? 'animate-pulse-slow' : ''}
      `}
    >
      {/* Header */}
      <div className="flex justify-between items-start mb-4">
        <div className="flex items-center space-x-2">
          <span className="text-xs bg-blue-600 text-white px-2 py-1 rounded">
            {opportunity.sport}
          </span>
          {isHighProfit && (
            <span className="text-xs bg-yellow-500 text-black px-2 py-1 rounded font-bold">
              HOT 🔥
            </span>
          )}
        </div>
        <div className={`text-xs font-mono ${
          timeLeft === 'Expired' ? 'text-red-500' :
          isExpiring ? 'text-yellow-500' : 'text-gray-400'
        }`}>
          {timeLeft}
        </div>
      </div>

      {/* Teams */}
      <div className="mb-4">
        <div className="text-white font-semibold text-lg">
          {opportunity.home_team} vs {opportunity.away_team}
        </div>
      </div>

      {/* Profit */}
      <div className="mb-4">
        <div className={`text-3xl font-bold ${
          isHighProfit ? 'text-yellow-400' : 'text-green-400'
        }`}>
          {opportunity.profit_percent.toFixed(2)}%
        </div>
        <div className="text-xs text-gray-500">Guaranteed Profit</div>
      </div>

      {/* Betting Details */}
      <div className="space-y-2 text-sm">
        <div className="flex justify-between text-gray-300">
          <span>{opportunity.home_team} @ {opportunity.bookmaker_home}</span>
          <span className="font-mono">{opportunity.home_odds.toFixed(2)}</span>
        </div>
        <div className="flex justify-between text-gray-300">
          <span>{opportunity.away_team} @ {opportunity.bookmaker_away}</span>
          <span className="font-mono">{opportunity.away_odds.toFixed(2)}</span>
        </div>
      </div>

      {/* Stakes */}
      <div className="mt-4 pt-4 border-t border-gray-700">
        <div className="grid grid-cols-2 gap-2 text-xs">
          <div>
            <div className="text-gray-500">Stake {opportunity.home_team}</div>
            <div className="text-white font-mono">${opportunity.home_stake.toFixed(2)}</div>
          </div>
          <div>
            <div className="text-gray-500">Stake {opportunity.away_team}</div>
            <div className="text-white font-mono">${opportunity.away_stake.toFixed(2)}</div>
          </div>
        </div>
        <div className="mt-2 text-center">
          <div className="text-gray-500 text-xs">Expected Return</div>
          <div className="text-green-400 font-mono font-semibold">
            ${opportunity.expected_return.toFixed(2)}
          </div>
        </div>
      </div>
    </div>
  );
}