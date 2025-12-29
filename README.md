# Kafka Stream

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache_2.0-green.svg)](LICENSE)
[![Medium](https://img.shields.io/badge/Read-Medium-black?logo=medium)](https://andriantriputra.medium.com/golang-x-kafka-9-stream-processing-206d70e0af19)

This project is a simple example of **stream processing** using **Golang** and **Apache Kafka**.  
The application reads messages from one topic, processes them, and writes the results to another topic.

---

## Features

- Reads messages from a Kafka topic  
- Processes data (e.g., text transformation, filtering, simple aggregation)  
- Sends processed results to another Kafka topic (`output-topic`)  
- Simple configuration via environment variables  
- Demonstrates a basic **consumer–processor–producer** implementation in Go  

---

## Project Structure
```bash
tree
.
├── main.go # Application entry point
├── stream/
│ ├── consumer.go  # Module for consuming data from Kafka
│ ├── producer.go  # Module for producing data to Kafka
│ ├── processor.go # Core stream processing logic
│ └── model.go     # Data model definitions (optional)
├── go.mod
└── go.sum
```

## 🔗 Referensi

Artikel lengkap:  
[Golang x Kafka #9 — Stream Processing](https://andriantriputra.medium.com/golang-x-kafka-9-stream-processing-206d70e0af19)

---

## Author

Andrian Tri Putra
- [Medium](https://andriantriputra.medium.com/)
- [andriantp](https://github.com/andriantp)
- [AndrianTriPutra](https://github.com/AndrianTriPutra)

---

## License
Licensed under the Apache License 2.0

