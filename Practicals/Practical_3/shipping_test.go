package shipping

import (
	"testing"
)

func TestCalculateShippingFee_EquivalencePartitioning(t *testing.T) {
	// Test cases for weight partitions
	t.Run("Invalid Weight (Negative)", func(t *testing.T) {
		_, err := CalculateShippingFee(-5.0, "Domestic")
		if err == nil {
			t.Error("Expected error for negative weight, got nil")
		}
	})

	t.Run("Valid Weight", func(t *testing.T) {
		fee, err := CalculateShippingFee(25.0, "Domestic")
		if err != nil {
			t.Errorf("Expected no error for valid weight, got: %v", err)
		}
		expectedFee := 30.0 // 5 base + (25 * 1.0)
		if fee != expectedFee {
			t.Errorf("Expected fee %.2f, got %.2f", expectedFee, fee)
		}
	})

	t.Run("Invalid Weight (Too Heavy)", func(t *testing.T) {
		_, err := CalculateShippingFee(51.0, "Domestic")
		if err == nil {
			t.Error("Expected error for weight > 50kg, got nil")
		}
	})

	// Test cases for zone partitions
	t.Run("Valid Zone (Domestic)", func(t *testing.T) {
		_, err := CalculateShippingFee(10.0, "Domestic")
		if err != nil {
			t.Errorf("Expected no error for Domestic zone, got: %v", err)
		}
	})

	t.Run("Invalid Zone", func(t *testing.T) {
		_, err := CalculateShippingFee(10.0, "InvalidZone")
		if err == nil {
			t.Error("Expected error for invalid zone, got nil")
		}
	})
}

func TestCalculateShippingFee_BoundaryValueAnalysis(t *testing.T) {
	tests := []struct {
		name        string
		weight      float64
		zone        string
		wantFee     float64
		wantErr     bool
		description string
	}{
		{
			name:        "Weight just below 0",
			weight:      -0.01,
			zone:        "Domestic",
			wantErr:     true,
			description: "Weight slightly below minimum",
		},
		{
			name:        "Weight exactly 0",
			weight:      0.0,
			zone:        "Domestic",
			wantErr:     true,
			description: "Weight at minimum boundary",
		},
		{
			name:        "Weight just above 0",
			weight:      0.01,
			zone:        "Domestic",
			wantFee:     5.01, // 5 base + (0.01 * 1.0)
			description: "Weight just inside valid range",
		},
		{
			name:        "Weight exactly 50",
			weight:      50.0,
			zone:        "Domestic",
			wantFee:     55.0, // 5 base + (50 * 1.0)
			description: "Weight at maximum boundary",
		},
		{
			name:        "Weight just above 50",
			weight:      50.01,
			zone:        "Domestic",
			wantErr:     true,
			description: "Weight slightly above maximum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateShippingFee(tt.weight, tt.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: CalculateShippingFee() error = %v, wantErr %v", tt.description, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantFee {
				t.Errorf("%s: CalculateShippingFee() = %.2f, want %.2f", tt.description, got, tt.wantFee)
			}
		})
	}
}

func TestCalculateShippingFee_DecisionTable(t *testing.T) {
	tests := []struct {
		name        string
		weight      float64
		zone        string
		wantFee     float64
		wantErr     bool
		description string
	}{
		{
			name:        "Rule 1: Invalid weight",
			weight:      -1.0,
			zone:        "Domestic",
			wantErr:     true,
			description: "Invalid weight should return error regardless of zone",
		},
		{
			name:        "Rule 2: Valid weight, Domestic zone",
			weight:      10.0,
			zone:        "Domestic",
			wantFee:     15.0, // 5 base + (10 * 1.0)
			description: "Valid domestic shipping calculation",
		},
		{
			name:        "Rule 3: Valid weight, International zone",
			weight:      10.0,
			zone:        "International",
			wantErr:     true,
			description: "International shipping should return error (not implemented)",
		},
		{
			name:        "Rule 4: Valid weight, Express zone",
			weight:      10.0,
			zone:        "Express",
			wantFee:     80.0, // 30 base + (10 * 5.0)
			description: "Valid express shipping calculation",
		},
		{
			name:        "Rule 5: Valid weight, Invalid zone",
			weight:      10.0,
			zone:        "Invalid",
			wantErr:     true,
			description: "Invalid zone should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateShippingFee(tt.weight, tt.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: CalculateShippingFee() error = %v, wantErr %v", tt.description, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantFee {
				t.Errorf("%s: CalculateShippingFee() = %.2f, want %.2f", tt.description, got, tt.wantFee)
			}
		})
	}
}
