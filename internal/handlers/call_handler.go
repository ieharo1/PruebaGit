package handlers

import (
	"net/http"
	"strconv"
	"time"

	"callflowmanager/internal/config"
	"callflowmanager/internal/models"
	"callflowmanager/internal/services"
	"callflowmanager/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authService *services.AuthService
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authService: authService, config: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role, h.config.JWTSecret, h.config.JWTExpiryHours)
	c.JSON(http.StatusCreated, models.AuthResponse{Token: token, User: *user})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, _ := utils.GenerateToken(user.ID, user.Email, user.Role, h.config.JWTSecret, h.config.JWTExpiryHours)
	c.JSON(http.StatusOK, models.AuthResponse{Token: token, User: *user})
}

type AgentHandler struct {
	agentService *services.AgentService
}

func NewAgentHandler(agentService *services.AgentService) *AgentHandler {
	return &AgentHandler{agentService: agentService}
}

func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var agent models.Agent
	c.ShouldBindJSON(&agent)
	created, _ := h.agentService.CreateAgent(c.Request.Context(), &agent)
	c.JSON(http.StatusCreated, created)
}

func (h *AgentHandler) GetAgents(c *gin.Context) {
	agents, _ := h.agentService.GetAllAgents(c.Request.Context())
	c.JSON(http.StatusOK, agents)
}

type CustomerHandler struct {
	customerService *services.CustomerService
}

func NewCustomerHandler(customerService *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var customer models.Customer
	c.ShouldBindJSON(&customer)
	created, _ := h.customerService.CreateCustomer(c.Request.Context(), &customer)
	c.JSON(http.StatusCreated, created)
}

func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	customers, _ := h.customerService.GetAllCustomers(c.Request.Context())
	c.JSON(http.StatusOK, customers)
}

type CallHandler struct {
	callService *services.CallService
}

func NewCallHandler(callService *services.CallService) *CallHandler {
	return &CallHandler{callService: callService}
}

func (h *CallHandler) CreateCall(c *gin.Context) {
	var call models.Call
	c.ShouldBindJSON(&call)
	h.callService.CreateCall(c.Request.Context(), &call)
	c.JSON(http.StatusCreated, call)
}

func (h *CallHandler) GetCalls(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	var startDate, endDate time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		startDate, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr := c.Query("end_date"); endStr != "" {
		endDate, _ = time.Parse("2006-01-02", endStr)
	}
	calls, total, _ := h.callService.GetCalls(c.Request.Context(), startDate, endDate, page, limit)
	c.JSON(http.StatusOK, gin.H{"data": calls, "total": total, "page": page, "limit": limit})
}

func (h *CallHandler) UpdateStatus(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var req struct{ Status string `json:"status"` }
	c.ShouldBindJSON(&req)
	h.callService.UpdateCallStatus(c.Request.Context(), id, models.CallStatus(req.Status))
	c.JSON(http.StatusOK, gin.H{"message": "Status updated"})
}

func (h *CallHandler) GetStats(c *gin.Context) {
	stats, _ := h.callService.GetStats(c.Request.Context())
	c.JSON(http.StatusOK, stats)
}

var _ = primitive.ObjectID{}
