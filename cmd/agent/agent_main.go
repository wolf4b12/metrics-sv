package main

import (
    "math/rand"
    "time"
    "github.com/wolf4b12/metrics-sv.git/internal/agent/agentmethods" // Импортируем пакет agentmethods
    "github.com/wolf4b12/metrics-sv.git/internal/agent/parseflags" // Импортируем пакет parseflags
)


func main() {

    rand.New(rand.NewSource(time.Now().UnixNano())) // Create new source for random numbers

    poll, report, addr := parseflags.ParseFlags()
    
    agent2 := agentmethods.NewAgent(poll, report, addr)


    go agent2.CollectMetrics()
    go agent2.SendJSONCollectedMetrics()
    go agent2.SendTextCollectedMetrics()
    




    select {} // Keep main goroutine alive  
}