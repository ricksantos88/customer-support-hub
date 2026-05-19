package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ricksantos88/customer-support-hub/internal/database"
	"github.com/ricksantos88/customer-support-hub/internal/models"
)

func main() {
	dsn := buildDSN()
	db, err := database.NewPostgresConnection(dsn)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := seed(db); err != nil {
		log.Fatalf("seed database: %v", err)
	}

	fmt.Println("seed completed successfully")
}

func seed(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Contact{}, &models.Agent{}, &models.Conversation{}, &models.Message{}); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}

	agentOne, err := upsertAgent(db, "Alice Silva", "alice@example.com", "alice-jwt-secret")
	if err != nil {
		return err
	}

	agentTwo, err := upsertAgent(db, "Bruno Costa", "bruno@example.com", "bruno-jwt-secret")
	if err != nil {
		return err
	}

	contactOne, err := upsertContact(db, "+5511999990001", "Cliente Um")
	if err != nil {
		return err
	}

	contactTwo, err := upsertContact(db, "+5511999990002", "Cliente Dois")
	if err != nil {
		return err
	}

	conversationOne, err := upsertConversation(db, contactOne.ID, agentOne.ID, models.ConversationStatusOpen)
	if err != nil {
		return err
	}

	conversationTwo, err := upsertConversation(db, contactTwo.ID, agentTwo.ID, models.ConversationStatusPending)
	if err != nil {
		return err
	}

	if err := upsertMessage(db, conversationOne.ID, contactOne.ID, models.MessageDirectionInbound, "Olá, preciso de ajuda com meu pedido."); err != nil {
		return err
	}
	if err := upsertMessage(db, conversationOne.ID, agentOne.ID, models.MessageDirectionOutbound, "Olá! Vou verificar isso para você."); err != nil {
		return err
	}
	if err := upsertMessage(db, conversationTwo.ID, contactTwo.ID, models.MessageDirectionInbound, "Consigo falar com um atendente?"); err != nil {
		return err
	}
	if err := upsertMessage(db, conversationTwo.ID, agentTwo.ID, models.MessageDirectionOutbound, "Claro, estou assumindo seu atendimento."); err != nil {
		return err
	}

	return nil
}

func upsertAgent(db *gorm.DB, name, email, jwtSecret string) (*models.Agent, error) {
	agent := models.Agent{}
	result := db.Where("email = ?", strings.ToLower(email)).First(&agent)
	if result.Error == nil {
		return &agent, nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find agent: %w", result.Error)
	}

	agent = models.Agent{Name: name, Email: email, JWTSecret: jwtSecret}
	if err := db.Create(&agent).Error; err != nil {
		return nil, fmt.Errorf("create agent: %w", err)
	}
	return &agent, nil
}

func upsertContact(db *gorm.DB, phone, name string) (*models.Contact, error) {
	contact := models.Contact{}
	result := db.Where("phone = ?", phone).First(&contact)
	if result.Error == nil {
		return &contact, nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find contact: %w", result.Error)
	}

	contact = models.Contact{Phone: phone, Name: name}
	if err := db.Create(&contact).Error; err != nil {
		return nil, fmt.Errorf("create contact: %w", err)
	}
	return &contact, nil
}

func upsertConversation(db *gorm.DB, contactID, agentID uuid.UUID, status string) (*models.Conversation, error) {
	conversation := models.Conversation{}
	result := db.Where("contact_id = ? AND status <> ?", contactID, models.ConversationStatusClosed).First(&conversation)
	if result.Error == nil {
		return &conversation, nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("find conversation: %w", result.Error)
	}

	agentIDPtr := &agentID
	conversation = models.Conversation{ContactID: contactID, AssignedAgentID: agentIDPtr, Status: status}
	if err := db.Create(&conversation).Error; err != nil {
		return nil, fmt.Errorf("create conversation: %w", err)
	}
	return &conversation, nil
}

func upsertMessage(db *gorm.DB, conversationID, senderID uuid.UUID, direction, content string) error {
	var message models.Message
	result := db.Where("conversation_id = ? AND sender_id = ? AND direction = ? AND content = ?", conversationID, senderID, direction, content).First(&message)
	if result.Error == nil {
		return nil
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("find message: %w", result.Error)
	}

	message = models.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Direction:      direction,
		Content:        content,
	}
	if err := db.Create(&message).Error; err != nil {
		return fmt.Errorf("create message: %w", err)
	}
	return nil
}

func buildDSN() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "support")
	password := getEnv("DB_PASSWORD", "support123")
	dbname := getEnv("DB_NAME", "customer_support")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
