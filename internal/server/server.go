package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

type Server struct {
	router  *fiber.App
	api     huma.API
	handler *APIHandler
	port    int
}

func NewServer(db *gorm.DB, port int) *Server {
	handler := &APIHandler{DB: db}

	router := fiber.New()

	api := humafiber.New(
		router, huma.DefaultConfig("A server to track file metadata discovered by agents.", "1.0.0"),
	)

	server := &Server{
		router:  router,
		api:     api,
		handler: handler,
		port:    port,
	}
	server.configureMiddleware()
	server.configureRoutes()

	return server
}

func (s *Server) configureMiddleware() {
	s.router.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05.999Z0700",
		TimeZone:   "Local",
		Format:     "${time} [INFO] ${locals:requestid} ${method} ${path} ${status} ${latency} ${error}​\n",
	}))

	s.router.Use(healthcheck.New())
	s.router.Use(helmet.New())

	s.router.Use(requestid.New())

	prom := fiberprometheus.NewWith("filo", "filo", "http")
	prom.RegisterAt(s.router, "/metrics")
	s.router.Use(prom.Middleware)

	s.router.Get("/service/metrics", monitor.New())
	s.router.Use(recover.New())
}

func (s *Server) configureRoutes() {
	huma.Register(s.api, huma.Operation{
		OperationID: "register-file",
		Summary:     "Register a discovered file",
		Description: "Create a new file record or update an existing one (based on directory path and filename).",
		Method:      fiber.MethodPost,
		Path:        "/files",
	}, s.handler.RegisterFile)

	huma.Register(s.api, huma.Operation{
		OperationID: "list-files",
		Summary:     "List discovered files",
		Description: "Get a paginated list of file records, with optional filters.",
		Method:      fiber.MethodGet,
		Path:        "/files",
	}, s.handler.ListFiles)

	huma.Register(s.api, huma.Operation{
		OperationID: "get-file-by-id",
		Summary:     "Get file by ID",
		Description: "Get a single file record by its unique database ID.",
		Method:      fiber.MethodGet,
		Path:        "/files/{id}",
	}, s.handler.GetFileByID)

	huma.Register(s.api, huma.Operation{
		OperationID: "delete-file",
		Summary:     "Delete a file record",
		Description: "Delete a file record using its directory path and filename.",
		Method:      fiber.MethodDelete,
		Path:        "/files",
	}, s.handler.DeleteFile)
}

func (s *Server) Start() error {
	serverErr := make(chan error, 1)

	go func() {
		if err := s.router.Listen(fmt.Sprintf(":%d", s.port)); !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	select {
	case err := <-serverErr:
		return err
	case sig := <-stop:
		fmt.Printf("Received signal: %s\n", sig)
		fmt.Println("Shutting down server...")
		err := s.router.Shutdown()
		if err != nil {
			fmt.Printf("Error shutting down server: %v\n", err)
		}
		done <- true
	}

	<-done
	fmt.Println("Server stopped")
	return nil
}
