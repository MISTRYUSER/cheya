package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

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

	vehicleID := "VIN-TEST-SIM-01"
	//起始位置 东方明珠
	lat := 31.2397
	lon := 121.4998

	log.Printf("🚀 Simulator started for vehicle: %s", vehicleID)

	for {
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
