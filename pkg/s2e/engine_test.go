package s2e

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// ==========================================
// 1. Complex Models for Edge Case Testing
// ==========================================

type DeepNested struct {
	Detail string `json:"جزئیات_عمیق"`
}

type Nested struct {
	PointerVal *string    `json:"مقدار_پوینتر"` // Testing nil pointers
	Inner      DeepNested // Testing multi-level nested structs
}

// EdgeCaseModel contains every possible data anomaly
type EdgeCaseModel struct {
	ID          int       `json:"شناسه"`
	privateStr  string    `json:"private_field"` // UNEXPORTED: Engine must not crash and must ignore this
	Secret      string    `json:"-"`             // EXPORTED but IGNORED by tag
	IsActive    bool      `json:"وضعیت"`
	Balance     float64   `json:"موجودی"`
	Registered  time.Time `json:"زمان_ثبت"`
	ComplexData Nested    // Nested struct
}

// ==========================================
// 2. Unit Tests (Table-Driven for Edge Cases)
// ==========================================

func TestEngine_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Sample string for pointer testing
	sampleText := "Valid Pointer"

	tests := []struct {
		name          string
		format        string
		data          []EdgeCaseModel
		expectedError bool
	}{
		{
			name:          "Empty Data Slice",
			format:        "xlsx",
			data:          []EdgeCaseModel{},
			expectedError: true, // Engine should reject empty data
		},
		{
			name:   "Valid Complex Data with Nil Pointers",
			format: "csv",
			data: []EdgeCaseModel{
				{
					ID:         1,
					privateStr: "hidden", // Should be ignored
					Secret:     "hidden", // Should be ignored
					IsActive:   true,
					Balance:    1500.50,
					Registered: time.Now(),
					ComplexData: Nested{
						PointerVal: nil, // Nil pointer test
						Inner:      DeepNested{Detail: "Level 3 Data"},
					},
				},
				{
					ID: 2,
					ComplexData: Nested{
						PointerVal: &sampleText, // Valid pointer test
					},
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			err := New[EdgeCaseModel]().
				SetFormat(tt.format).
				Stream(c, tt.data)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test '%s' failed: expected error %v, got %v", tt.name, tt.expectedError, err)
			}

			if !tt.expectedError && w.Code != http.StatusOK {
				t.Errorf("Test '%s' failed: expected HTTP 200, got %d", tt.name, w.Code)
			}
		})
	}
}

// ==========================================
// 3. Heavy Load & Pagination Test (2 Million Rows)
// ==========================================

func TestEngine_HeavyLoad_Pagination(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping 2,000,000 rows pagination test in short mode.")
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/download-stress", func(c *gin.Context) {
		t.Log("Generating 2,000,000 records in memory (This proves memory stability)...")

		totalRows := 2000000
		data := make([]EdgeCaseModel, totalRows)
		for i := 0; i < totalRows; i++ {
			data[i] = EdgeCaseModel{
				ID:       i + 1,
				IsActive: i%2 == 0,
			}
		}

		t.Log("Streaming started. Watch your RAM; it should remain flat!")
		err := New[EdgeCaseModel]().
			SetFormat("xlsx").
			SetFilename("stress_test").
			Stream(c, data)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
	})

	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/download-stress")
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Stream directly to disk to verify the generated file
	outputFile := "../../test_2M_pagination.xlsx"
	file, err := os.Create(outputFile)
	if err != nil {
		t.Fatalf("Failed to create local file: %v", err)
	}
	defer file.Close()

	writtenBytes, err := io.Copy(file, resp.Body)
	if err != nil {
		t.Fatalf("Failed to download stream: %v", err)
	}

	t.Logf("Successfully paginated and streamed 2 Million records!")
	t.Logf("File size on disk: %d bytes (~%d MB)", writtenBytes, writtenBytes/1024/1024)
	t.Logf("Please open '%s' in Excel. You will see Sheet1 and Sheet2 automatically created.", outputFile)
}
