package queue

// Pusher go-queue的kq
type Pusher interface {
	Name() string
	Push(string) error
}
