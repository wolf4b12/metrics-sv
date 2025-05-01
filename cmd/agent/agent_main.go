package main

import (
    "math/rand"
    "time"
    "github.com/wolf4b12/metrics-sv/internal/agent/agentmethods" // Импортируем пакет agentmethods
    "github.com/wolf4b12/metrics-sv/internal/agent/parseflags" // Импортируем пакет parseflags
)


func main() {

    rand.New(rand.NewSource(time.Now().UnixNano())) // Create new source for random numbers

    poll, report, addr := parseflags.ParseFlags()
    
    agent := agentmethods.NewAgent(poll, report, addr)

go   agent.StartCollectingMetrics()

go  agent.SendJSONCollectedMetrics()

go   agent.SendTextCollectedMetrics()
    




    select {} // Keep main goroutine alive  
}