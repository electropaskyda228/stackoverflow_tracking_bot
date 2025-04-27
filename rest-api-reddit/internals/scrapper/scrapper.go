package scrapper

import (
	"log"
	"net/http"
	"net/url"
	. "rest-api-reddit/internals/api"
	. "rest-api-reddit/internals/common"
	. "rest-api-reddit/internals/db/postgres"
	"strconv"
	"time"
)

func waitMotherFucker() {
	time.Sleep(3 * time.Second)
}

/*
ЧТО ПРОИСХОДИТ В ЭТОЙ ФУНКЦИИ
Это функция, которая смотрит обновления на сайте. Предполагается, что это функция будет запущена в отдельной горутине.

Для начала создается новое постоянное соединение с базой данных. Потом запускается бесконечный цикл.
В нем на каждой новой итерации выгружается информация из базы данных в виде map[uint][]uint, где ключ - это id вопроса,
а значения списка - id пользователей (они же id чатов с этими пользователями в tg), кроме первого значения, обозначающего
время последней активности на данном вопросе. Я решил выгружать данные в таком странном виде для уменьшения занимаемой памяти
с сохранением удобной для последующей работы структурой.

Далее для каждого вопроса делается запрос на api stackoverflow с целью обновления информации.
Если время последней активности не изменилось, то переходим к следующему вопросу. В противном случае
передаем в канал обмена сначала "куклу", в поле QuestionId которой записано количество пользователей.
После передаем в канал id-шники этих пользователей. "Кукла" сделана для того, чтобы читатель канала
сразу записал все значения для данного вопроса, а структура Tracking выбрана для того, чтобы спокойно передавать
до двух значений.

Далее информация о вопросе обновляется в бд.

После обработки каждого вопроса функция зависает на 3 секунды, чтобы избежать чрезмерных запросов в базу данных и
в api stackoverflow. API вообще дает 30k запросов в день, то есть [3600 * 24 / 30k] + 1 = 3 секунды на один запрос.
К тому же обновления на каждом вопросе не слишком частое явление. Плюс пользователю нашего приложения не так уже и
важно сразу получать уведомление об обновлении на вопросе (далее вообще планируется дайджест событий). Поэтому
на мой взгляд (пока я в пьяном угаре), что данная ленивая обработка обновлений является оптимальной.
*/
func ListenTrackings(in chan<- Tracking) {
	connection := NewConnection()
	for {
		log.Println("Сделал выгрузку базы. Начинаю обновлять информацию")
		trackings := GetAllTracking(connection)
		for questionId, usersId := range trackings {
			newQuestion, err := GetQuestion(UintToString(questionId))
			if err != nil {
				waitMotherFucker()
				continue
			}
			if newQuestion.LastActivityDate != usersId[0] {
				in <- Tracking{QuestionId: uint(len(usersId) - 1), UserId: 0}
				for _, userId := range usersId[1:] {
					in <- Tracking{QuestionId: questionId, UserId: userId}
				}
				UpdateQuestion(connection, newQuestion)
				log.Println("В вопросе", questionId, "есть обновление")
			}
			log.Println("Работа с вопросом", questionId, "закончена")
			waitMotherFucker()
		}
	}
}

const uri = "http://localhost:8081/"

func SendInfoToBot(out <-chan Tracking) {
	for {
		doll := <-out
		for i := uint(0); i < doll.QuestionId; i++ {
			msg := <-out
			client := http.Client{Timeout: 3 * time.Second}

			data := url.Values{}
			data.Add("question", strconv.FormatUint(uint64(msg.QuestionId), 10))
			data.Add("chat_id", strconv.FormatUint(uint64(msg.UserId), 10))

			resp, err := client.PostForm(uri+"bot/send", data)
			if err != nil {
				log.Println("Собшение об обновлении вопроса", msg.QuestionId, "не отправлено пользователю", msg.UserId)
				return
			}
			if resp.StatusCode != 200 {
				log.Println("Собшение об обновлении вопроса", msg.QuestionId, "не отправлено пользователю", msg.UserId)
				return
			}
			log.Println("Собшение об обновлении вопроса", msg.QuestionId, "отправлено пользователю", msg.UserId)
		}
	}
}
