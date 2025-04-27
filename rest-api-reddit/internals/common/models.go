package common

type Question struct {
	ID               uint   `json:"-" gorm:"primaryKey"`
	IsAnswered       bool   `json:"is_answered" gorm:"column:is_answered"`
	AnswerCount      int    `json:"answer_count" gorm:"column:answer_count"`
	ClosedDate       uint   `json:"closed_date" gorm:"column:closed_date"`
	LastActivityDate uint   `json:"last_activity_date" gorm:"column:last_activity_date"`
	Link             string `json:"link" gorm:"column:link"`
	Title            string `json:"title" gorm:"column:title"`
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	UserName string `gorm:"column:user_name"`
}

type Tracking struct {
	UserId     uint `gorm:"column:user_id"`
	QuestionId uint `gorm:"column:question_id"`
}
