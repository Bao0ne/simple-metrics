package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var AccessCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_requests_total",
	},
	[]string{"method", "path"},
	)
var QueueGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "queue_num_total",
	},
	[]string{"name"},
	)
var HttpDurationsHistogram = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:	"http_durations_histogram_seconds",
		Buckets: []float64{0.2, 0.5, 1, 2, 5, 10, 30},
	},
	[]string{"path"},
	)
var HttpDUrations = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name:	"http_durations_seconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
	[]string{"path"},
	)

func init()  {
	prometheus.MustRegister(AccessCounter, QueueGauge, HttpDurationsHistogram, HttpDUrations)
}

func main() {
	engine := gin.New()
	engine.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "豹一")
	})
	engine.GET("/counter", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		AccessCounter.With(prometheus.Labels{
			"method": c.Request.Method,
			"path": purl.Path,
		}).Add(1)
	})
	engine.GET("/queue", func(c *gin.Context) {
		num := c.Query("num")
		fnum, _ := strconv.ParseFloat(num, 32)
		QueueGauge.With(prometheus.Labels{"name": "queue_eddycjy"}).Set(fnum)
	})
	engine.GET("/histogram", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		HttpDurationsHistogram.With(prometheus.Labels{"path": purl.Path}).Observe(float64(rand.Intn(30)))
	})
	engine.GET("/summary", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		HttpDUrations.With(prometheus.Labels{"path": purl.Path}).Observe(float64(rand.Intn(30)))
	})
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	engine.Run(":10001")
}
