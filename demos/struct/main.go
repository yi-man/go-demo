package main

import (
	"fmt"
	"reflect"
	// "time"
)

type Greeter struct {
	Hello string
}

type Model struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true;comment:自增id" json:"id"` // 自增id
	// Name        string    `gorm:"column:name;not null;comment:消息类型名称" json:"name"`                                     // 消息类型名称
	// Description *string   `gorm:"column:description;comment:描述" json:"description"`                                    // 描述
	// ClientID    string    `gorm:"column:client_id;not null;comment:服务端标识" json:"client_id"`                            // 服务端标识
	// Type        string    `gorm:"column:type;not null;comment:消息类型" json:"type"`                                       // 消息类型
	// IsDeleted   bool      `gorm:"column:is_deleted;not null;comment:是否删除" json:"is_deleted"`                           // 是否删除
	// Priority    int32     `gorm:"column:priority;not null;comment:0: 未设置 1: 紧急 2: 非紧急" json:"priority"`                // 0: 未设置 1: 紧急 2: 非紧急
	// CreatedAt   time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	// UpdatedAt   time.Time `gorm:"column:updated_at;not null;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}
type CreateModelParams struct {
	*Model
	TemplateIds []int64
}

func main() {
	greeter := &Greeter{
		Hello: "hello",
	}

	fmt.Println(*greeter, &greeter, reflect.TypeOf((greeter)), reflect.TypeOf((*greeter)))

	params := &CreateModelParams{
		&Model{
			ID: 1,
		},
		[]int64{1, 2, 3},
	}
	fmt.Println(params)
}
