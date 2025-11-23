package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

// 远程测量/远程监控数据
type TelemetryData struct {
	VehicleID string  `json:"vehicle_id"`
	Timestamp int64   `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
}

// StartTelemetryConsumer 启动消费者
func StartTelemetryConsumer(ctx context.Context, brokers []string, topic string, rdb *redis.Client) {
	//1.配置 Reader （消费者）
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  "telemetry-service-group",
		MinBytes: 1,                      // 最小 1 字节就返回，实时性更好
		MaxBytes: 10e6,                   // 最大 10MB
		MaxWait:  100 * time.Millisecond, // 最多等待 100ms
	})

	
	
	//函数结束后关闭
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("❌ Failed to close Kafka reader: %v", err)
		}
	}()

	log.Printf("🎧 Listening on Kafka topic: %s ...", topic)

	//2.循环读取消息
	for {
		//检查上下文是否取消
		select {
		case <-ctx.Done():
			log.Println("🛑 Stopping Kafka consumer...")
			return
		default:
		}

		//阻塞读取一条消息
		m, err := r.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("⚠️ Read error: %v", err)
			continue
		}

		//3.处理消息（反序列化）
		var data TelemetryData
		if err := json.Unmarshal(m.Value, &data); err != nil {
			log.Printf("⚠️ JSON parse error: %v", err)
			continue
		}

		// 发布到 Redis
		if err := rdb.Publish(ctx, "vehicle:update", m.Value).Err(); err != nil {
			log.Printf("⚠️ Redis Publish Error: %v", err)
		} else {
			log.Printf("✅ Published to Redis")
		}

		// 打印接收到的数据
		log.Printf("🚛 [Recv] Truck=%s Lat=%.6f Lon=%.6f Speed=%.1f km/h Time=%s",
			data.VehicleID,
			data.Latitude,
			data.Longitude,
			data.Speed,
			time.Unix(data.Timestamp, 0).Format("15:04:05"),
		)
		// TODO: 将数据存储到数据库或进行其他处理
	}
}
