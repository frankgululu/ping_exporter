package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//编写metric分为一下几个步骤
// 需要弄清楚应用程序暴露出来的指标意味着什么

//1. 定义metrics，这里分为4种类型的metrics
//2. 处理Request/recodrd metrics 重难点
//3. 注册metrics到prometheus
//4. 启动服务，路由metrics给promhttp.Handler()处理

// 1. 定义metrics
var (
	RequestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "Total number of HTTP requests received",
	})

	RequestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.LinearBuckets(0.001, 0.005, 10),
	})
)

type pingResponse struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("handling new get request")
		startTime := time.Now()
		RequestCounter.Inc()

		//Generate random number between 1 and 10
		randomNumber := rand.Intn(10) + 1

		var response pingResponse
		response.Message = "hello"

		if randomNumber < 5 {
			//sleep 5 s
			time.Sleep(3 * time.Second)
		}

		elapsed := time.Since(startTime)
		RequestDuration.Observe(elapsed.Seconds())

		w.Header().Set("Context-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			fmt.Println("Error encoding response:", err)
			return
		}
		slog.Info("processed request", "duration", elapsed)
	})

	//3. Resister metrics handler，有几个metric，注册几个
	prometheus.MustRegister(RequestCounter, RequestDuration)

	//4. Start http server
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8880", nil)
}
