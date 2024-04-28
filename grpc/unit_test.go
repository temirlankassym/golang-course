package test

import (
	"course/grpc/client"
	"testing"
)

func TestClient(t *testing.T) {
	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Test AddUser",
			test: func(t *testing.T) {
				_, err := c.AddUser("Test", "test@test.com")
				if err != nil {
					t.Errorf("AddUser failed: %v", err)
				}
			},
		},
		{
			name: "Test GetUser",
			test: func(t *testing.T) {
				_, err := c.GetUser(1)
				if err != nil {
					t.Errorf("GetUser failed: %v", err)
				}
			},
		},
		{
			name: "Test ListUsers",
			test: func(t *testing.T) {
				_, err := c.ListUsers()
				if err != nil {
					t.Errorf("ListUsers failed: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
