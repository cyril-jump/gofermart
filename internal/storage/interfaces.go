package storage

type Logger interface {
	Close()
}

type DB interface {
	Ping() error
	Close() error
}
