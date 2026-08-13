package dto

import "time"

type CreateAgentRequest struct {
	Name     string `json:"name"     example:"Carlos Souza"       description:"Agent name"`
	Email    string `json:"email"    example:"carlos@example.com" description:"Agent email address"`
	Password string `json:"password" example:"strong-password-123" description:"Agent password"`
	Role     string `json:"role"     example:"agent"              description:"Agent role: admin or agent"`
}

type UpdateAgentRequest struct {
	Name     string `json:"name"     example:"Carlos Souza Updated" description:"Agent name"`
	Email    string `json:"email"    example:"carlos@example.com"   description:"Agent email address"`
	Password string `json:"password" example:"new-password-123"    description:"Agent password (optional)"`
	Role     string `json:"role"     example:"admin"                description:"Agent role: admin or agent"`
}

type AgentResponse struct {
	ID         string    `json:"id"          format:"uuid"                description:"Agent UUID"`
	Name       string    `json:"name"        example:"Carlos Souza"       description:"Agent name"`
	Email      string    `json:"email"       example:"carlos@example.com" description:"Agent email address"`
	Role       string    `json:"role"        example:"agent"              description:"Agent role: admin or agent"`
	CreatedAt  time.Time `json:"created_at"                               description:"Agent creation timestamp"`
	LastActive time.Time `json:"last_active"                              description:"Agent last active timestamp"`
}

type SessionResponse struct {
	ID        string    `json:"id"         format:"uuid"                description:"Session UUID"`
	AgentID   string    `json:"agent_id"   format:"uuid"                description:"Agent UUID"`
	AgentName string    `json:"agent_name" example:"Alice Silva"        description:"Agent name"`
	IPAddress string    `json:"ip_address" example:"127.0.0.1"          description:"IP address of the session"`
	UserAgent string    `json:"user_agent" example:"Mozilla/5.0..."     description:"User agent of the session"`
	CreatedAt time.Time `json:"created_at"                              description:"Session creation timestamp"`
	ExpiresAt time.Time `json:"expires_at"                              description:"Session expiration timestamp"`
}

type SystemStatusResponse struct {
	Database            string `json:"database"             example:"connected"  description:"Database connection status"`
	Redis               string `json:"redis"                example:"connected"  description:"Redis connection status"`
	WhatsAppIntegration string `json:"whatsapp_integration" example:"configured" description:"WhatsApp API configuration status"`
}
