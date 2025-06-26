package instance

import "errors"

var (
	ErrCannotParseConfig    = errors.New("unable to get messaging config")
	ErrWrongPostgresConn    = errors.New("unable to init Postgres connection")
	ErrWrongMongoConn       = errors.New("unable to init MongoDB connection")
	ErrWrongMongoCollection = errors.New("unable to get MongoDB collection")
	ErrWrongRedisConn       = errors.New("unable to init Redis connection")
	ErrWrongProducerConn    = errors.New("unable to init Kafka producer")
	ErrWrongProducerTopic   = errors.New("unable to init Kafka producer for topic")
	ErrWrongConsumerConn    = errors.New("unable to init Kafka consumer")
	ErrWrongConsumerTopic   = errors.New("unable to init Kafka consumer for topic")
)
