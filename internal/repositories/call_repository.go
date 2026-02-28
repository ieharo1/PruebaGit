package repositories

import (
	"context"
	"fmt"
	"time"

	"callflowmanager/internal/config"
	"callflowmanager/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	db := config.GetDatabase()
	return &UserRepository{collection: db.Database.Collection("users")}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)}})
	return err
}

type AgentRepository struct {
	collection *mongo.Collection
}

func NewAgentRepository() *AgentRepository {
	db := config.GetDatabase()
	return &AgentRepository{collection: db.Database.Collection("agents")}
}

func (r *AgentRepository) Create(ctx context.Context, agent *models.Agent) error {
	agent.CreatedAt = time.Now()
	agent.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, agent)
	if err != nil {
		return err
	}
	agent.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *AgentRepository) FindAll(ctx context.Context) ([]models.Agent, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var agents []models.Agent
	cursor.All(ctx, &agents)
	return agents, nil
}

func (r *AgentRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Agent, error) {
	var agent models.Agent
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&agent)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &agent, nil
}

func (r *AgentRepository) CreateIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{{Keys: bson.D{{Key: "email", Value: 1}}}})
	return err
}

type CustomerRepository struct {
	collection *mongo.Collection
}

func NewCustomerRepository() *CustomerRepository {
	db := config.GetDatabase()
	return &CustomerRepository{collection: db.Database.Collection("customers")}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, customer)
	if err != nil {
		return err
	}
	customer.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *CustomerRepository) FindAll(ctx context.Context) ([]models.Customer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var customers []models.Customer
	cursor.All(ctx, &customers)
	return customers, nil
}

func (r *CustomerRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Customer, error) {
	var customer models.Customer
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) CreateIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{{Keys: bson.D{{Key: "phone", Value: 1}}}})
	return err
}

type CallRepository struct {
	collection *mongo.Collection
}

func NewCallRepository() *CallRepository {
	db := config.GetDatabase()
	return &CallRepository{collection: db.Database.Collection("calls")}
}

func (r *CallRepository) Create(ctx context.Context, call *models.Call) error {
	call.CreatedAt = time.Now()
	call.UpdatedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, call)
	if err != nil {
		return err
	}
	call.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *CallRepository) Find(ctx context.Context, startDate, endDate time.Time, page, limit int) ([]models.Call, int64, error) {
	query := bson.M{}
	if !startDate.IsZero() && !endDate.IsZero() {
		query["scheduled_at"] = bson.M{"$gte": startDate, "$lte": endDate}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit

	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "scheduled_at", Value: -1}})
	cursor, err := r.collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var calls []models.Call
	if err := cursor.All(ctx, &calls); err != nil {
		return nil, 0, err
	}
	return calls, total, nil
}

func (r *CallRepository) Update(ctx context.Context, call *models.Call) error {
	call.UpdatedAt = time.Now()
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": call.ID}, bson.M{"$set": call})
	return err
}

func (r *CallRepository) GetStats(ctx context.Context) (*models.DashboardStats, error) {
	stats := &models.DashboardStats{CallsByStatus: make(map[string]int)}

	total, _ := r.collection.CountDocuments(ctx, bson.M{})
	stats.TotalCalls = int(total)

	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	todayCount, _ := r.collection.CountDocuments(ctx, bson.M{"scheduled_at": bson.M{"$gte": startOfDay, "$lt": endOfDay}})
	stats.TodayCalls = int(todayCount)

	pending, _ := r.collection.CountDocuments(ctx, bson.M{"status": models.CallScheduled})
	stats.PendingCalls = int(pending)

	completed, _ := r.collection.CountDocuments(ctx, bson.M{"status": models.CallCompleted})
	stats.CompletedCalls = int(completed)

	pipeline := []bson.M{{"$group": bson.M{"_id": "$status", "count": bson.M{"$sum": 1}}}}
	cursor, _ := r.collection.Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var result struct {
			ID    string
			Count int
		}
		cursor.Decode(&result)
		stats.CallsByStatus[result.ID] = result.Count
	}

	stats.ActiveAgents = 5

	return stats, nil
}

func (r *CallRepository) CreateIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "scheduled_at", Value: 1}}},
		{Keys: bson.D{{Key: "agent_id", Value: 1}}},
		{Keys: bson.D{{Key: "customer_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	return err
}
