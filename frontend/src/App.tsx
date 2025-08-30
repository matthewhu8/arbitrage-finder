import React, { useEffect, useState } from 'react';
import { ArbitrageDashboard } from './components/ArbitrageDashboard';
import { useWebSocket } from './hooks/useWebSocket';
import { ArbitrageOpportunity } from './types';
import { Toaster } from 'react-hot-toast';
import './index.css';

function App() {
  const [opportunities, setOpportunities] = useState<ArbitrageOpportunity[]>([]);
  const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'disconnected'>('connecting');

  // Connect to WebSocket for real-time updates
  const { lastMessage, readyState } = useWebSocket('ws://localhost:8080/ws', {
    onOpen: () => setConnectionStatus('connected'),
    onClose: () => setConnectionStatus('disconnected'),
    onError: () => setConnectionStatus('disconnected'),
    reconnectInterval: 3000,
  });

  // Handle incoming WebSocket messages
  useEffect(() => {
    if (lastMessage) {
      try {
        const data = JSON.parse(lastMessage);
        if (data.type === 'arbitrage') {
          const newOpportunity = data.data as ArbitrageOpportunity;
          
          setOpportunities(prev => {
            // Remove duplicates and add new opportunity
            const filtered = prev.filter(opp => opp.id !== newOpportunity.id);
            return [newOpportunity, ...filtered].slice(0, 50); // Keep last 50
          });
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    }
  }, [lastMessage]);

  // Load initial opportunities
  useEffect(() => {
    fetch('/api/arbitrage')
      .then(res => res.json())
      .then(data => setOpportunities(data))
      .catch(err => console.error('Error fetching arbitrage:', err));
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-blue-900 to-gray-900">
      <Toaster position="top-right" />
      
      <header className="bg-black/50 backdrop-blur-sm border-b border-gray-700">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <h1 className="text-3xl font-bold text-white">
                ⚡ Sports Arbitrage Finder
              </h1>
              <span className="text-xs text-gray-400">Real-Time</span>
            </div>
            
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <div className={`w-3 h-3 rounded-full ${
                  connectionStatus === 'connected' ? 'bg-green-500 animate-pulse' :
                  connectionStatus === 'connecting' ? 'bg-yellow-500 animate-pulse' :
                  'bg-red-500'
                }`} />
                <span className="text-sm text-gray-300">
                  {connectionStatus === 'connected' ? 'Live' :
                   connectionStatus === 'connecting' ? 'Connecting...' :
                   'Disconnected'}
                </span>
              </div>
              
              <div className="text-sm text-gray-400">
                {opportunities.length} Active Opportunities
              </div>
            </div>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <ArbitrageDashboard opportunities={opportunities} />
      </main>

      <footer className="mt-16 py-8 text-center text-gray-500 text-sm">
        <p>Arbitrage opportunities update in real-time via WebSocket</p>
        <p className="mt-2">Data streams from Kafka → Detector → WebSocket → Dashboard</p>
      </footer>
    </div>
  );
}

export default App;