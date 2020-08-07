package main

import "fmt"

//verify required car fields are filled
func verifyCarItem(c *car) error {
	if c.Name == "" {
		return fmt.Errorf("Name field required")
	}
	if c.Description == "" {
		return fmt.Errorf("Description field required")
	}
	if c.Charge.Currency == "" {
		return fmt.Errorf("Currency field in charge required")
	}
	if c.Charge.Per == "" {
		return fmt.Errorf("Per field in charge required")
	}
	if c.Make == "" {
		return fmt.Errorf("Make field required")
	}
	if c.Merchant == "" {
		return fmt.Errorf("Merchant field required")
	}

	return nil
}
