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

    go agent.StartCollectingMetrics(ctx)
    go agent.SendJSONCollectedMetrics()
    go agent.SendTextCollectedMetrics(ctx)
    go agent.CollectAndSendBatches(ctx)

    select {} // Keep main goroutine alive
}
