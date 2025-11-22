package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Vehicle struct{
	ent.Schema
}

// Fields 定义数据库字段
// 对应白皮书 6.1 章节的 vehicles 表设计
func (Vehicle) Fields() []ent.Field {
	return []ent.Field{
		// 1. VIN 码 (车架号)
		// 定义: VARCHAR(17) UNIQUE NOT NULL
		field.String("vin").
			Unique().
			NotEmpty(),

		// 2. 车牌号
		// 定义: VARCHAR(20) NOT NULL
		field.String("license_plate").
			NotEmpty(),

		// 3. 车辆状态
		// 定义: VARCHAR(20) DEFAULT 'Offline'
		field.String("status").
			Default("Offline"),

		// 4. 最后心跳时间
		// 定义: TIMESTAMP
		// 允许为空 (Optional, Nillable)，因为新车可能还没发过心跳
		field.Time("last_heartbeat").
			Optional().
			Nillable(),

		// 5. 当前位置 (JSONB)
		// 定义: location JSONB
		// 使用结构体存储，Ent 会自动转为 JSON 存入数据库
		field.JSON("location", &Location{}).
			Optional(),

		// 6. 最新遥测数据 (JSONB)
		// 定义: telemetry JSONB
		// 用于缓存该车辆最后一次上报的完整数据
		field.JSON("telemetry", map[string]interface{}{}).
			Optional(),

		// 7. 创建时间
		// 定义: TIMESTAMP DEFAULT NOW()
		field.Time("created_at").
			Default(time.Now).
			Immutable(), // 创建后不可修改

		// 8. 更新时间
		// 定义: TIMESTAMP DEFAULT NOW()
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now), // 每次更新数据自动刷新此时间
	}
}

// Edges 定义关联关系
func (Vehicle) Edges() []ent.Edge {
	return nil
}

// Indexes 定义索引
func (Vehicle) Indexes() []ent.Index {
	return []ent.Index{
		// 对应白皮书: CREATE INDEX idx_vehicles_status ON vehicles(status);
		// 优化 "查询所有在线车辆" 的速度
		index.Fields("status"),
	}
}

// Location 结构体用于配合 JSON 字段
// 这不是数据库表，只是 JSON 的数据结构
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}