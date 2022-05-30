package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	// "strings"
	"syscall"

	corelog "log"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/microservices-demo/user/api"
	"github.com/microservices-demo/user/db"
	"github.com/microservices-demo/user/db/mongodb"
	stdopentracing "github.com/opentracing/opentracing-go"
	// zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	commonMiddleware "github.com/weaveworks/common/middleware"

	// AppDynamics Go SDK Agent
	"strconv"
	appd "appdynamics"
)

var (
	port string
	// zip  string
)

var (
	HTTPLatency = stdprometheus.NewHistogramVec(stdprometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: stdprometheus.DefBuckets,
	}, []string{"method", "path", "status_code", "isWS"})
)

const (
	ServiceName = "user"
)

//AppD middleware to enclose the routing handling functions and monitor Business Transactions
func appdynamicsMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				btHandle := appd.StartBT(r.URL.Path, "")
				next.ServeHTTP(w, r)
				appd.EndBT(btHandle)
				fmt.Printf("AppD Middleware sucessfully instrumented BT with handle: %x and URL.Path: %s\n", btHandle, r.URL.Path)
		})
}

func init() {
	stdprometheus.MustRegister(HTTPLatency)
	// flag.StringVar(&zip, "zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
	flag.StringVar(&port, "port", "8084", "Port on which to run")
	db.Register("mongodb", &mongodb.Mongo{})
}

func main() {

	flag.Parse()
	// Mechanical stuff.
	errc := make(chan error)

	// AppDynamics SDK initialization
	cfg := appd.Config{}
	cfg.AppName = os.Getenv("APPD_APPNAME") // "exampleapp"
	cfg.TierName = os.Getenv("APPD_TIERNAME")
	var (
		appd_hostname string
		appd_usessl bool
		appd_port uint64
		err error
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

	// AppDynamics Backend - MySQL Database
	backendName := "user-db"
	backendType := appd.APPD_BACKEND_DB
	backendProperties := map[string]string{
		"HOST":"user-db",
		"PORT":"27017",
		"DATABASE":"users",
		"VENDOR":"MongoDB",
		"VERSION":"3.0",
	}
	resolveBackend := false

	if err := appd.AddBackend(backendName, backendType, backendProperties, resolveBackend); err != nil {
		fmt.Printf("Error adding the AppDynamics backend\n")
	} else {
		fmt.Printf("Added AppDynamics backend successfully\n")
	}

	// Log domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Find service local IP.
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	// localAddr := conn.LocalAddr().(*net.UDPAddr)
	// host := strings.Split(localAddr.String(), ":")[0]
	defer conn.Close()

	var tracer stdopentracing.Tracer
	{
		// if zip == "" {
			tracer = stdopentracing.NoopTracer{}
		// } else {
		// 	logger := log.With(logger, "tracer", "Zipkin")
		// 	logger.Log("addr", zip)
		// 	collector, err := zipkin.NewHTTPCollector(
		// 		zip,
		// 		zipkin.HTTPLogger(logger),
		// 	)
		// 	if err != nil {
		// 		logger.Log("err", err)
		// 		os.Exit(1)
		// 	}
		// 	tracer, err = zipkin.NewTracer(
		// 		zipkin.NewRecorder(collector, false, fmt.Sprintf("%v:%v", host, port), ServiceName),
		// 	)
		// 	if err != nil {
		// 		logger.Log("err", err)
		// 		os.Exit(1)
		// 	}
		// }
		stdopentracing.InitGlobalTracer(tracer)
	}

	// AppD Start Backend Test
	btHandle := appd.StartBT("MongoDB Test", "")
	exitHandle := appd.StartExitcall(btHandle,"user-db")
	fmt.Printf("AppD Middleware sucessfully started BT with handle: %x and name: %s\n", btHandle, "MongoDB Test")
	fmt.Printf("AppD Middleware sucessfully started Exit Call with handle: %x and backend: %s\n", exitHandle, "user-db")

	dbconn := false
	for !dbconn {
		err := db.Init()
		if err != nil {
			if err == db.ErrNoDatabaseSelected {
				corelog.Fatal(err)
			}
			corelog.Print(err)
			// AppD Log Backend Error
			appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Unable to connect to Database", true)
		} else {
			dbconn = true
		}
	}

	// AppD Close Backend Test
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)

	fieldKeys := []string{"method"}
	// Service domain.
	var service api.Service
	{
		service = api.NewFixedService()
		service = api.LoggingMiddleware(logger)(service)
		service = api.NewInstrumentingService(
			kitprometheus.NewCounterFrom(
				stdprometheus.CounterOpts{
					Namespace: "microservices_demo",
					Subsystem: "user",
					Name:      "request_count",
					Help:      "Number of requests received.",
				},
				fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "microservices_demo",
				Subsystem: "user",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			service,
		)
	}

	// Endpoint domain.
	endpoints := api.MakeEndpoints(service, tracer)

	// HTTP router
	router := api.MakeHTTPHandler(endpoints, logger, tracer)

	// Inject AppDynamics Middleware
	router.Use(appdynamicsMiddleware)

	httpMiddleware := []commonMiddleware.Interface{
		commonMiddleware.Instrument{
			Duration:     HTTPLatency,
			RouteMatcher: router,
		},
	}

	// Handler
	handler := commonMiddleware.Merge(httpMiddleware...).Wrap(router)

	// Create and launch the HTTP server.
	go func() {
		logger.Log("transport", "HTTP", "port", port)
		errc <- http.ListenAndServe(fmt.Sprintf(":%v", port), handler)
	}()

	// Capture interrupts.
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("exit", <-errc)
}
