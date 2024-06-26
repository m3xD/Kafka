# Kafka

## So sánh với RabbitMQ và Redis
| Kafka         | RabbitMQ | Redis | 
| ------------- | ------------- | --------- |
| Có thể gửi up to 1 triệu msg/s  | Có thể gửi 50000 msg/s  | Có thể gửi up to 1 triệu msg/s |
| Có cơ chế khôi phục dữ liệu và tính sẵn dùng  | Tương tự như kafka  | Không có cơ chế bảo vệ dữ liệu |
| One to many | One to one, One to many | One to one, One to many |
| High-throughput | Đa dạng về msg pattern | Độ trễ thấp, hiệu năng cao |

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
#### Kafka Message:
![Elements of Kafka Message](https://www.conduktor.io/kafka/_next/image/?url=https%3A%2F%2Fimages.ctfassets.net%2Fo12xgu4mepom%2F2TuJ55uK20OUVLQgZ17yUU%2F9bb611597f4914e971d85e3938856968%2FKafka_Producers_3.png&w=1920&q=75)
- Key: optional, key sẽ được serialize thành binary format.
- Value: có thể null, là nội dụng của msg, value sẽ được serialize thành binary format.
- Compression type: optional, msg có thể nén thành các dạng: gzip, lz4,..
- Headers: là các cặp key:value, thường được sử dụng để mô tả metadata về msg.
- Partition + offset: một msg được định danh bởi số thứ tự phân vùng và offset id.
- Timestamp: có thể thêm bởi người dùng hoặc system.
#### Kafka Message Serializers
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
#### Kafka Consumer Offsets
- Nhằm xác định được consumer đã đọc msg tới đâu tại hiện thời, consumer thường commit offset đã đọc được.
- Nếu có lỗi xảy ra, consumer sẽ biết bắt đầu từ đâu.
- Nếu có một consumer mới được thêm, thông tin offset cũng được sử dụng để bắt đầu đọc.
### Kafka brokers
- Mỗi phân vùng lại nằm trên các server khác nhau, cũng được biết đến với tên khác là broker.
- Một kafka server thì được gọi là kafka broker
#### Kafka cluster
- Một tập hợp các brokers thì được gọi là một kafka cluster.
- Các client muốn nhận hay gửi msg thì có thể kết nối đến bất kì broker nào. Các broker đều có metadata của các broker khác .
#### Kafka brokers và Topics
- Để đạt được thông lượng và mức độ mở rộng thì các phân vùng được phân bố đều trên các brokers.
#### Cách client kết nối tới Kafka Cluster
- Broker trong cluster có tên gọi khác là bootstrap server.
- Khi client kết nối đến 1 bootstrap server, nó sẽ trả về list các broker có trong cluter.
### Kafka Topic Replication:
- Phòng tránh thất lạc dữ liệu khi broker bị hỏng bằng cách replica các phân vùng.
- Mỗi phân vùng thì có một phân vùng là leaader và các phân vùng còn lại là replica
####  In-Sync Replicas:
- Các replica được cập nhật liên tục cùng với phân vùng leader thì được gọi là ISR.
#### ACK:
- Cung cấp giải pháp ghi tối thiểu vào replica để được coi là ghi thành công từ producer.
- ACK config:
  * ack = 0: msg được gửi đi mà không cần phản hồi lại.
  * ack = 1: msg được gửi đi nhưng cần leader partition phản hồi là ghi thành công mới tính là thành công.
  * ack = all: msg được gửi leader và ISR ghi thành công và phản hồi lại, số ISR tối thiểu cần đáp ứng được config qua **min.insync.replicas**.
#### Kafka Topic Durability & Availability: 
```
Nếu có N replica, thì ta có thể mất N-1 brokers mà vẫn có thể khôi phục dữ liệu.
```
#### Kafka Consumers Replicas Fetching
- Consumer có thể đọc trực tiếp từ ISR mà không cần qua leader partition để tăng hiệu năng và giảm chi phí.

## Advanced Kafka
### Producers:
#### Acks:
- Số brokers phải respone lại thì được coi là ghi thành công.
- Ack = 0:
  * Chỉ request lên broker mà không cần chờ ack.
  * Có thể xảy ra mất mát dữ liệu.
  * Không thể biết được request đã tới nơi chưa.
- Ack = 1:
  * Producer gửi request lên leader broker và chờ ack chỉ từ leader broker.
  * Nếu không nhận được ack, producer sẽ thử gửi lại request.
  * Có thể xảy ra mất mát dữ liệu nếu leader broker offline và replica chưa kịp được backup.
- Ack = all:
  * Producer gửi request lên, bao gồm cả leader broker và các ISR phải phản hồi lại ack để được tính là ghi thành công.
  * Leader broker sẽ kiểm tra số lượng ISR cần thiết để ghi data, chờ đợi ack từ các ISR, nếu các ISR tối thiểu đã phản hồi thì sucessful ack sẽ được gửi lại producer.
 - Tính bền vững và sẵn sàng:
   * Nếu có N replication factor, broker có thể sập đến N - 1 mà vẫn có thể khôi phục dữ liệu.
   * Nếu ack = all, replication factor = N, min.insync.replicas = M, broker có thể offline tới N - M mà vẫn có thể đảm bảo tính sẵn sàng.
 - Safety settings (default in kafka ver >= 3.0): acks = all và min.insync.replicas = 2 là cài đặt thông dụng nhất, broker có thể unavailable nhiều nhất một để đảm bảo tính sẵn sàng và bền vững.

#### Retries:
- Gửi lại request khi broker gửi error code, và các trường hợp thử lại là các "Retriable error", có thể khắc phục khi thử lại.
- Retries: số lần attempt khi broker gửi error code trước khi được mark là fail.
- Delivery.timeout.ms:
  * Record được mark là fail nếu không thể gửi tin trong thời gian **Delivery.timeout.ms** được config
  * `delivery.timeout.ms >= (linger.ms + retry.backoffms + request.timeout.ms)`
 - retry.backoff.ms: producer sẽ đợi trước khi thử lại sau thời gian này.
 - max.in.flight.request.per.connection:
   * số lượng record gửi đi trong 1 connection
   * nếu config > 1, order của record có thể không được đảm bảo.
   * nếu config = 1, throughput có thể suy giảm nhiều.
> [!NOTE]
> nếu **enable=idempotence=true**, thì **max.in.flight.requests.per.connection** <= 5 để đảm bảo order của record.

#### Idempotent:
- Có một khả năng nhỏ khi cả 2 record đều được ghi thành công broker, dẫn tới duplicate.
![Duplicate msg](https://www.conduktor.io/kafka/_next/image/?url=https%3A%2F%2Fimages.ctfassets.net%2Fo12xgu4mepom%2F1oXO3anfY5Bm3Uvz5xZRRZ%2F379b11bbbba199b3b001bba6f3e96493%2FAdv_Idempotent_Producer_1.png&w=1920&q=75)
- Producer idempotence đảm bảo các record trùng nhau sẽ không được commit vào broker.
![Non-dup mgs](https://www.conduktor.io/kafka/_next/image/?url=https%3A%2F%2Fimages.ctfassets.net%2Fo12xgu4mepom%2F4XZYhrPGmkGENxjfWfDbrW%2F8157319ea1dbc0c1952d51f591e30b0a%2FAdv_Idempotent_Producer_2.png&w=1920&q=75)
- Cách hoạt động:
  * Mỗi producer được gán **Producer ID (PID)** mỗi lần producer gửi tin đến broker.
  * Mỗi msg được gán thêm một **sequences number**, mỗi lần tin được gửi đi, broker sẽ kiểm tra PID-Seq_No xem có lớn hơn PID-Seq_No hiện tại không, nếu không thì discard msg, nếu có thì commit và ghi nhận PID-Seq_No
  * Beside, này nghe giống mmt thế nhể =)) cur_sequences_num += size(msg).
  
 #### Message Compression:
 - Kafka hỗ trợ 2 kiểu nén dữ liệu ở: producer và broker.
 - Producer-side:
   * Có thể config thông qua tham số **compression.type**
   * Có các lựa chọn nén là: none, gzip, lz4, snappy, zstd.
 - Broker-side:
   * Nếu **compression.type=producer** được thiết lập, broker sẽ trực tiếp ghi lên topic mà không nén thêm.
   * Nếu topic có config khác trên, data sẽ được giải nén và nén theo thiết lập hiện tại.
 - Producer-level msg compression:
   * Producer sẽ group các msg thành các batch, các batch này có thể nén lại nếu được thiết lập và gửi tới kafka.
   * Batch có kích thước bé, batch cũng được transfer nhanh hơn.
>[!NOTE]
>snappy hay lz4 có sự cân bằng giữa speed và compression ratio.

#### Producer Batching:
- Kafka sẽ cố gắng để gửi các record một cách liên tục, người dùng có thể config số lượng msg batch thông qua một tham số.
- Để config batching thường có 2 tham số tham gia là **linger.ms** và **batch.size**
- Linger.ms:
  * Là thời gian producer chờ đợi msg trước khi gửi batch đi.
  * Tăng thời gian chờ đợi để tăng khả năng các msg được gửi đi cùng một batch.
  * Nếu batch đã đầy trước khi hết thời gian thì batch cũng được gửi đi ngay.
- Batch.size:
  * Là size tối đa mà batch có thể chứa các msg.
- Tất cả các tham số trên đều có thể giúp tăng throughput cũng như gửi ít request hơn các msg được đóng thành từng batch.

#### Round robin partitioner and sticky partitioner:
- Round robin:
  * Các data lần lượt được sặp xếp vào các partition theo qui tắc lần lượt từng partition nhận một data cho đến partition thì quay vòng về partition đầu tiên.
  * Có độ trễ cao
- Sticky:
  * Các data được lấp đầy 1 batch hoặc linger.ms timeout thì mới tới batch kế tiếp.
  * Có độ trễ thấp hơn
 
### Consumer:
#### Delivery Semantics:
- At most once:
  * Offset được commit ngay khi poll được gọi, nếu có xảy ra lỗi, offset không thể đọc lại do đã commit.
  * Phù hợp với các hệ thống có thể chấp nhận mất mát dữ liệu.
 ![at_most_once](https://www.conduktor.io/kafka/_next/image/?url=https%3A%2F%2Fimages.ctfassets.net%2Fo12xgu4mepom%2F3HWZKKcNQLUiVaZBHOQKWt%2F422c797f8ee23205502b6ddc6e321e54%2FAdv_Delivery_Semantics_for_Consumers_1_2x.png&w=3840&q=75)
- At least once:
  * Các offset được commit ngay khi được process bởi consumer.
  * Có thể có các offset được đọc lại nhiều lần nếu có lỗi xảy ra.
  * Phù hợp với các hệ thống không chấp nhận mất mát dữ liệu.
 ![at_least_once](https://www.conduktor.io/kafka/_next/image/?url=https%3A%2F%2Fimages.ctfassets.net%2Fo12xgu4mepom%2F3X2ZPpl4wkTmvnly4e60cu%2F6d3b7404b93030bdfff9012f465ac560%2FAdv_Delivery_Semantics_for_Consumers_2_2x.png&w=3840&q=75)
- Ta có thể tự động commit sau khoảng thời gian nhất định thông qua 2 tham số là **enable.auto.commit** và **auto.commit.interval.ms** sau khi gọi hành động poll.

#### Auto Offsets Reset:
- auto.offset.reset:
  * lastest: đọc dữ liệu từ offset mới nhất phân vùng.
  * earliest: đọc từ offset cũ nhất trong phân vùng.
  * none: throw exception nếu không có offset trước đó.
- offset.retention.minutes: giúp kéo dài thời gian trước khi bị reset offset.


#### Consumer Read from Closest Replica
- Mặc định consumer sẽ lấy dữ liệu từ leader partition.
- Từ những version mới hơn của kafka (>=2.4), consumer có thể lấy data từ các ISR gần nhất giúp giảm độ trễ, chi phí.

#### Incremental Rebalance & Static Group Membership:
- Khi consumer vào group, ra khỏi group hay một phân vùng được thêm vào topic thì cần được rebalance group.
- Eager Rebalancing:
  * khi xảy ra các hiện tưởng kể trên, các consumers sẽ tạm dừng lấy data từ kafka và tạm thời dừng liên kết với các partition hiện tại.
  * các consumer khi được rebalance thì có thể không còn kết nối với các partition trước đấy.
- Cooperative Rebalance:
  * Chỉ có một tập nhỏ partition sẽ ngưng liên kết và gán cho consumer mới
- Static Group Membership:
  * Khi consumer rời group và rejoin lại, nó sẽ được cấp một member_id mới và có thể được gán với phân vùng mới.
  * Nếu được config thì consumer có thể được liên kết lại đến phân vùng cũ mà không cần hành động rebalance.

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
