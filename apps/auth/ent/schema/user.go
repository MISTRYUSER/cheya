package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User 用户实体
type User struct {
	ent.Schema
}

// Fields 定义用户字段
func (User) Fields() []ent.Field {
	return []ent.Field{
		// 1. 用户名
		field.String("username").
			Unique().
			NotEmpty(),
		// 2. 密码
		field.String("password").NotEmpty(),
		// 3. 邮箱（可选）
		field.String("email").Optional(),
		// 4. 创建时间
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}
//Edges 
func (User) Edges() []ent.Edge {
	return nil
}

//Indexes 定义索引
func (User) Indexes() []ent.Index{
	return []ent.Index{
		index.Fields("username"),
	}
}
