package services

import (
	"context"
	"errors"
	"time"

	"callflowmanager/internal/models"
	"callflowmanager/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repositories.NewUserRepository()}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.User, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrUserExists
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	role := req.Role
	if role == "" {
		role = "user"
	}
	user := &models.User{Email: req.Email, Password: string(hash), Name: req.Name, Role: role, Active: true}
	s.userRepo.Create(ctx, user)
	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.User, error) {
	user, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if user == nil {
		return nil, ErrUserNotFound
	}
	bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	user.Password = ""
	return user, nil
}

type AgentService struct {
	agentRepo *repositories.AgentRepository
}

func NewAgentService() *AgentService {
	return &AgentService{agentRepo: repositories.NewAgentRepository()}
}

func (s *AgentService) CreateAgent(ctx context.Context, agent *models.Agent) (*models.Agent, error) {
	err := s.agentRepo.Create(ctx, agent)
	return agent, err
}

func (s *AgentService) GetAllAgents(ctx context.Context) ([]models.Agent, error) {
	return s.agentRepo.FindAll(ctx)
}

type CustomerService struct {
	customerRepo *repositories.CustomerRepository
}

func NewCustomerService() *CustomerService {
	return &CustomerService{customerRepo: repositories.NewCustomerRepository()}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *models.Customer) (*models.Customer, error) {
	err := s.customerRepo.Create(ctx, customer)
	return customer, err
}

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	return s.customerRepo.FindAll(ctx)
}

type CallService struct {
	callRepo     *repositories.CallRepository
	agentRepo    *repositories.AgentRepository
	customerRepo *repositories.CustomerRepository
}

func NewCallService() *CallService {
	return &CallService{
		callRepo:     repositories.NewCallRepository(),
		agentRepo:    repositories.NewAgentRepository(),
		customerRepo: repositories.NewCustomerRepository(),
	}
}

func (s *CallService) CreateCall(ctx context.Context, call *models.Call) error {
	customer, _ := s.customerRepo.FindByID(ctx, call.CustomerID)
	if customer != nil {
		call.CustomerName = customer.Name
	}

	agent, _ := s.agentRepo.FindByID(ctx, call.AgentID)
	if agent != nil {
		call.AgentName = agent.Name
	}

	call.Status = models.CallScheduled
	return s.callRepo.Create(ctx, call)
}

func (s *CallService) GetCalls(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]models.Call, int64, error) {
	return s.callRepo.Find(ctx, startDate, endDate, page, limit)
}

func (s *CallService) UpdateCallStatus(ctx context.Context, id primitive.ObjectID, status models.CallStatus) error {
	calls, _, err := s.callRepo.Find(ctx, time.Time{}, time.Time{}, 1, 100)
	if err != nil {
		return err
	}

	for _, c := range calls {
		if c.ID == id {
			c.Status = status
			if status == models.CallInProgress {
				now := time.Now()
				c.StartedAt = &now
			}
			if status == models.CallCompleted {
				now := time.Now()
				c.EndedAt = &now
				if c.StartedAt != nil {
					c.Duration = int(now.Sub(*c.StartedAt).Minutes())
				}
			}
			return s.callRepo.Update(ctx, &c)
		}
	}
	return nil
}

func (s *CallService) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	return s.callRepo.GetStats(ctx)
}
