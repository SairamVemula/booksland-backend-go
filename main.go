package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SairamVemula/booksland-backend-go/pkg/models"
	"github.com/SairamVemula/booksland-backend-go/pkg/routes"
	"github.com/SairamVemula/booksland-backend-go/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	logger := utils.NewLogger()
	configs := utils.NewConfigurations(logger)

	// validator contains all the methods that are need to validate the user json in request
	validator := models.NewValidation()

	models.Connect(configs.MONGO_URI, configs.DBName, logger)
	var dir string
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	flag.StringVar(&dir, "dir", "./assets", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(dir))))

	//Main routes
	routes.RegisterUsersRoutes(r, logger, configs, validator)
	routes.RegisterAuthToutes(r, logger, configs, validator)
	routes.RegisterCourseRoutes(r, logger, configs, validator)
	routes.RegisterBooksRoutes(r, logger, configs, validator)
	routes.RegisterMediaRoutes(r, logger, configs, validator)
	routes.RegisterStocksRoutes(r, logger, configs, validator)
	routes.RegisterFeedsRoutes(r, logger, configs, validator)

	srv := &http.Server{
		Addr: configs.ServerAddress,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      cors.AllowAll().Handler(r), // Pass our instance of gorilla/mux in.

	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		logger.Debug("Listening on ", configs.ServerAddress)
		if err := srv.ListenAndServe(); err != nil {
			logger.Error(err.Error(), err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Debug("gracefully shut down")
	os.Exit(0)
}

// func main() {
// 	r := mux.NewRouter()
// 	cors := cors.New(cors.Options{
// 		AllowedOrigins:         []string{"*"},
// 		AllowOriginRequestFunc: func(r *http.Request, origin string) bool { return true },
// 		AllowedMethods:         []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
// 		AllowedHeaders:         []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
// 		ExposedHeaders:         []string{"Link"},
// 		AllowCredentials:       true,
// 		OptionsPassthrough:     true,
// 		MaxAge:                 3599, // Maximum value not ignored by any of major browsers
// 	})

// 	r.Use(cors.Handler)

// 	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		w.Write([]byte("Hello World"))
// 	})
// 	http.ListenAndServe(":8080", r)
// }
