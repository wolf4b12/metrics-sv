package main

import (
    "context"
    "math/rand"
    "time"
    "github.com/wolf4b12/metrics-sv/internal/agent/agentmethods" // Импортируем пакет agentmethods
    "github.com/wolf4b12/metrics-sv/internal/agent/parseflags" // Импортируем пакет parseflags
)

func main() {
    rand.New(rand.NewSource(time.Now().UnixNano())) // Create new source for random numbers

    poll, report, addr := parseflags.ParseFlags()

    agent := agentmethods.NewAgent(poll, report, addr)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    agent.StartCollectingMetrics(ctx)
    agent.SendJSONCollectedMetrics()
    agent.SendTextCollectedMetrics()
    agent.CollectAndSendBatches(ctx)

    select {} // Keep main goroutine alive
}
