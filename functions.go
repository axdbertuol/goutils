package main

import "time"

// MapFunc is a function type for mapping elements

// Map applies the given function to each element of the slice and returns a new slice
func Map2[T, U any](data []T, f func(T) U) []U {

	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

func Filter[T any](data []T, f func(T) bool) []T {

	fltd := make([]T, 0, len(data))

	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}

// Wait function waits for the specified duration
func Wait(duration time.Duration) {
	// Create a channel to signal when the wait is done
	done := make(chan bool)

	// Start a goroutine to wait for the specified duration
	go func() {
		time.Sleep(duration)
		// Send a signal on the channel when the wait is done
		done <- true
	}()

	// Wait for the signal on the channel
	<-done
}

func RemoveElement[T any](slice []T, index int) []T {
	if index < 0 || index >= len(slice) {
		// Index out of range
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}
