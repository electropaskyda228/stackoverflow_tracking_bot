package postgres

import (
	"log"
	"os"
	. "rest-api-reddit/internals/common"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	db  *gorm.DB
	dsn string
}

var dbGlobal *Connection = nil

func checkConstants() {
	if os.Getenv("SERVER_HOST") == "" || os.Getenv("POSTGRES_USER") == "" || os.Getenv("POSTGRES_PASSWORD") == "" {
		log.Println("Os's constants have not been found")
		os.Exit(1)
	}
}

func MakeDB() {
	checkConstants()
	dsn := "host=" + os.Getenv("SERVER_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Question{})
	db.AutoMigrate(&Tracking{})

	dbGlobal = &Connection{db, dsn}
}

func NewConnection() *Connection {
	var newConnection Connection
	newConnection.dsn = dbGlobal.dsn
	newConnection.db = dbGlobal.db.Session(&gorm.Session{})
	return &newConnection
}

func AddQuestion(connection *Connection, question *Question) error {
	return connection.db.Create(question).Error
}

func AddUser(connection *Connection, user *User) error {
	return connection.db.Create(user).Error
}

func AddTracking(connection *Connection, userIds uint, questionIds uint) error {
	var tracking Tracking = Tracking{UserId: userIds, QuestionId: questionIds}
	return connection.db.Create(tracking).Error
}

func FindQuestion(connection *Connection, ids string) (*Question, error) {
	var question Question
	err := connection.db.First(&question, ids).Error
	return &question, err
}

func FindUser(connection *Connection, ids string) (*User, error) {
	var user User
	err := connection.db.First(&user, ids).Error
	return &user, err
}

func FindTracking(connection *Connection, questionIds string, userIds string) (*Tracking, error) {
	var tracking Tracking
	err := connection.db.Where("user_id = ? AND question_id = ?", userIds, questionIds).First(&tracking).Error
	return &tracking, err
}

func FindAllTrackingsByUser(connection *Connection, userIds string) ([]Tracking, error) {
	trackings := make([]Tracking, 0)
	err := connection.db.Where("user_id = ?", userIds).Find(&trackings).Error
	return trackings, err
}

func DeleteQuestion(connection *Connection, ids string) error {
	question, err := FindQuestion(connection, ids)
	if err != nil {
		return err
	}

	return connection.db.Delete(question).Error
}

func DeleteUser(connection *Connection, ids string) error {
	user, err := FindUser(connection, ids)
	if err != nil {
		return err
	}

	return connection.db.Delete(user).Error
}

func DeleteTracking(connection *Connection, questionIds string, userIds string) error {
	return connection.db.Where("user_id = ? AND question_id = ?", userIds, questionIds).Delete(&Tracking{}).Error
}

func CountQuestion(connection *Connection, questionIds string) int64 {
	var count int64 = -1
	connection.db.Model(&Tracking{}).Where("question_id = ?", questionIds).Count(&count)
	return count
}

func CountUser(connection *Connection, userIds string) int64 {
	var count int64 = -1
	connection.db.Model(&Tracking{}).Where("user_id = ?", userIds).Count(&count)
	return count
}

func CheckExistingUser(connection *Connection, userIds string) bool {
	user, err := FindUser(connection, userIds)
	return err == nil && user != nil
}

func CheckExistingQuestion(connection *Connection, questionIds string) bool {
	question, err := FindQuestion(connection, questionIds)
	return err == nil && question != nil
}

func CheckExistingTracking(connection *Connection, userIds string, questionIds string) bool {
	tracking, err := FindTracking(connection, questionIds, userIds)
	return err == nil && tracking != nil
}

func GetAllTracking(connection *Connection) map[uint][]uint {
	questions := GetAllQuestions(connection)
	if questions == nil {
		return nil
	}
	answer := make(map[uint][]uint)
	for _, question := range questions {
		answer[question.ID] = make([]uint, 1)
		answer[question.ID][0] = question.LastActivityDate
	}

	questions = nil

	var trackings []Tracking
	err := connection.db.Find(&trackings).Error
	if err != nil {
		return nil
	}

	for _, track := range trackings {
		answer[track.QuestionId] = append(answer[track.QuestionId], track.UserId)
	}

	return answer

}

func GetAllQuestions(connection *Connection) []*Question {
	var questions []Question
	err := connection.db.Find(&questions).Error
	if err != nil {
		return nil
	}
	answer := make([]*Question, len(questions))
	for i, question := range questions {
		answer[i] = &question
	}
	return answer
}

func UpdateQuestion(connection *Connection, question *Question) {
	connection.db.Save(question)
}
