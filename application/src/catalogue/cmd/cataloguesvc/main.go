package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	// zipkin "github.com/openzipkin/zipkin-go-opentracing"

	"net/http"

	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/microservices-demo/catalogue"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weaveworks/common/middleware"
	"golang.org/x/net/context"

	// AppDynamics Go SDK Agent
	"strconv"
	appd "appdynamics"
)

const (
	ServiceName = "catalogue"
)

var (
	HTTPLatency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Time (in seconds) spent serving HTTP requests.",
		Buckets: prometheus.DefBuckets,
	}, []string{"service", "method", "route", "status_code"})
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
	// Prometheus
	prometheus.MustRegister(HTTPLatency)
}

func main() {
	var (
		port   = flag.String("port", "8081", "Port to bind HTTP listener") // TODO(pb): should be -addr, default ":8081"
		images = flag.String("images", "./images/", "Image path")
		dsn    = flag.String("DSN", "catalogue_user:default_password@tcp(catalogue-db:3306)/socksdb", "Data Source Name: [username[:password]@][protocol[(address)]]/dbname")
		// zip    = flag.String("zipkin", os.Getenv("ZIPKIN"), "Zipkin address")
	)
	flag.Parse()

	fmt.Fprintf(os.Stderr, "images: %q\n", *images)
	abs, err := filepath.Abs(*images)
	fmt.Fprintf(os.Stderr, "Abs(images): %q (%v)\n", abs, err)
	pwd, err := os.Getwd()
	fmt.Fprintf(os.Stderr, "Getwd: %q (%v)\n", pwd, err)
	files, _ := filepath.Glob(*images + "/*")
	fmt.Fprintf(os.Stderr, "ls: %q\n", files) // contains a list of all files in the current directory

	// Mechanical stuff.
	errc := make(chan error)
	ctx := context.Background()

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

	// AppDynamics Backend - MySQL Database
	backendName := "catalogue-db"
	backendType := "APPD_BACKEND_DB"
	backendProperties := map[string]string{
		"HOST":"catalogue-db",
		"PORT":"3306",
		"DATABASE":"MySQL",
		"VENDOR":"MySQL",
		"VERSION":"5.7",
	}
	resolveBackend := false

	if err := appd.AddBackend(backendName, backendType, backendProperties, resolveBackend); err != nil {
		fmt.Printf("Error adding the AppDynamics MySQL backend\n")
	} else {
		fmt.Printf("Added AppDynamics MySQL backend successfully\n")
	}

	// Log domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	var tracer stdopentracing.Tracer
	{
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

	// Data domain.

	// AppD Backend Test
	btHandle := appd.StartBT("MySQL Test", "")
	exitHandle := appd.StartExitcall(btHandle,"catalogue-db")
	fmt.Printf("AppD Middleware sucessfully started BT with handle: %x and name: %s\n", btHandle, "MySQL Test")
	fmt.Printf("AppD Middleware sucessfully started Exit Call with handle: %x and backend: %s\n", exitHandle, "catalogue-db")

	db, err := sqlx.Open("mysql", *dsn)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	defer db.Close()

	// Check if DB connection can be made, only for logging purposes, should not fail/exit
	err = db.Ping()
	if err != nil {
		logger.Log("Error", "Unable to connect to Database", "DSN", dsn)
		appd.AddExitcallError(exitHandle, appd.APPD_LEVEL_ERROR, "Unable to connect to Database", true)
	}
	appd.EndExitcall(exitHandle)
	appd.EndBT(btHandle)


	// Service domain.
	var service catalogue.Service
	{
		service = catalogue.NewCatalogueService(db, logger)
		service = catalogue.LoggingMiddleware(logger)(service)
	}

	// Endpoint domain.
	endpoints := catalogue.MakeEndpoints(service, tracer)

	// HTTP router
	router := catalogue.MakeHTTPHandler(ctx, endpoints, *images, logger, tracer)

	router.Use(appdynamicsMiddleware)

	httpMiddleware := []middleware.Interface{
		middleware.Instrument{
			Duration:     HTTPLatency,
			RouteMatcher: router,
		},
	}

	// Handler
	handler := middleware.Merge(httpMiddleware...).Wrap(router)

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

	// Exit AppD Agent
	// appd.TerminateSDK()

	logger.Log("exit", <-errc)
}
