# Learning Plan - Sports Arbitrage Project

## Phase 1: Go Basics (This Week)
- [ ] Run the basic main.go
- [ ] Understand structs and JSON parsing
- [ ] Add goroutines to fetch from multiple endpoints
- [ ] Implement proper error handling
- [ ] Add context for timeouts

## Phase 2: Concurrent Programming (Next Week)
- [ ] Build worker pool pattern
- [ ] Use channels for communication
- [ ] Implement rate limiting
- [ ] Add circuit breaker pattern

## Phase 3: Real-Time Features
- [ ] WebSocket server with Gorilla
- [ ] Redis pub/sub integration
- [ ] Message serialization with MessagePack
- [ ] Connection management

## Phase 4: Production Features
- [ ] Add Prometheus metrics
- [ ] Structured logging
- [ ] Configuration management
- [ ] Docker containerization

## Resources to Study

### Go Fundamentals
1. **Tour of Go**: https://go.dev/tour/
2. **Effective Go**: https://go.dev/doc/effective_go
3. **Go by Example**: https://gobyexample.com/

### Concurrency
1. **Concurrency in Go** by Katherine Cox-Buday
2. Rob Pike's talks on concurrency

### Real-Time Systems
1. Gorilla WebSocket examples
2. Redis Pub/Sub patterns

### Practice Projects Before Main Project
1. Build a concurrent web scraper
2. Create a simple chat server
3. Make a rate limiter
4. Build a connection pool

## Daily Practice
- Morning: Read one chapter/article
- Code: 2-3 hours implementing
- Evening: Review and refactor

## Key Concepts to Master
- [ ] Goroutines vs Threads
- [ ] Channel patterns (fan-in, fan-out)
- [ ] Context package
- [ ] Interface design
- [ ] Error handling patterns
- [ ] Testing in Go
- [ ] Benchmarking

## Milestones
- Week 1: Working single-threaded arbitrage detector
- Week 2: Concurrent multi-bookmaker fetcher
- Week 3: Real-time WebSocket updates
- Week 4: Production-ready with monitoring