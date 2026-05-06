package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goshop/internal/notification/service"
	"goshop/pkg/response"
)

type Handler struct {
	svc service.PreferenceService
}

func NewHandler(svc service.PreferenceService) *Handler {
	return &Handler{svc: svc}
}

// ListPreferences godoc
//
//	@Summary	List notification preferences for the current user
//	@Tags		notifications
//	@Produce	json
//	@Success	200	{object}	response.Response
//	@Router		/me/notification-preferences [get]
//	@Security	ApiKeyAuth
func (h *Handler) ListPreferences(c *gin.Context) {
	userID := c.GetString("userId")
	prefs, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err, "list preferences")
		return
	}
	response.JSON(c, http.StatusOK, prefs)
}

type setPreferenceRequest struct {
	EventType string `json:"event_type" binding:"required"`
	Channel   string `json:"channel" binding:"required"`
	Enabled   bool   `json:"enabled"`
}

// SetPreference godoc
//
//	@Summary	Set a single notification preference
//	@Tags		notifications
//	@Accept		json
//	@Produce	json
//	@Param		body	body		setPreferenceRequest	true	"preference"
//	@Success	200		{object}	response.Response
//	@Router		/me/notification-preferences [put]
//	@Security	ApiKeyAuth
func (h *Handler) SetPreference(c *gin.Context) {
	var req setPreferenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err, "invalid body")
		return
	}
	userID := c.GetString("userId")
	pref, err := h.svc.Set(c.Request.Context(), userID, req.EventType, req.Channel, req.Enabled)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err, "set preference")
		return
	}
	response.JSON(c, http.StatusOK, pref)
}
