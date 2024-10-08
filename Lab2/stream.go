package main

import (
	"fmt"
)

type Stream[T Displayable] struct {
	elements []T
}

func CreateStream[T Displayable](elements []T) Stream[T] {
	return Stream[T]{elements}
}

func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	var filtered []T
	for _, e := range s.elements {
		if predicate(e) {
			filtered = append(filtered, e)
		}
	}
	return Stream[T]{filtered}
}

func (s Stream[T]) Map(transform func(T) T) Stream[T] {
	var transformed []T
	for _, e := range s.elements {
		transformed = append(transformed, transform(e))
	}
	return Stream[T]{transformed}
}

func (s Stream[T]) Max(less func(T, T) bool) *T {
	if len(s.elements) == 0 {
		return nil
	}
	max := s.elements[0]
	for _, e := range s.elements[1:] {
		if less(max, e) {
			max = e
		}
	}
	return &max
}

func (s Stream[T]) Reduce(initialValue int, accumulator func(int, T) int) int {
	result := initialValue
	for _, e := range s.elements {
		result = accumulator(result, e)
	}
	return result
}

func (s Stream[T]) Distinct() Stream[T] {
	seen := make(map[string]struct{})
	var distinct []T
	for _, e := range s.elements {
		displayStr := e.display()
		if _, exists := seen[displayStr]; !exists {
			seen[displayStr] = struct{}{}
			distinct = append(distinct, e)
		}
	}
	return Stream[T]{distinct}
}

func (s Stream[T]) Display() {
	for _, e := range s.elements {
		fmt.Println(e.display())
	}
}