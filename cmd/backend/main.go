package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/garlicgarrison/chessvars-backend/graph"
	"github.com/garlicgarrison/chessvars-backend/graph/generated"
	"github.com/garlicgarrison/chessvars-backend/graph/resolver"
	"github.com/garlicgarrison/chessvars-backend/middleware"
	"github.com/garlicgarrison/chessvars-backend/pkg/elo"
	"github.com/garlicgarrison/chessvars-backend/pkg/firestore"
	"github.com/garlicgarrison/chessvars-backend/pkg/game"
	"github.com/garlicgarrison/chessvars-backend/pkg/users"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
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
	fmt.Printf("config: %v", cfg)

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

	elo, err := elo.NewService(elo.Config{
		Firestore: fs,
	})
	if err != nil {
		fmt.Printf("failed to init users service: %s", err)
		os.Exit(1)
	}

	game, err := game.NewService(game.Config{
		Firestore:  fs,
		EloService: elo,
	})
	if err != nil {
		fmt.Printf("failed to init users service: %s", err)
		os.Exit(1)
	}

	resolver, err := graph.NewResolver(graph.Config{
		Services: &resolver.Services{
			Users: users,
			Game:  game,
			Elo:   elo,
		},
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

	graphql.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				log.Printf("checkorigin %v", true)
				return true
			},
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			log.Printf("init payload %v", initPayload)
			return ctx, nil
		},
	})

	/* end section: initialize server */

	/* start section: register routes */

	mux := http.NewServeMux()

	// mux.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	mux.HandleFunc("/explorer", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r,
			fmt.Sprintf("https://sandbox.apollo.dev/?endpoint=%s", cfg.Address+"/graphql"),
			http.StatusSeeOther,
		)
	})
	mux.Handle("/graphql",
		middleware.NewLogger(
			middleware.NewCors(
				middleware.NewAuth(client, graphql),
			),
		),
	)
	mux.Handle("/subscriptions", graphql)

	/* end section: register routes */

	handler := middleware.NewRecover(
		mux,
		// middleware.NewLogger(
		// 	middleware.NewCors(
		// 		mux,
		// 	),
		// ),
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
