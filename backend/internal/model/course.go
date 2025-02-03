package model

import (
	"time"
)

// Course 课程信息
type Course struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Name        string    `gorm:"size:64" json:"name"`          // 课程名称
	Description string    `gorm:"type:text" json:"description"` // 课程描述
	TrainerID   uint      `json:"trainer_id"`                   // 教练ID
	Type        string    `gorm:"size:32" json:"type"`          // 课程类型
	Capacity    int       `json:"capacity"`                     // 课程容量
	Duration    int       `json:"duration"`                     // 课程时长(分钟)
	Price       float64   `json:"price"`                        // 课程价格
	Status      int       `gorm:"default:1" json:"status"`      // 状态：1-正常，2-下架
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CourseSchedule 课程排期
type CourseSchedule struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CourseID  uint      `json:"course_id"`               // 课程ID
	StartTime time.Time `json:"start_time"`              // 开始时间
	EndTime   time.Time `json:"end_time"`                // 结束时间
	Status    int       `gorm:"default:1" json:"status"` // 状态：1-未开始，2-进行中，3-已结束
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CourseOrder 课程订单
type CourseOrder struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	OrderNo   string    `gorm:"size:32;unique_index" json:"order_no"` // 订单编号
	UserID    uint      `json:"user_id"`                              // 用户ID
	CourseID  uint      `json:"course_id"`                            // 课程ID
	Amount    float64   `json:"amount"`                               // 订单金额
	Status    int       `gorm:"default:1" json:"status"`              // 状态：1-待支付，2-已支付，3-已取消
	PayMethod string    `gorm:"size:16" json:"pay_method"`            // 支付方式
	PayTime   time.Time `json:"pay_time"`                             // 支付时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CourseEnrollment 课程报名记录
type CourseEnrollment struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CourseID   uint      `json:"course_id"`               // 课程ID
	ScheduleID uint      `json:"schedule_id"`             // 课程排期ID
	UserID     uint      `json:"user_id"`                 // 用户ID
	Status     int       `gorm:"default:1" json:"status"` // 状态：1-已报名，2-已签到，3-已取消
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
