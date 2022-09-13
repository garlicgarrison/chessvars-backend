package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	firebase "firebase.google.com/go/v4"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"github.com/pafkiuq/backend/graph"
	"github.com/pafkiuq/backend/graph/generated"
	"github.com/pafkiuq/backend/middleware"
	"github.com/pafkiuq/backend/pkg/firestore"
	"github.com/pafkiuq/backend/pkg/users"
)

type Config struct {
	Port    int    `envconfig:"PORT" default:"8080"`
	Address string `envconfig:"ADDRESS" default:"http://localhost:8080"`

	Firestore firestore.Config
}

func main() {
	ctx := context.Background()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		fmt.Printf("failed to process configs: %s\n", err)
		os.Exit(1)
	}

	/* start section: third party */

	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Printf("error in initializing firebase: %s\n", err)
		os.Exit(1)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error in initializing firebase auth: %s\n", err)
		os.Exit(1)
	}

	fs, err := firestore.NewClient(ctx, &cfg.Firestore)
	if err != nil {
		log.Printf("error in intitializing firestore: %s \n", err)
		os.Exit(1)
	}

	/* end section: third party */

	/* start section: initialize server */

	users, err := users.NewService(users.Config{
		Firestore: fs,
	})
	if err != nil {
		fmt.Printf("failed to init users service: %s", err)
		os.Exit(1)
	}

	resolver, err := graph.NewResolver(graph.Config{
		UsersService: users,
	})
	if err != nil {
		fmt.Printf("failed to init resolver: %s\n", err)
		os.Exit(1)
	}

	graphql := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	/* end section: initialize server */

	/* start section: register routes */

	mux := http.NewServeMux()

	mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	mux.HandleFunc("/explorer", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r,
			fmt.Sprintf("https://sandbox.apollo.dev/?endpoint=%s", cfg.Address+"/graphql"),
			http.StatusSeeOther,
		)
	})
	mux.Handle("/graphql", middleware.NewAuth(client, graphql))

	/* end section: register routes */

	handler := middleware.NewRecover(
		middleware.NewLogger(
			middleware.NewCors(
				mux,
			),
		),
	)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	done := make(chan struct{})
	go func() {
		close(done)

		<-signalChan
		log.Printf("signaling server shutdown\n")
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error in signaling shutdown: %s\n", err)
			return
		}
		log.Printf("successfully signaled shutdown")
	}()

	log.Printf("connect to http://localhost:%d/ for GraphQL playground\n", cfg.Port)
	log.Printf("connect to http://localhost:%d/explorer for Apollo Explorer\n", cfg.Port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("error in running server: %s\n", err)
		os.Exit(1)
	}
	<-done
	log.Printf("successfully shutdown server :)\n")
}
