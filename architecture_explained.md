# Sports Arbitrage Architecture - Simple Explanation

## What We're Building
A system that finds risk-free betting opportunities by comparing odds across different sportsbooks in real-time.

## How Arbitrage Works
```
DraftKings:  Lakers 2.10 (+110)  | Celtics 1.95 (-105)
FanDuel:     Lakers 1.92 (-108)  | Celtics 2.15 (+115)

If you bet:
- $48.78 on Lakers at DraftKings (2.10)
- $51.22 on Celtics at FanDuel (2.15)

Total bet: $100
Guaranteed return: $102.44
Profit: $2.44 (2.44%)
```

## System Components - Start Simple

### Version 1: Basic (Start Here)
```
[Sportsbook APIs] → [Your Go Program] → [Terminal Output]
```
- Fetch odds from APIs
- Calculate if arbitrage exists
- Print results

### Version 2: Add Storage
```
[Sportsbook APIs] → [Go Fetcher] → [Database] → [Calculator] → [Terminal]
```
- Store odds in database
- Separate fetching from calculating
- Track historical opportunities

### Version 3: Real-Time (Week 2)
```
[Sportsbook APIs] → [Go Fetcher] → [Message Queue] → [Calculator] → [WebSocket] → [Browser]
```
- Push updates immediately
- Multiple users can connect
- Real-time notifications

### Version 4: Production (Final)
```
                    ┌─────────────┐
                    │ Load        │
                    │ Balancer    │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌─────▼────┐      ┌─────▼────┐
   │Fetcher 1│       │Fetcher 2 │      │Fetcher 3 │
   └────┬────┘       └─────┬────┘      └─────┬────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                    ┌──────▼──────┐
                    │   Message   │
                    │    Queue     │
                    │   (Kafka)    │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼──────┐    ┌──────▼──────┐   ┌──────▼──────┐
   │Calculator │    │ Calculator  │   │ Calculator  │
   │     1     │    │      2      │   │      3      │
   └────┬──────┘    └──────┬──────┘   └──────┬──────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                    ┌──────▼──────┐
                    │  WebSocket  │
                    │   Server    │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │   Browser   │
                    │     App     │
                    └─────────────┘
```

## Data Flow - Step by Step

### 1. Fetching Odds (Every 5-30 seconds)
```go
// Concurrent fetching from multiple sportsbooks
go fetchDraftKings()  // Goroutine 1
go fetchFanDuel()     // Goroutine 2  
go fetchBetMGM()      // Goroutine 3
// All run at the same time!
```

### 2. Message Queue (Why We Need It)
```
Without Queue:
Fetcher → Calculator (tightly coupled, if calculator crashes, we lose data)

With Queue:
Fetcher → Queue → Calculator (decoupled, queue stores messages)
```

Benefits:
- Fetchers don't wait for calculators
- Can replay messages if something crashes
- Scale fetchers and calculators independently

### 3. Arbitrage Detection
```
For each game:
1. Collect all odds for both teams
2. Find best odds for each outcome
3. Calculate: 1/odds1 + 1/odds2
4. If sum < 1.0, arbitrage exists!
5. Send alert to users
```

### 4. WebSocket (Real-Time Updates)
```
Traditional HTTP: Browser asks server every second "Any updates?"
WebSocket: Server pushes to browser instantly when arbitrage found
```

## Why Each Technology?

### Go
- **Goroutines**: Fetch from 100 sportsbooks simultaneously
- **Channels**: Pass odds data between components safely
- **Speed**: Process millions of odds quickly

### Kafka/Pulsar (Message Queue)
- **Buffer**: Store odds if calculator is busy
- **Reliability**: Don't lose data if system crashes
- **Scale**: Add more calculators during busy times

### Redis/KeyDB (Cache)
- **Speed**: Access current odds in microseconds
- **TTL**: Auto-delete old odds
- **Pub/Sub**: Notify components of updates

### PostgreSQL/ScyllaDB (Database)
- **History**: Track all arbitrage opportunities
- **Analytics**: Which sportsbooks have best odds
- **Users**: Store preferences and alerts

### WebSocket
- **Real-time**: Push updates instantly
- **Efficiency**: One connection for all updates
- **Bidirectional**: Users can send filters

## Start Simple, Grow Complex

### Week 1: Build This
```go
func main() {
    odds1 := fetchDraftKings()
    odds2 := fetchFanDuel()
    
    if calculateArbitrage(odds1, odds2) {
        fmt.Println("ARBITRAGE FOUND!")
    }
}
```

### Week 2: Add Concurrency
```go
func main() {
    ch := make(chan Odds)
    
    go fetchDraftKings(ch)
    go fetchFanDuel(ch)
    go fetchBetMGM(ch)
    
    for odds := range ch {
        checkArbitrage(odds)
    }
}
```

### Week 3: Add Real-Time
```go
func main() {
    // WebSocket server
    http.HandleFunc("/ws", handleWebSocket)
    
    // Send updates to all connected clients
    for arb := range arbitrageChannel {
        broadcast(arb)
    }
}
```

## Common Questions

**Q: Why not just poll APIs constantly?**
A: Rate limits. Most sportsbooks limit to 10-100 requests/minute.

**Q: Why separate fetcher and calculator?**
A: Fetcher waits on network (slow). Calculator does math (fast). Don't mix.

**Q: Why message queue instead of direct connection?**
A: If calculator crashes, fetcher keeps working. Messages wait in queue.

**Q: How fast do odds change?**
A: Major changes every 30-60 seconds. Critical to be fast.

**Q: What's the hardest part?**
A: Getting reliable data feeds. Many sportsbooks block scrapers.

## Your Learning Path

1. **Understand the problem** (arbitrage math) ✓
2. **Build single-threaded version** (main.go)
3. **Add concurrency** (goroutines)
4. **Add persistence** (database)
5. **Add real-time** (websockets)
6. **Add message queue** (Kafka)
7. **Scale horizontally** (multiple instances)

Focus on steps 1-3 first. The rest is optimization.