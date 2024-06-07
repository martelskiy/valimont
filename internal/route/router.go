package route

import (
	"net/http"

	_ "github.com/martelskiy/valimont/api/docs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Router struct {
	router *http.ServeMux
}

func NewRouter() *Router {
	router := http.NewServeMux()
	return &Router{
		router: router,
	}
}

func (r *Router) WithAPIDocumentation() *Router {
	r.router.Handle("/swagger/*", httpSwagger.WrapHandler)
	return r
}

// @Summary	Prometheus metrics
// @Tags		status
// @Accept		json
// @Produce	json
// @Success	200
// @Router		/metrics [get]
func (r *Router) WithPrometheusMetrics() *Router {
	r.router.Handle("/metrics", promhttp.Handler())
	return r
}

func (r *Router) WithRoute(route Route) *Router {
	r.router.HandleFunc(route.name, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != string(route.httpVerb) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		h := route.Handler()
		h(w, r)
	})
	return r
}

func (r *Router) GetRouter() *http.ServeMux {
	return r.router
}
