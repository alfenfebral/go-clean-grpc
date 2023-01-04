package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"go-clean-grpc/pkg/config"
	"go-clean-grpc/pkg/logger"
	pkgmongodb "go-clean-grpc/pkg/mongodb"
	pkgvalidator "go-clean-grpc/pkg/validator"
	todogrpcdelivery "go-clean-grpc/todo/delivery/grpc"
	todoproto "go-clean-grpc/todo/delivery/grpc/proto"
	todohttpdelivery "go-clean-grpc/todo/delivery/http"
	todorepository "go-clean-grpc/todo/repository"
	todoservice "go-clean-grpc/todo/service"
	responseutil "go-clean-grpc/utils/response"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger, // Log API request calls
		// middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	return router
}

// PrintAllRoutes - printing all routes
func PrintAllRoutes(router *chi.Mux) {
	logger.Println("===========================")
	logger.Println("REST API routes")
	logger.Println("===========================")
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logger.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Error(err)
	}
	logger.Println("===========================")
}

func main() {
	pkgvalidator.New()

	// Load environment variables
	err := config.LoadConfig()
	if err != nil {
		logger.Error(err)
	}

	// Init MongoDB
	_, cancel, client := pkgmongodb.InitMongoDB()
	defer cancel()

	go func() {
		startRESTServer(client)
	}()

	go func() {
		startGRPCServer(client)
	}()

	// catch shutdown
	done := make(chan bool, 1)
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		// graceful shutdown
		// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		// defer cancel()
		// server.GracefulStop(ctx, done)
	}()
	// wait for graceful shutdown
	<-done
}

func startRESTServer(client *mongo.Client) {
	router := Routes()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, responseutil.H{
			"success": "true",
			"code":    200,
			"message": "Services run properly",
		})
	})

	// Repository
	todoRepo := todorepository.New(client)

	// Service
	todoService := todoservice.New(todoRepo)

	// Delivery
	todoHandler := todohttpdelivery.New(todoService)
	todoHandler.RegisterRoutes(router)

	// Print
	PrintAllRoutes(router)

	addr := fmt.Sprintf("%s%s", ":", os.Getenv("REST_API_PORT"))
	logger.Info("REST API server started on port " + os.Getenv("REST_API_PORT"))
	err := http.ListenAndServe(addr, router)
	if err != nil {
		logger.Error(err)
	}
}

func startGRPCServer(client *mongo.Client) {
	server := grpc.NewServer()

	// Repository
	todoRepo := todorepository.New(client)
	// Service
	todoService := todoservice.New(todoRepo)
	// Delivery
	todoGrpcDelivery := todogrpcdelivery.New(todoService)

	reflection.Register(server)
	todoproto.RegisterTodoServer(server, todoGrpcDelivery)

	addr := fmt.Sprintf("%s%s", ":", os.Getenv("GRPC_PORT"))
	tl, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(err)
	}

	logger.Info("gRPC server started on port 8765")

	err = server.Serve(tl)
	if err != nil {
		logger.Error(err)
	}
}
