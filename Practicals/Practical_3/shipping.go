// shipping.go
package shipping

import (
	"errors"
	"fmt"
)

// CalculateShippingFee calculates the fee based on weight and zone.
func CalculateShippingFee(weight float64, zone string) (float64, error) {
	// Rule #1 and #4: Check weight validity
	if weight <= 0 || weight > 50 {
		return 0, errors.New("invalid weight: must be between 0 and 50 kg")
	}

	// Rule #2 and #5: Check zone validity
	switch zone {
	case "Domestic":
		return 5 + (weight * 1.0), nil
	case "Express":
		return 30 + (weight * 5.0), nil
	case "International":
		return 0, errors.New("international shipping not implemented")
	default:
		return 0, fmt.Errorf("invalid zone: %s", zone)
	}
}
