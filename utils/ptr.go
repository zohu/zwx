package utils

func Ptr[T any](arg T) *T {
	return &arg
}
