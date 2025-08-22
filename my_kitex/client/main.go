package main

import (
	"context"
	"my_kitex/kitex_gen/math_service"
	"my_kitex/kitex_gen/math_service/math"
	"os"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/warmup"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	fout, err := os.OpenFile("log/client.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	klog.SetOutput(fout)
	klog.SetLevel(klog.LevelInfo)

	resolver, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		klog.Fatalf("connect to etcd failed: %s", err)
	}

	failurePolicy := retry.NewFailurePolicy()
	failurePolicy.WithMaxRetryTimes(2) //最多重试2次（不包含首次）,不是立即重试，中间要休息一段时间

	serviceCBSuit := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string { return ri.To().ServiceName() })                         //服务粒度的熔断控制
	methodCBSuit := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string { return ri.To().ServiceName() + "/" + ri.To().Method() }) //方法粒度的熔断控制

	client, err := math.NewClient(
		"dqq.math",
		//client.WithHostPorts("127.0.0.1:5678"),
		client.WithResolver(resolver),
		client.WithMiddleware(TimerMW),
		// client.WithMiddleware(AuthMW),
		client.WithConnectTimeout(100*time.Millisecond), //连接超时
		client.WithRPCTimeout(200*time.Millisecond),     //请求超时
		client.WithWarmingUp(&warmup.ClientOption{ //预先初始化服务发现和连接池的相关组件，避免在首次请求时产生较大的延迟
			ResolverOption: &warmup.ResolverOption{
				Dests: []*rpcinfo.EndpointBasicInfo{
					{
						ServiceName: "dqq.math",
						// Tags: map[string]string{
						// 	"cluster": "default",
						// 	// "env": "xxx"
						// },
					},
				},
			},
		}),
		client.WithFailureRetry(failurePolicy),                                          //重试机制
		client.WithLoadBalancer(loadbalance.NewWeightedRandomWithAliasMethodBalancer()), //负载均衡策略
		client.WithCircuitBreaker(serviceCBSuit),                                        //熔断机制，控制粒度为：下游服务
		client.WithCircuitBreaker(methodCBSuit),                                         //熔断机制，控制粒度为：下游服务+接口
		// client.WithFallback(fallback.NewFallbackPolicy(func(ctx context.Context, args utils.KitexArgs, result utils.KitexResult, err error) (fbErr error) {
		// 	result.SetSuccess(&math_service.AddResponse{Sum: -1}) //一般不使用client级别的降级策略，因为所有接口的返回数据类型不可能全是一样的
		// 	return
		// })), //RPC取不到正常结果时，采用降级方案
	)
	if err != nil {
		klog.Fatalf("create rpc client fail: %s", err)
	}

	//可以在任意位置调用UpdateServiceCBConfig更新熔断策略，更新后对后续的调用立即生效。注意这里的“立即”是近似的，因为理想情况下每一次调用都要开启一个滚动的统计窗口，而且要维护所有的滚动窗口，才能满足随时可能发生变化的熔断统计策略。实际上是划分成了很多固定长度的时间窗口，在每个时间窗口内单独统计调用数和失败数，切换窗口的时候会存在抖动问题。
	//熔断器开启后所有请求都直接被拒绝，这段时间被称为冷却期。
	//冷却期持续一段时间，然后熔断器进入半开启状态，即会放过一部分请求，当连续成功"若干数目"的请求后，熔断器将变为关闭，这是一个逐渐试探下游的过程。
	methodCBSuit.UpdateServiceCBConfig("Math/Add", circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   0.3,
		MinSample: 200,
	})
	methodCBSuit.UpdateServiceCBConfig("Math/Sub", circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   0.4,
		MinSample: 200,
	})
	serviceCBSuit.UpdateServiceCBConfig("Math", circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   0.3,
		MinSample: 500,
	})

	request := math_service.SubRequest{Left: 7, Right: 5}
	response, err := client.Sub(context.Background(),
		&request,
		//直连访问，不走服务发现
		//  callopt.WithURL("http://myserverdomain.com:8888"),
		//callopt.WithHostPort("127.0.0.1:5678"),           //很多client.Option和call.Option功能是重复的，只不过控制的粒度不一样
		callopt.WithConnectTimeout(100*time.Millisecond), //连接超时
		callopt.WithRPCTimeout(200*time.Millisecond),     //请求超时
		callopt.WithFallback( //接口级别的降级策略
			fallback.TimeoutAndCBFallback( //当发生超时或熔断error时
				fallback.UnwrapHelper(
					func(ctx context.Context, req, resp interface{}, err error) (fbResp interface{}, fbErr error) {
						fbResp = &math_service.SubResponse{Diff: -1} //统一返回-1
						fbErr = nil
						return
					},
				),
			),
		),
	)
	if err != nil {
		klog.Error(err)
	} else {
		klog.Info(response.Diff)
	}

}

// go run ./client
