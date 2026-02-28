package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Name      string             `bson:"name" json:"name"`
	Role      string             `bson:"role" json:"role"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Agent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Email       string             `bson:"email" json:"email"`
	Phone       string             `bson:"phone" json:"phone"`
	Department  string             `bson:"department" json:"department"`
	Status      string             `bson:"status" json:"status"`
	Active      bool               `bson:"active" json:"active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Customer struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Email       string             `bson:"email" json:"email"`
	Phone       string             `bson:"phone" json:"phone"`
	Company     string             `bson:"company" json:"company"`
	Notes       string             `bson:"notes" json:"notes"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CallStatus string

const (
	CallScheduled CallStatus = "scheduled"
	CallInProgress CallStatus = "in_progress"
	CallCompleted CallStatus = "completed"
	CallMissed    CallStatus = "missed"
	CallCancelled CallStatus = "cancelled"
)

type Call struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID  primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	CustomerName string            `bson:"customer_name" json:"customer_name"`
	AgentID     primitive.ObjectID `bson:"agent_id" json:"agent_id"`
	AgentName   string             `bson:"agent_name" json:"agent_name"`
	ScheduledAt time.Time          `bson:"scheduled_at" json:"scheduled_at"`
	StartedAt   *time.Time         `bson:"started_at,omitempty" json:"started_at,omitempty"`
	EndedAt     *time.Time         `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
	Status      CallStatus         `bson:"status" json:"status"`
	Notes       string             `bson:"notes" json:"notes"`
	Duration    int                `bson:"duration" json:"duration"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type CallLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CallID      primitive.ObjectID `bson:"call_id" json:"call_id"`
	AgentID     primitive.ObjectID `bson:"agent_id" json:"agent_id"`
	Action      string             `bson:"action" json:"action"`
	Description string             `bson:"description" json:"description"`
	Timestamp   time.Time          `bson:"timestamp" json:"timestamp"`
}

type Schedule struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AgentID     primitive.ObjectID `bson:"agent_id" json:"agent_id"`
	DayOfWeek   int                `bson:"day_of_week" json:"day_of_week"`
	StartTime   string             `bson:"start_time" json:"start_time"`
	EndTime     string             `bson:"end_time" json:"end_time"`
	Active      bool               `bson:"active" json:"active"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

type DashboardStats struct {
	TotalCalls     int              `json:"total_calls"`
	TodayCalls     int              `json:"today_calls"`
	PendingCalls   int              `json:"pending_calls"`
	CompletedCalls int              `json:"completed_calls"`
	ActiveAgents   int              `json:"active_agents"`
	CallsByStatus  map[string]int   `json:"calls_by_status"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
