# Kafka

## Khái niệm:
### Kafka Topics:
- Giống với SQL, topic giống với table nhưng không thể query
- Thay vì thế sẽ có producers gửi data đến topic và consumers sẽ lấy data theo thứ tự
- Kafka có thể send data dưới nhiều định dạng như JSON,...
- Khi đã khi data vào topic thì không thể thay đổi 
### Kafka Partitions:
- Topics thì được chia thành các phân vùng (partitions). Một topic có thể có nhiều hơn 1 phân vùng
- Số lượng phân vùng thường thấy của 1 topic: 100
- Có thể config được số lượng phân vùng của 1 topic
- Mỗi message thì được gán với 1 offset duy nhất
### Kafka Offsets
- Biểu thị vị trí của từng message trong Kafka Partition.
- Bắt đầu từ 0
- Chỉ mang ý nghĩa của từng phân vùng.
> [!NOTE]
> Message trong từng phân vùng thì có thể đảm bảo về thứ tự, điều đó không xảy ra nếu xét hơn MỘT phân vùng
- Offset không thể tái sử dụng
### Kafka Producers
- Các ứng dụng mà gửi data đến cho topic thì được gọi là Producer.
- Được viết trên nhiều ngôn ngữ: Go, Java, Python,...
- Khi data được gửi đến topic, message được phân phối cho các phân vùng dựa trên một thuật toán nào đó.
- Thuật toán murmur2 là một cách để phân phối msg trên các phân vùng:
 ```
 targetPartition = Math.abs(Utils.murmur2(keyBytes)) % (numPartitions - 1)
 ```
- Message Keys:
  * Key = null: Các message được phân bổ đồng đều trên các phân vùng trong topic theo thuật toán round-robin.
  * Key != null: Các message có cùng key xếp chung một phân vùng
### Kafka Message:
![Elements of Kafka Message](https://www.conduktor.io/kafka/_next/image/?url=https%3A%2F%2Fimages.ctfassets.net%2Fo12xgu4mepom%2F2TuJ55uK20OUVLQgZ17yUU%2F9bb611597f4914e971d85e3938856968%2FKafka_Producers_3.png&w=1920&q=75)
- Key: optional, key sẽ được serialize thành binary format.
- Value: có thể null, là nội dụng của msg, value sẽ được serialize thành binary format.
- Compression type: optional, msg có thể nén thành các dạng: gzip, lz4,..
- Headers: là các cặp key:value, thường được sử dụng để mô tả metadata về msg.
- Partition + offset: một msg được định danh bởi số thứ tự phân vùng và offset id.
- Timestamp: có thể thêm bởi người dùng hoặc system.
### Kafka Message Serializers
- Vì broker chỉ nhận dữ liệu byte nên cần có quá trình serialize
- Consumer lấy data sẽ một lần nữa deserialize
- Trong quá trình serialize và deserialize thì kiểu dữ liệu không nên thay đổi, nếu có thì cách tốt nhất là tạo một topic mới.
### Kafka Consumers:
- Các ứng dụng lấy dữ liệu từ topic thì được gọi là Consumer.
- Các consumer đọc từ offset thấp đến cao, và không để đọc ngược lại.
### Kafka Consumer Group
- Group_id giúp định dạnh consumer group
- Consumer là một phần của ứng dụng và thực hiện cùng một logic thì có thể xếp thành một nhóm Consumers.
- Các consumers có trong group có thể phân chia công việc đọc các phân vùng.
- Sử dụng **GroupCoordinator** và **ConsumerCoordinator** để gán consumer với partition và đảm bảo load balacing với tất cả consumer.
- Với mỗi phân vùng thì chỉ được gán với 1 consumer, nhưng 1 consumer có thể gán với nhiều phân vùng.
- Có thể có nhiều consumer group cùng truy cập vào 1 topic tại cùng 1 thời điểm.
- Nếu có nhiều consumer hơn phân vùng thì các consumer bị dư sẽ ở trạng thái inactive.
### Kafka Consumer Offsets
- Nhằm xác định được consumer đã đọc msg tới đâu tại hiện thời, consumer thường commit offset đã đọc được.
- Nếu có lỗi xảy ra, consumer sẽ biết bắt đầu từ đâu.
- Nếu có một consumer mới được thêm, thông tin offset cũng được sử dụng để bắt đầu đọc.

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
### Kafka Producer
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
### Kafka Consumer
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
