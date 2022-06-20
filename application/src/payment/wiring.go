package payment

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"

	"github.com/microservices-demo/payment/middleware"
	stdopentracing "github.com/opentracing/opentracing-go"

	// AppDynamics Go SDK Agent
	"fmt"
	appd "appdynamics"
)

//AppD middleware to enclose the routing handling functions and monitor Business Transactions
func appdynamicsMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hdr := r.Header.Get(appd.APPD_CORRELATION_HEADER_NAME)
				btHandle := appd.StartBT(r.URL.Path, hdr)
				next.ServeHTTP(w, r)
				appd.EndBT(btHandle)
				fmt.Printf("AppD Middleware sucessfully instrumented BT with handle: %x and URL.Path: %s\n", btHandle, r.URL.Path)
		})
}

func WireUp(ctx context.Context, declineAmount float32, tracer stdopentracing.Tracer, serviceName string) (http.Handler, log.Logger) {
	// Log domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	// Service domain.
	var service Service
	{
		service = NewAuthorisationService(declineAmount)
		service = LoggingMiddleware(logger)(service)
	}

	// Endpoint domain.
	endpoints := MakeEndpoints(service, tracer)

	router := MakeHTTPHandler(ctx, endpoints, logger, tracer)

	// Inject AppDynamics Middleware
	router.Use(appdynamicsMiddleware)

	httpMiddleware := []middleware.Interface{
		middleware.Instrument{
			Duration:     middleware.HTTPLatency,
			RouteMatcher: router,
			Service:      serviceName,
		},
	}

	// Handler
	handler := middleware.Merge(httpMiddleware...).Wrap(router)

	return handler, logger
}
