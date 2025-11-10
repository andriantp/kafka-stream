GO_FILE    = docker/go.docker-compose.yml
FLINK_FILE = docker/flink.docker-compose.yml

network:
	@if ! docker network inspect kafka_network > /dev/null 2>&1; then \
		echo "Network kafka_network not found, creating..."; \
		docker network create --driver=bridge kafka_network; \
	else \
		echo "Network kafka_network already exists."; \
	fi

# ======================== prune ========================  
clean:system volume
	@echo "🧹 Docker cleanup completed."
volume:
	docker volume prune -a -f 
system: 
	docker system prune -a -f

# ======================== GO ========================
go-up:network
	@echo "🐳 Starting (Kafka & AKHQ) with SASL-PLAIN containers Go ..."
	chmod -R 777 docker/kafka
	docker compose -f $(GO_FILE) up --force-recreate -d --build 
	@echo "✅ (Kafka & AKHQ) with SASL-PLAIN Go are up"

go-down:
	@echo "🛑 Stopping (Kafka + AKHQ) with SASL-PLAIN containers Go ..."
	docker compose -f $(GO_FILE) down
	@echo "✅ Containers Go stopped"

# ======================== FLINK ========================
flink-up:
	@echo "🐳 Starting FLINK ..."
	docker compose -f $(FLINK_FILE) up --force-recreate -d --build 
	@echo "✅ FLINK are up"

flink-down:
	@echo "🛑 Stopping FLINK ..."
	docker compose -f $(FLINK_FILE) down
	@echo "✅ Containers FLINK stopped"

# ======================== log ========================  
ps:
	@echo "📋 Checking container status..."
	docker ps -a

kafka:
	@echo "📜 Showing Kafka logs..."
	docker logs -f kafka

akhq:
	@echo "📜 Showing AKHQ logs..."
	docker logs -f akhq

flink:
	@echo "📜 Showing Flink logs..."
	docker ps --filter "name=flink"

jobmanager:
	@echo "📜 Showing jobmanager logs..."
	docker logs -f flink-jobmanager

taskmanager:
	@echo "📜 Showing taskmanager logs..."
	docker logs -f flink-taskmanager	


