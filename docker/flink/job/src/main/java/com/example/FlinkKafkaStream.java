package com.example;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.api.common.functions.ReduceFunction;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import org.apache.flink.connector.base.DeliveryGuarantee;
import org.apache.flink.connector.kafka.sink.KafkaRecordSerializationSchema;
import org.apache.flink.connector.kafka.sink.KafkaSink;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.api.java.tuple.Tuple2;

/**
 * Contoh job sederhana untuk membaca data dari Kafka, menghitung rata-rata suhu,
 * dan menulis hasilnya ke topik Kafka lain.
 */
public class FlinkKafkaStream {

    public static void main(String[] args) throws Exception {

        // === Setup environment ===
        final StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        // Aktifkan checkpoint setiap 10 detik
        // env.enableCheckpointing(10000);

        // // konfigurasi tambahan (opsional)
        // env.getCheckpointConfig().setMinPauseBetweenCheckpoints(5000); // jeda minimal antar checkpoint
        // env.getCheckpointConfig().setCheckpointTimeout(60000);         // timeout per checkpoint 60 detik
        // env.getCheckpointConfig().setMaxConcurrentCheckpoints(1);      // satu checkpoint dalam satu waktu
        // env.getCheckpointConfig().setTolerableCheckpointFailureNumber(3); // boleh gagal 3x sebelum job gagal
        
        // === Kafka Source ===
        KafkaSource<String> source = KafkaSource.<String>builder()
                .setBootstrapServers("kafka:29092")
                .setTopics("stream-go")
                .setGroupId("flink-group")
                .setValueOnlyDeserializer(new SimpleStringSchema())
                .build();

        DataStream<String> input = env.fromSource(
                source,
                WatermarkStrategy.noWatermarks(),
                "Kafka Source"
        );

        // === Transformasi sederhana: hitung rata-rata suhu per sensor ===
        // reduce tanpa windowing
        DataStream<String> averaged = input
                .map((MapFunction<String, SensorReading>) SensorReading::fromJson)
                .keyBy(r -> r.sensor)
                .reduce((ReduceFunction<SensorReading>) (r1, r2) -> {
                    double avg = (r1.temp + r2.temp) / 2.0;
                    return new SensorReading(r1.sensor, avg);
                })
                .map((MapFunction<SensorReading, String>) SensorReading::toJson);
        // reduce dengan windowing x menit
        // DataStream<String> averaged = input
        //     .map((MapFunction<String, SensorReading>) SensorReading::fromJson)
        //     .keyBy(r -> r.sensor)
        //     .window(org.apache.flink.streaming.api.windowing.assigners.TumblingProcessingTimeWindows.of(
        //             org.apache.flink.streaming.api.windowing.time.Time.minutes(5)
        //     ))
        //     .aggregate(new org.apache.flink.api.common.functions.AggregateFunction<SensorReading, Tuple2<Double, Integer>, String>() {

        //         @Override
        //         public Tuple2<Double, Integer> createAccumulator() {
        //             return new Tuple2<>(0.0, 0);
        //         }

        //         @Override
        //         public Tuple2<Double, Integer> add(SensorReading value, Tuple2<Double, Integer> accumulator) {
        //             return new Tuple2<>(accumulator.f0 + value.temp, accumulator.f1 + 1);
        //         }

        //         @Override
        //         public String getResult(Tuple2<Double, Integer> accumulator) {
        //             double avg = accumulator.f0 / accumulator.f1;
        //             return String.format("{\"sensor\":\"A1\", \"avg_temp\":%.2f}", avg);
        //         }

        //         @Override
        //         public Tuple2<Double, Integer> merge(Tuple2<Double, Integer> a, Tuple2<Double, Integer> b) {
        //             return new Tuple2<>(a.f0 + b.f0, a.f1 + b.f1);
        //         }
        //     });

        // === Kafka Sink ===
        KafkaSink<String> sink = KafkaSink.<String>builder()
                .setBootstrapServers("kafka:29092")
                .setRecordSerializer(
                        KafkaRecordSerializationSchema.builder()
                                .setTopic("stream-go-agg")
                                .setValueSerializationSchema(new SimpleStringSchema())
                                .build()
                )
                .setDeliverGuarantee(DeliveryGuarantee.AT_LEAST_ONCE)
                .build();

        averaged.sinkTo(sink);

        // === Eksekusi job ===
        env.execute("Flink Kafka Stream Average");
    }

    // === Data Model ===
    public static class SensorReading {
        public String sensor;
        public double temp;

        public SensorReading() {}

        public SensorReading(String sensor, double temp) {
            this.sensor = sensor;
            this.temp = temp;
        }

        public static SensorReading fromJson(String json) {
            try {
                json = json.trim().replace("{", "").replace("}", "");
                String[] parts = json.split(",");
                String id = parts[0].split(":")[1].replace("\"", "").trim();
                double temp = Double.parseDouble(parts[1].split(":")[1].trim());
                return new SensorReading(id, temp);
            } catch (Exception e) {
                return new SensorReading("unknown", 0.0);
            }
        }

        public String toJson() {
            return String.format("{\"sensor\":\"%s\", \"avg_temp\":%.2f}", sensor, temp);
        }
    }
}
