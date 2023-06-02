package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// struct for getHealthz and getReadyz functions
type healthCheckResponse struct {
	Status string `json:"status"`
}

func (s *Server) getHealthz(context *gin.Context) {
	err := s.store.Ping(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	context.JSON(http.StatusOK, healthCheckResponse{Status: "Ok"})
}

func (s *Server) getReadyz(context *gin.Context) {
	err := s.store.Ping(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	context.JSON(http.StatusOK, healthCheckResponse{Status: "Ok"})
}
