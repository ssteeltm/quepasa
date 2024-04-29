package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	models "github.com/nocodeleaks/quepasa/models"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// swagger embed files
	httpSwagger "github.com/swaggo/http-swagger"
)

func QPWebServerStart() error {
	r := newRouter()
	webAPIPort := os.Getenv(models.ENV_WEBAPIPORT)
	webAPIHost := os.Getenv(models.ENV_WEBAPIHOST)
	if len(webAPIPort) == 0 {
		webAPIPort = "31000"
	}

	var timeout = 30 * time.Second
	server := http.Server{
		Addr:         webAPIHost + ":" + webAPIPort,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Handler:      r,
	}

	log.Infof("starting web server on port: %s", webAPIPort)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func NormalizePathsToLower(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "" {
			r.URL.Path = strings.ToLower(r.URL.Path)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func newRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(NormalizePathsToLower)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	if models.ENV.HttpLogs() {
		r.Use(middleware.Logger)
	}

	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// web routes
	// authenticated web routes
	r.Group(RegisterFormAuthenticatedControllers)

	// unauthenticated web routes
	r.Group(RegisterFormControllers)

	// api routes
	addAPIRoutes(r)

	// static files
	workDir, _ := os.Getwd()
	assetsDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(assetsDir))

	// Swagger Ui
	ServeSwaggerUi(r)

	// Metrics
	ServeMetrics(r)
	return r
}

func addAPIRoutes(r chi.Router) {
	r.Group(RegisterAPIControllers)
	r.Group(RegisterAPIV2Controllers)
	r.Group(RegisterAPIV3Controllers)
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func ServeSwaggerUi(r chi.Router) {
	log.Debug("starting swaggerUi service")
	r.Mount("/swagger", httpSwagger.WrapHandler)
}

func ServeMetrics(r chi.Router) {
	log.Debug("starting metrics service")
	r.Handle("/metrics", promhttp.Handler())
}
