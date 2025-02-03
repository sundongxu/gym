package model

import (
	"time"
)

// MembershipType 会员卡类型
type MembershipType int

const (
	MembershipTypeExperience MembershipType = iota + 1 // 体验卡（次卡）
	MembershipTypeMonthly                              // 月卡
	MembershipTypeQuarterly                            // 季卡
	MembershipTypeYearly                               // 年卡
)

// MembershipStatus 会员状态
type MembershipStatus int

const (
	MembershipStatusActive   MembershipStatus = 1 // 生效中
	MembershipStatusExpired  MembershipStatus = 2 // 已过期
	MembershipStatusUsed     MembershipStatus = 3 // 已使用（仅针对体验卡）
	MembershipStatusCanceled MembershipStatus = 4 // 已取消
)

// MembershipBenefit 会员权益
type MembershipBenefit struct {
	Name        string `json:"name"`        // 权益名称
	Description string `json:"description"` // 权益描述
	Value       string `json:"value"`       // 权益值
}

// MembershipCard 会员卡
type MembershipCard struct {
	ID          uint           `gorm:"primary_key" json:"id"`
	Name        string         `gorm:"size:64" json:"name"`          // 会员卡名称
	Type        MembershipType `json:"type"`                         // 会员卡类型
	Price       float64        `json:"price"`                        // 价格
	Duration    int            `json:"duration"`                     // 有效期（天）
	Description string         `gorm:"type:text" json:"description"` // 描述
	Benefits    string         `gorm:"type:json" json:"benefits"`    // 会员权益（JSON格式存储，MembershipBenefit数组）
	Status      int            `gorm:"default:1" json:"status"`      // 状态：1-正常，2-下架
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// UserMembership 用户会员信息
type UserMembership struct {
	ID               uint             `gorm:"primary_key" json:"id"`
	UserID           uint             `json:"user_id"`                              // 用户ID
	MembershipCardID uint             `json:"membership_card_id"`                   // 会员卡ID
	OrderNo          string           `gorm:"size:32;unique_index" json:"order_no"` // 订单编号
	StartTime        time.Time        `json:"start_time"`                           // 会员开始时间
	EndTime          time.Time        `json:"end_time"`                             // 会员结束时间
	Status           MembershipStatus `gorm:"default:1" json:"status"`              // 会员状态
	PayAmount        float64          `json:"pay_amount"`                           // 支付金额
	PayMethod        string           `gorm:"size:16" json:"pay_method"`            // 支付方式
	PayTime          time.Time        `json:"pay_time"`                             // 支付时间
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}
