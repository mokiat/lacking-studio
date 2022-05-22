package data

import "strings"

type Filter[T any] func(T) bool

func FilterAny[T any]() Filter[T] {
	return func(T) bool {
		return true
	}
}

func FilterNone[T any]() Filter[T] {
	return func(T) bool {
		return false
	}
}

func FilterAnd[T any](filters ...Filter[T]) Filter[T] {
	return func(item T) bool {
		for _, filter := range filters {
			if !filter(item) {
				return false
			}
		}
		return true
	}
}

func FilterOr[T any](filters ...Filter[T]) Filter[T] {
	if len(filters) == 0 {
		return FilterAny[T]()
	}
	return func(item T) bool {
		for _, filter := range filters {
			if filter(item) {
				return true
			}
		}
		return false
	}
}

func FilterWithKind(kind ResourceKind) Filter[*Resource] {
	return func(resource *Resource) bool {
		return resource.kind == kind
	}
}

func FilterWithSimilarName(name string) Filter[*Resource] {
	if name == "" {
		return FilterAny[*Resource]()
	}
	name = strings.ToLower(name)
	return func(resource *Resource) bool {
		return strings.Contains(strings.ToLower(resource.name), name)
	}
}
