import React, { useState, useEffect } from 'react';
import { ArbitrageOpportunity } from '../types';
import { OpportunityCard } from './OpportunityCard';
import { ProfitCalculator } from './ProfitCalculator';
import toast from 'react-hot-toast';
import useSound from 'use-sound';

interface Props {
  opportunities: ArbitrageOpportunity[];
}

export function ArbitrageDashboard({ opportunities }: Props) {
  const [filter, setFilter] = useState<'all' | 'high-profit'>('all');
  const [selectedOpportunity, setSelectedOpportunity] = useState<ArbitrageOpportunity | null>(null);
  const [prevOpportunities, setPrevOpportunities] = useState<ArbitrageOpportunity[]>([]);
  
  // Sound notification (you'll need to add a sound file)
  const [playNotification] = useSound('/notification.mp3', { volume: 0.5 });

  // Detect new opportunities for notifications
  useEffect(() => {
    if (opportunities.length > 0 && prevOpportunities.length > 0) {
      const newOpps = opportunities.filter(
        opp => !prevOpportunities.find(prev => prev.id === opp.id)
      );

      newOpps.forEach(opp => {
        if (opp.profit_percent >= 2.0) {
          // High profit alert
          toast.success(
            `🔥 HIGH PROFIT: ${opp.profit_percent.toFixed(2)}% - ${opp.home_team} vs ${opp.away_team}`,
            { duration: 10000 }
          );
          playNotification();
        } else {
          // Normal alert
          toast(
            `New arbitrage: ${opp.profit_percent.toFixed(2)}% - ${opp.home_team} vs ${opp.away_team}`,
            { duration: 5000 }
          );
        }
      });
    }
    setPrevOpportunities(opportunities);
  }, [opportunities, prevOpportunities, playNotification]);

  const filteredOpportunities = opportunities.filter(opp => {
    if (filter === 'high-profit') return opp.profit_percent >= 2.0;
    return true;
  });

  const stats = {
    total: opportunities.length,
    avgProfit: opportunities.length > 0 
      ? opportunities.reduce((acc, opp) => acc + opp.profit_percent, 0) / opportunities.length
      : 0,
    highProfit: opportunities.filter(opp => opp.profit_percent >= 2.0).length,
  };

  return (
    <div className="space-y-6">
      {/* Stats Bar */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="bg-gray-800/50 backdrop-blur rounded-lg p-6 border border-gray-700">
          <div className="text-3xl font-bold text-white">{stats.total}</div>
          <div className="text-sm text-gray-400">Active Opportunities</div>
        </div>
        <div className="bg-gray-800/50 backdrop-blur rounded-lg p-6 border border-gray-700">
          <div className="text-3xl font-bold text-green-400">
            {stats.avgProfit.toFixed(2)}%
          </div>
          <div className="text-sm text-gray-400">Average Profit</div>
        </div>
        <div className="bg-gray-800/50 backdrop-blur rounded-lg p-6 border border-gray-700">
          <div className="text-3xl font-bold text-yellow-400">{stats.highProfit}</div>
          <div className="text-sm text-gray-400">High Profit (≥2%)</div>
        </div>
      </div>

      {/* Filters */}
      <div className="flex items-center space-x-4">
        <button
          onClick={() => setFilter('all')}
          className={`px-4 py-2 rounded-lg font-medium transition ${
            filter === 'all'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
          }`}
        >
          All Opportunities
        </button>
        <button
          onClick={() => setFilter('high-profit')}
          className={`px-4 py-2 rounded-lg font-medium transition ${
            filter === 'high-profit'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
          }`}
        >
          High Profit Only
        </button>
      </div>

      {/* Opportunities Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        {filteredOpportunities.length === 0 ? (
          <div className="col-span-full text-center py-12">
            <div className="text-gray-500">
              <div className="text-6xl mb-4">📊</div>
              <div className="text-xl">No arbitrage opportunities at the moment</div>
              <div className="text-sm mt-2">Opportunities will appear here in real-time</div>
            </div>
          </div>
        ) : (
          filteredOpportunities.map(opp => (
            <OpportunityCard
              key={opp.id}
              opportunity={opp}
              onClick={() => setSelectedOpportunity(opp)}
            />
          ))
        )}
      </div>

      {/* Profit Calculator Modal */}
      {selectedOpportunity && (
        <ProfitCalculator
          opportunity={selectedOpportunity}
          onClose={() => setSelectedOpportunity(null)}
        />
      )}
    </div>
  );
}