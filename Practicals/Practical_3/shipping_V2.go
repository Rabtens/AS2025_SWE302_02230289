// shipping_V2.go
package shipping

import (
	"errors"
	"fmt"
	"math"
)

// ShippingCalculator represents a shipping fee calculator with configurable rates
type ShippingCalculator struct {
	domesticBaseRate     float64
	domesticPerKgRate   float64
	expressBaseRate     float64
	expressPerKgRate    float64
	internationalRates  map[string]float64
	volumeDiscounts     map[string]float64
	insuranceThreshold float64
	insuranceRate      float64
}

// NewShippingCalculator creates a new instance of ShippingCalculator with default rates
func NewShippingCalculator() *ShippingCalculator {
	return &ShippingCalculator{
		domesticBaseRate:     5.0,
		domesticPerKgRate:   1.0,
		expressBaseRate:     30.0,
		expressPerKgRate:    5.0,
		internationalRates:  make(map[string]float64),
		volumeDiscounts:     map[string]float64{"SUMMER10": 0.10, "BULK20": 0.20},
		insuranceThreshold: 20.0,
		insuranceRate:      0.05,
	}
}

// CalculateShippingFeeV2 calculates shipping fee with support for discounts and insurance
func (sc *ShippingCalculator) CalculateShippingFeeV2(weight float64, zone string, discountCode string) (float64, error) {
	// Validate weight
	if weight <= 0 || weight > 50 {
		return 0, errors.New("invalid weight: must be between 0 and 50 kg")
	}

	// Calculate base fee
	var baseFee float64
	switch zone {
	case "Domestic":
		baseFee = sc.domesticBaseRate + (weight * sc.domesticPerKgRate)
	case "Express":
		baseFee = sc.expressBaseRate + (weight * sc.expressPerKgRate)
	case "International":
		if rate, exists := sc.internationalRates[zone]; exists {
			baseFee = rate * weight
		} else {
			return 0, fmt.Errorf("shipping rate not available for zone: %s", zone)
		}
	default:
		return 0, fmt.Errorf("invalid zone: %s", zone)
	}

	// Apply insurance if weight exceeds threshold
	if weight > sc.insuranceThreshold {
		insuranceFee := baseFee * sc.insuranceRate
		baseFee += insuranceFee
	}

	// Apply discount if valid
	if discount, exists := sc.volumeDiscounts[discountCode]; exists {
		baseFee = baseFee * (1 - discount)
	} else if discountCode != "" {
		return 0, fmt.Errorf("invalid discount code: %s", discountCode)
	}

	// Round to 2 decimal places
	return math.Round(baseFee*100) / 100, nil
}

// SetInternationalRate sets the rate for a specific international zone
func (sc *ShippingCalculator) SetInternationalRate(zone string, ratePerKg float64) {
	sc.internationalRates[zone] = ratePerKg
}

// AddDiscountCode adds a new discount code with specified percentage
func (sc *ShippingCalculator) AddDiscountCode(code string, discountPercentage float64) error {
	if discountPercentage <= 0 || discountPercentage >= 1 {
		return errors.New("discount percentage must be between 0 and 1")
	}
	sc.volumeDiscounts[code] = discountPercentage
	return nil
}
