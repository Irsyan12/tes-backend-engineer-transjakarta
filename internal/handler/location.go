package handler

import (
	"net/http"
	"strconv"

	"fleet-management/internal/model"
	"fleet-management/internal/service"

	"github.com/gin-gonic/gin"
)

type LocationAPIHandler struct {
	service service.LocationService
}

func NewLocationAPIHandler(s service.LocationService) *LocationAPIHandler {
	return &LocationAPIHandler{service: s}
}

func (h *LocationAPIHandler) GetLatestLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")

	loc, err := h.service.GetLatestByVehicleID(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get location or not found"})
		return
	}

	c.JSON(http.StatusOK, loc)
}

func (h *LocationAPIHandler) GetHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	start, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start timestamp"})
		return
	}

	end, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end timestamp"})
		return
	}

	history, err := h.service.GetHistory(c.Request.Context(), vehicleID, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get history"})
		return
	}

	if history == nil {
		history = make([]model.VehicleLocation, 0)
	}

	c.JSON(http.StatusOK, history)
}
