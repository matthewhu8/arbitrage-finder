#!/bin/bash

echo "🚀 Starting Sports Arbitrage Finder..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

echo "📦 Building services..."
docker-compose build

echo "🔄 Starting services..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

echo "✅ Services are starting up!"
echo ""
echo "📊 Dashboard: http://localhost:3000"
echo "🔌 API: http://localhost:8080"
echo "📝 Logs: docker-compose logs -f"
echo ""
echo "Services running:"
echo "  - Kafka (real-time streaming)"
echo "  - 3 Fetchers (DraftKings, FanDuel, BetMGM)"
echo "  - Detector (arbitrage detection)"
echo "  - WebSocket API"
echo "  - React Dashboard"
echo ""
echo "⚡ Arbitrage opportunities will appear in real-time on the dashboard!"