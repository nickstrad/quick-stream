KAFKA_DIR=bin/kafka_2.11-0.10.0.1

PHONY: api/start
api/start:
	go run ./cmd

PHONY: test
test:
	go test ./...

PHONY: test-verbose
test-verbose:
	go test -v ./...

PHONY: testdox
testdox:
	cd cmd/hybrid_message_logger && gotestdox


PHONY: service/start/kafka
service/start/kafka:
	kafka-server-start /opt/homebrew/etc/kafka/server.properties

PHONY: service/start/zookeeper
service/start/zookeeper:
	zookeeper-server-start /opt/homebrew/etc/kafka/zookeeper.properties

PHONY: service/create-topic/kafka
service/create-topic/kafka:
	/opt/homebrew/opt/kafka/bin/kafka-topics \
	--bootstrap-server localhost:9092 \
	--create \
	--topic meetup-raw-rsvps \
	--partitions 1 \
	--replication-factor 1
