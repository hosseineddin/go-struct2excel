package s2e

import (
	"os"
	"testing"
	"time"
)

// 1. Define a Struct with JSON tags representing business logic
type UserLog struct {
	ID        int       `json:"log_id"`               // Column 1
	Username  string    `json:"user_name"`            // Column 2
	IsActive  bool      `json:"active_status"`        // Column 3
	Secret    string    `json:"-"`                    // Will be IGNORED
	CreatedAt time.Time `json:"created_at,omitempty"` // Column 4: "created_at"
}

func TestEngine_GenerateXLSX(t *testing.T) {
	// 2. Generate dummy data
	var data []UserLog
	for i := 1; i <= 100; i++ {
		data = append(data, UserLog{
			ID:        i,
			Username:  "User_" + string(rune(i)),
			IsActive:  i%2 == 0,
			Secret:    "HiddenPassword123",
			CreatedAt: time.Now(),
		})
	}

	// 3. Create an output file (this file is ignored by .gitignore)
	file, err := os.Create("../../test_output.xlsx")
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// 4. Run the engine
	err = New[UserLog]().
		SetOutputWriter(file).
		SetFormat("xlsx").
		Generate(data)

	if err != nil {
		t.Fatalf("Engine failed: %v", err)
	}

	t.Log("Successfully generated test_output.xlsx! Check your root directory.")
}
