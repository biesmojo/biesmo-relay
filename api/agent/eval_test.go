//go:build !integration

package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock Anthropic client for unit tests
type MockClient struct {
	Response anthropic.MessagesCreateResponse
	Error    error
}

func (m *MockClient) MessagesCreate(ctx context.Context, req *anthropic.MessagesCreateRequest) (*anthropic.MessagesCreateResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return &m.Response, nil
}

func TestAgentHappyInquiry(t *testing.T) {
	pool := &pgxpool.Pool{} // mock pool
	client := &MockClient{
		Response: anthropic.MessagesCreateResponse{
			Content: []anthropic.ContentBlock{
				{
					Type: "text",
					Text: "Oli mesin X mengandung aditif anti-wear dan detergen.",
				},
				{
					Type: "tool_use",
					Name: "search_kb",
					Input: map[string]interface{}{"query": "oli mesin X"},
				},
			},
		},
	}
	// Setup mock tool responses etc
	// Test POST /api/chat expect search_kb called, Indonesian reply, sentiment Positif
}

func TestAgentComplaint(t *testing.T) {
	// Similar, expect create_ticket, sentiment Negatif/Sangat Negatif
}

func TestAgentEscalation(t *testing.T) {
	// Multi-turn vague → escalate after 1 clarify
}

// ... other fixtures

// //go:build integration
// func TestAgentIntegration(t *testing.T) {
// 	// Real API calls with real Supabase/Anthropic
// }

type Fixture struct {
	Name string
	Input []string
	ExpectedTools []string
	ExpectedSentiment string
}

fixtures := []Fixture{
	{
		Name: "Happy inquiry",
		Input: []string{"Halo, saya mau tanya kandungan oli mesin X"},
		ExpectedTools: []string{"search_kb"},
		ExpectedSentiment: "Positif",
	},
	{
		Name: "Complaint",
		Input: []string{"Produk kamu rusak! Saya sangat kecewa dan minta refund segera!"},
		ExpectedTools: []string{"create_ticket"},
		ExpectedSentiment: "Sangat Negatif",
	},
	{
		Name: "Human request",
		Input: []string{"Saya mau bicara dengan manusia saja"},
		ExpectedTools: []string{"escalate_to_human"},
	},
	{
		Name: "Lookup",
		Input: []string{"Cek status saya, nomor hp saya 08123456789"},
		ExpectedTools: []string{"lookup_customer"},
	},
}

func TestAgentFixtures(t *testing.T) {
	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			// Setup DB fixtures
			// Mock Anthropic sequence
			// Call handler
			// Assert tools called match
			// Assert reply contains Indonesian
			// Assert session sentiment matches
		})
	}
}

