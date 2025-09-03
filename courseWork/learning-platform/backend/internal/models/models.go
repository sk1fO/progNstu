package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string   `gorm:"uniqueIndex;not null"`
	Email    string   `gorm:"uniqueIndex;not null"`
	Password string   `gorm:"not null"`
	Role     string   `gorm:"not null"` // "teacher" или "student"
	Courses  []Course `gorm:"foreignKey:TeacherID"`
}

type Course struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	TeacherID   uint `gorm:"not null"`
	Teacher     User
	Lessons     []Lesson
	Enrollments []Enrollment
}

type Enrollment struct {
	gorm.Model
	UserID   uint `gorm:"not null"`
	CourseID uint `gorm:"not null"`
	User     User
	Course   Course
}

type Lesson struct {
	gorm.Model
	Title         string `gorm:"not null"`
	Order         int    `gorm:"not null"`
	TheoryContent string `gorm:"type:text"`
	CourseID      uint   `gorm:"not null"`
	Course        Course
	Assignments   []Assignment
}

type Assignment struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"type:text"`
	StarterCode string `gorm:"type:text"`
	Language    string `gorm:"not null;default:'python'"`
	LessonID    uint   `gorm:"not null"`
	Lesson      Lesson
	TestCases   []TestCase
	Solutions   []Solution
}

type TestCase struct {
	gorm.Model
	InputData      string `gorm:"type:text"`
	ExpectedOutput string `gorm:"type:text;not null"`
	IsVisible      bool   `gorm:"default:false"`
	AssignmentID   uint   `gorm:"not null"`
	Assignment     Assignment
}

type Solution struct {
	gorm.Model
	Code            string `gorm:"type:text;not null"`
	AssignmentID    uint   `gorm:"not null"`
	Assignment      Assignment
	UserID          uint `gorm:"not null"`
	User            User
	Status          string `gorm:"default:'sent'"`
	PassedAutotests bool   `gorm:"default:false"`
	TeacherComment  string `gorm:"type:text"`
	Score           int
}
