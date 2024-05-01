package main

import (
	"final_project/initializers"
	"final_project/internal/router"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	initializers.GetKeysInEnv()
	initializers.ConnectDb()
}

func main() {
	var pingCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ping_request_count",
			Help: "No of request handled by Ping handler",
		},
	)
	prometheus.MustRegister(pingCounter)
	router := router.SetupRouter()
	router.Run(":8092")

}
