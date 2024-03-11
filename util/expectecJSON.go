package util

import "fmt"

func CheckExpectedFields(data map[string]string, expectedFields []string) error {
	for key := range data {
		found := false
		for _, field := range expectedFields {
			if key == field {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("Unexpected field '%s' in request body", key)
		}
	}
	return nil
}
