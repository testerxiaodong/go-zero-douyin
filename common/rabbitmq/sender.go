package rabbitmq

type Sender interface {
	Send(exchange string, routeKey string, msg []byte) error
}
