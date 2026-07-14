package shared

type Response[T any] struct {
	Data T
}

type Empty struct{}
