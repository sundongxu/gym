package database

import (
	"gym/internal/model"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&model.User{},
		&model.TrainerProfile{},
		&model.CustomerProfile{},
		&model.CheckInRecord{},
		&model.Course{},
		&model.CourseSchedule{},
		&model.CourseOrder{},
		&model.CourseEnrollment{},
		&model.MembershipCard{},
		&model.UserMembership{},
	)
	return err
}
