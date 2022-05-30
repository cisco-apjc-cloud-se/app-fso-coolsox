package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/microservices-demo/payment"
	stdopentracing "github.com/opentracing/opentracing-go"
	// zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"golang.org/x/net/context"

	// AppDynamics Go SDK Agent
	"strconv"
	appd "appdynamics"
)

const (
	ServiceName = "payment"
)

func main() {
	var (
		port          = flag.String("port", "8080", "Port to bind HTTP listener")
		// zip           = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
		declineAmount = flag.Float64("decline", 100, "Decline payments over certain amount")
	)
	flag.Parse()

	// AppDynamics SDK initialization
	cfg := appd.Config{}
	cfg.AppName = os.Getenv("APPD_APPNAME") // "exampleapp"
	cfg.TierName = os.Getenv("APPD_TIERNAME")
	var (
		appd_hostname string
		appd_usessl bool
		appd_port uint64
		// err error
	)
	appd_hostname, err = os.Hostname()
	if err != nil {
		fmt.Printf("Error determining hostname\n")
	}
	cfg.NodeName = appd_hostname // os.Getenv("APPD_NODENAME")
	cfg.Controller.Host = os.Getenv("APPD_CONTROLLER_HOST") // "my-appd-controller.example.org"
	appd_port, err = strconv.ParseUint(os.Getenv("APPD_CONTROLLER_PORT"), 10 ,16)
	if err != nil {
		fmt.Printf("Error converting AppDynamics Controller Port from environmental variable\n")
	}
	cfg.Controller.Port = uint16(appd_port) // 8090
	appd_usessl, err = strconv.ParseBool(os.Getenv("APPD_CONTROLLER_USE_SSL"))
	if err != nil {
		fmt.Printf("Error converting AppDynamics Controller SSL use from environmental variable\n")
	}
	cfg.Controller.UseSSL = appd_usessl // false
	cfg.Controller.Account = os.Getenv("APPD_CONTROLLER_ACCOUNT") // "customer1"
	cfg.Controller.AccessKey = os.Getenv("APPD_CONTROLLER_ACCESS_KEY")// "secret"
	cfg.InitTimeoutMs = 5000  // Wait up to 1s for initialization to finish // Needs to be >1s <5s for Backend..

	if err := appd.InitSDK(&cfg); err != nil {
		fmt.Printf("Error initializing the AppDynamics SDK\n")
	} else {
		fmt.Printf("Initialized AppDynamics SDK successfully\n")
	}

	var tracer stdopentracing.Tracer
	{
		// Log domain.
		var logger log.Logger
		{
			logger = log.NewLogfmtLogger(os.Stderr)
			logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
			logger = log.NewContext(logger).With("caller", log.DefaultCaller)
		}
		// if *zip == "" {
			tracer = stdopentracing.NoopTracer{}
		// } else {
		// 	logger := log.NewContext(logger).With("tracer", "Zipkin")
		// 	logger.Log("addr", zip)
		// 	collector, err := zipkin.NewHTTPCollector(
		// 		*zip,
		// 		zipkin.HTTPLogger(logger),
		// 	)
		// 	if err != nil {
		// 		logger.Log("err", err)
		// 		os.Exit(1)
		// 	}
		// 	tracer, err = zipkin.NewTracer(
		// 		zipkin.NewRecorder(collector, false, fmt.Sprintf("localhost:%v", port), ServiceName),
		// 	)
		// 	if err != nil {
		// 		logger.Log("err", err)
		// 		os.Exit(1)
		// 	}
		// }
		stdopentracing.InitGlobalTracer(tracer)

	}
	// Mechanical stuff.
	errc := make(chan error)
	ctx := context.Background()

	handler, logger := payment.WireUp(ctx, float32(*declineAmount), tracer, ServiceName)

	// Create and launch the HTTP server.
	go func() {
		logger.Log("transport", "HTTP", "port", *port)
		errc <- http.ListenAndServe(":"+*port, handler)
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
