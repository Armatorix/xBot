package main

import (
	"fmt"
	"testing"
)

func TestFollowCount(t *testing.T) {
	total := 0
	for i := range 48 {
		count := followCount(i/2, 200)
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

func TestUnfollowCount(t *testing.T) {
	total := 0
	for i := range 48 {
		count := unfollowCount(i/2, 500, 1000) // Example values for followers and following
		total += count
		fmt.Printf("Hour: %d, Unfollow Count: %d\n", i/2, count)
		if count < 0 {
			t.Errorf("Unfollow count should not be negative, got %d for hour %d", count, i)
		}
	}

	fmt.Printf("Total unfollow count over 24 hours: %d\n", total)
	if total < 0 {
		t.Error("Total unfollow count should not be negative")
	}
	// as it is run twice a hour
	if total > 400 { // Arbitrary limit for testing purposes
		t.Error("Total unfollow count exceeds expected limit")
	}
}
