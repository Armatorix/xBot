package main

import (
	"fmt"
	"testing"
)

func TestFollowCount(t *testing.T) {
	total := 0
	for i := 0; i < 48; i++ {
		count := followCount(i / 2)
		total += count
		fmt.Printf("Hour: %d, Follow Count: %d\n", i/2, count)
		if count < 0 {
			t.Errorf("Follow count should not be negative, got %d for hour %d", count, i)
		}
	}

	fmt.Printf("Total follow count over 24 hours: %d\n", total)
	if total < 0 {
		t.Error("Total follow count should not be negative")
	}
	// as it is run twice a hour
	if total > 400 { // Arbitrary limit for testing purposes
		t.Error("Total follow count exceeds expected limit")
	}
}
