package model

import (
	"time"
)

// UserRole 用户角色
type UserRole int

const (
	RoleCustomer UserRole = iota + 1 // 普通用户
	RoleTrainer                      // 教练
	RoleManager                      // 店长
)

// LoginType 登录类型
type LoginType = string

const (
	LoginTypePassword  = "password"    // 密码登录
	LoginTypePhone     = "phone"       // 手机验证码登录
	LoginTypeWxMiniApp = "wx_mini_app" // 微信小程序登录
	LoginTypeWxApp     = "wx_app"      // 微信APP登录
)

// User 用户基本信息
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Phone     string     `gorm:"size:16;unique_index" json:"phone"` // 手机号
	Password  string     `gorm:"size:128" json:"-"`                 // 密码（可选）
	Nickname  string     `gorm:"size:32" json:"nickname"`           // 昵称
	Avatar    string     `gorm:"size:255" json:"avatar"`            // 头像
	Gender    int        `gorm:"default:0" json:"gender"`           // 性别：0-未知，1-男，2-女
	Role      UserRole   `gorm:"default:1" json:"role"`             // 用户角色
	Status    int        `gorm:"default:1" json:"status"`           // 状态：1-正常，2-禁用
	LastLogin *time.Time `gorm:"default:null" json:"last_login"`    // 最后登录时间
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// UserAuth 用户授权信息
type UserAuth struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	UserID     uint      `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user_id"` // 关联用户ID
	LoginType  LoginType `gorm:"size:32;index;not null;comment:'登录类型'" json:"login_type"`            // 登录类型
	AuthKey    string    `gorm:"size:64;index;not null;comment:'认证标识'" json:"auth_key"`              // 认证标识（手机号/openid）
	AuthSecret string    `gorm:"size:255;comment:'认证密钥'" json:"-"`                                   // 认证密钥（密码hash/验证码/session_key/token）
	Extra      string    `gorm:"type:text;comment:'额外认证信息'" json:"-"`                                // 额外认证信息（JSON格式，unionid/scope等）
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User User `gorm:"foreignkey:UserID" json:"user"` // 关联用户信息
}

// TrainerProfile 教练信息
type TrainerProfile struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `gorm:"unique_index" json:"user_id"`  // 关联用户ID
	Description string    `gorm:"type:text" json:"description"` // 个人简介
	Experience  int       `json:"experience"`                   // 教龄(年)
	Speciality  string    `gorm:"size:255" json:"speciality"`   // 专长领域
	Certificate string    `gorm:"size:255" json:"certificate"`  // 资质证书
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	User User `gorm:"foreignkey:UserID" json:"user"` // 关联用户信息
}

// CustomerProfile 会员信息
type CustomerProfile struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	UserID       uint      `gorm:"unique_index" json:"user_id"`   // 关联用户ID
	Age          int       `json:"age"`                           // 年龄
	Height       float64   `json:"height"`                        // 身高(cm)
	Weight       float64   `json:"weight"`                        // 体重(kg)
	HealthStatus string    `gorm:"size:255" json:"health_status"` // 健康状况
	Goals        string    `gorm:"size:255" json:"goals"`         // 健身目标
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	User User `gorm:"foreignkey:UserID" json:"user"` // 关联用户信息
}

// CheckInRecord 签到记录
type CheckInRecord struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `json:"user_id"`     // 关联用户ID
	CheckInAt time.Time `json:"check_in_at"` // 签到时间
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignkey:UserID" json:"user"` // 关联用户信息
}
