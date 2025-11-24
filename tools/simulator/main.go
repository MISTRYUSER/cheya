package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

// TelemetryData 结构体需与消费者一致
type TelemetryData struct {
	VehicleID string  `json:"vehicle_id"`
	Timestamp int64   `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
}

func main() {
	//配置 kafka producer
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "telemetry.raw",
		Balancer: &kafka.LeastBytes{},
	}
	defer w.Close()
	//2.redis client
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	//VIN
	vehicleID := "VIN-TEST-SIM-01"

	//控制标志位
	isRunning := true
	wasStopped := false // 用于跟踪是否已经打印过停止日志

	//启动指令监听协程
	go func() {
		log.Println("👂 Listening for commands on Redis channel: vehicle:commands")
		sub := rdb.Subscribe(context.Background(), "vehicle:commands")
		ch := sub.Channel()

		for msg := range ch {
			if msg.Payload == "STOP:"+vehicleID {
				log.Println("🛑 收到远程停车指令！！！")
				isRunning = false
			} else if msg.Payload == "START:"+vehicleID {
				log.Println("▶️ 收到远程启动指令")
				isRunning = true
				wasStopped = false // 重置标志
			}
		}
	}()
	//起始位置 东方明珠
	lat := 31.2397
	lon := 121.4998

	log.Printf("🚀 Simulator started for vehicle: %s", vehicleID)

	for {
		if !isRunning {
			if !wasStopped {
				log.Println("⏸️  车辆已停止，等待恢复指令...")
				wasStopped = true
			}
			time.Sleep(1 * time.Second)
			continue
		}
		
		// 恢复运行时打印日志
		if wasStopped {
			log.Println("✅ 车辆已恢复运行")
			wasStopped = false
		}
		//1.模拟移动
		lat += (rand.Float64() - 0.5) * 0.001
		lon += (rand.Float64() - 0.5) * 0.001
		speed := 40.0 + (rand.Float64() * 40.0)

		//2.组装数据
		data := TelemetryData{
			VehicleID: vehicleID,
			Timestamp: time.Now().Unix(),
			Latitude:  lat,
			Longitude: lon,
			Speed:     speed,
		}
		jsonData, _ := json.Marshal(data)

		//3.把数据发送到 kafka
		err := w.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(vehicleID), //保证同一辆车有序
				Value: jsonData,
			},
		)

		if err != nil {
			log.Printf("❌ Failed to write messages: %v", err)
		} else {
			fmt.Printf("📤 Sent: %s\n", string(jsonData))
		}

		time.Sleep(1 * time.Second)
	}
}
