# Kafka

## Kafka CLI
### Kafka Topics:
- Để thiết lập kafka qua command line:
  ```
  kafka-topics.sh --command-config playground.config --bootstrap-server <endpoint>
  ```
- Với playground.config là properties của cluster, bootstrap-server là endpoint của cluster.
#### Để tạo một topic:
  ```
  --create --topic <name> --partitions <số_par> --replication-factor <số_replica>
  ```
- Với partitions là số phân vùng, replication-factor là số replica.
> [!NOTE]
> Số replica không vượt quá số brokers đang tồn tại.
#### Liệt kê các topics đang có
  ```
  --list
  ```
#### Mô tả thông tin của một topic:
  ```
  --topic <tên_topic> --describe 
  ```
#### Xóa một topic:
```
  --topic <tên_topic> --delete
```
### Kafka producer
- Để thiết lập kafka producer:
```
kafka-console-producer.sh --producer.config playground.config --bootstrap-server <endpoint>
```
#### Để gửi message đến topic:
```
--topic first_topic
```
#### Properties:
```
--producer-property acks=all
```
#### Key:value 
```
--topic <tên_topic> --property parse.key=true --property key.separator=<dấu ngăn cách giữa key và value>
```
### Kafka consumer
- Thiết lập consumer:
```
kafka-console-consumer.sh --consumer.config playground.config --bootstrap-server <endpoint>
```
#### Để xem tất cả các message từ khi khởi tạo:
```
--topic <tên_topic> --from-beginning
```
#### Để xem chi tiết các thông tin của các msg (key, value, timestamp):
```
--topic <tên_topic> --formatter kafka.tools.DefaultMessageFormatter --property print.timestamp=true --property print.key=true --property print.value=true --property print.partition=true --from-beginning
```
#### Sử dụng arg --group: Nhiều consumer live một lúc
```
 --topic <tên_topic> --group my-first-application
```
- Để khai thác thông tin của group:
```
kafka-consumer-groups.sh --command-config playground.config --bootstrap-server <endpoint>
```
#### List tất cả các group đã tạo:
```
 --list
```
#### Mô tả thông tin của một group:
```
--describe --group my-second-application
```
#### Để reset offset:
```
--group <tên_group> --reset-offsets --to-earliest --topic third_topic --execute
```
