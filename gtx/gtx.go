package gtx

func Zero[T any]() (t T) {
	return
}

func IsZero[T comparable](t T) bool {
	return t == *new(T)
}
