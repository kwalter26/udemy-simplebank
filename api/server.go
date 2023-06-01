package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/token"
	"github.com/kwalter26/udemy-simplebank/util"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	app        *newrelic.Application
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maketer: %w", err)
	}

	server := &Server{store: store, tokenMaker: maker, config: config}

	if config.NewRelicAppEnabled {
		app, err := newrelic.NewApplication(
			newrelic.ConfigAppName(config.NewRelicAppName),
			newrelic.ConfigLicense(config.NewRelicLicenseKey),
			newrelic.ConfigAppLogForwardingEnabled(config.NewRelicLogForwardingEnabled),
			newrelic.ConfigDistributedTracerEnabled(config.NewRelicDistributedTracingEnabled),
		)
		if err != nil {
			return nil, fmt.Errorf("cannot create new relic app: %w", err)
		}
		server.app = app
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	router.Use(nrgin.Middleware(s.app))
	router.POST("/users", s.CreateUser)
	router.POST("/users/login", s.loginUser)
	router.GET("readyz", s.getReadyz)
	router.GET("healthz", s.getHealthz)

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts", s.listAccounts)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.PUT("/accounts/:id", s.updateAccount)
	authRoutes.DELETE("/accounts/:id", s.deleteAccount)

	authRoutes.POST("/transfers", s.createTransfer)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
