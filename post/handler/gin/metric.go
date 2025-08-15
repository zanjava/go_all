package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const SERVICE = "post"

var (
	// Counter是一个积累量（单调增），跟历史值有关
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{Name: "request_counter"}, []string{"service", "interface"}) //此处指定了2个Label
	// Gauge是每个记录是独立的
	requestTimer = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "request_timer"}, []string{"service", "interface"})
)

func Metric(ctx *gin.Context) {
	begin := time.Now()
	ctx.Next()
	ifc := mappingUrl(ctx)                             //以请求的url path作为接口标识
	requestCounter.WithLabelValues(SERVICE, ifc).Inc() //WithLabelValues()的参数跟Label一一对应
	requestTimer.WithLabelValues(SERVICE, ifc).Set(float64(time.Since(begin).Milliseconds()))
}

var (
	// 检查一下所有的gin路由，把restful参数全写在这个map里
	restfulMapping = map[string]string{"id": ":id"}
)

// restful里的参数映射成泛化的url路径
func mappingUrl(ctx *gin.Context) string {
	// ctx.Request.RequestURI 包含url里的get参数，所以要用ctx.Request.URL.Path，但是restful情况下ctx.Request.URL.Path里也包含具体的参数值，需要做映射替换
	url := ctx.Request.URL.Path    //   /news/3   -->  /news/:id
	for _, p := range ctx.Params { //遍历restful参数
		if value, exists := restfulMapping[p.Key]; exists {
			url = strings.Replace(url, p.Value, value, 1) //把restful参数的值替换成一个泛化的字符串
		}
	}
	return url
}

// import "github.com/prometheus/client_golang/prometheus/promhttp"
// engine.GET("/metrics", func(ctx *gin.Context) {
// 		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
// 	})

// sum(rate(request_counter{ job="dqq_golang", service="post",interface!~".*metric.*"}[1m])) by (interface)
// avg(request_timer{ job="dqq_golang", service="post",interface!="/metrics"}) by (interface)
