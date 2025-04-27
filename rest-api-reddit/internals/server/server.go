package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	. "rest-api-reddit/internals/api"
	. "rest-api-reddit/internals/common"
	. "rest-api-reddit/internals/db/postgres"
	. "rest-api-reddit/internals/scrapper"
)

func trackQuestion(w http.ResponseWriter, req *http.Request) {
	body := "Method is not allowed"
	status := http.StatusBadRequest
	w.Header().Set("Content-Type", "text/plain")
	if req.Method != "POST" {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}
	body = "Failed to track question"
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	params := req.Form
	ids := params.Get("id")
	user := params.Get("user")
	connection := NewConnection()

	body = fmt.Sprintf("Failed to track question with ids %s", ids)

	if isTracking := CheckExistingTracking(connection, user, ids); isTracking {
		body = fmt.Sprintf("Question %s is already tracking", ids)
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	userUint, err1 := StringToUint(user)
	idsUint, err2 := StringToUint(ids)
	if err1 != nil || err2 != nil {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	if isUser := CheckExistingUser(connection, user); !isUser {
		AddUser(connection, &User{ID: userUint, UserName: "admin"})
	}

	if isQuestion := CheckExistingQuestion(connection, ids); !isQuestion {
		question, err := GetQuestion(ids)
		if err != nil {
			body = fmt.Sprintf("StackOverFlow is dead %s", ids)
			w.WriteHeader(status)
			w.Write([]byte(body))
			return
		}
		question.ID = idsUint
		AddQuestion(connection, question)
	}

	err = AddTracking(connection, userUint, idsUint)
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}
	body = fmt.Sprintf("Tracked question with ids %s", ids)
	status = http.StatusOK
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func untrackQuestion(w http.ResponseWriter, req *http.Request) {
	body := "Method is not allowed"
	status := http.StatusBadRequest
	w.Header().Set("Content-Type", "text/plain")
	if req.Method != "POST" {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}
	body = "Failed to untrack question"
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	params := req.Form
	ids := params.Get("ids")
	user := params.Get("user")
	connection := NewConnection()

	body = fmt.Sprintf("Failed to untrack question with ids %s", ids)

	if isTracking := CheckExistingTracking(connection, user, ids); !isTracking {
		body = fmt.Sprintf("Question %s is already untracked", ids)
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	err = DeleteTracking(connection, ids, user)
	if err != nil {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	if CountQuestion(connection, ids) == 0 {
		DeleteQuestion(connection, ids)
	}

	if CountUser(connection, user) == 0 {
		DeleteUser(connection, user)
	}

	body = fmt.Sprintf("Untracked question with ids %s", ids)
	status = http.StatusOK
	w.WriteHeader(status)
	w.Write([]byte(body))
}

func checkUserExisting(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	user := params.Get("user")
	connection := NewConnection()

	body := "false"
	status := http.StatusOK
	w.Header().Set("Content-Type", "text/plain")

	if isUser := CheckExistingUser(connection, user); isUser {
		body = "true"
	}

	w.WriteHeader(status)
	w.Write([]byte(body))

}

func getTrackingsListByUser(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	user := params.Get("user")
	connection := NewConnection()

	var body []byte
	status := http.StatusOK
	w.Header().Set("Content-Type", "application/json")

	trackings, err := FindAllTrackingsByUser(connection, user)
	if err != nil {
		status = http.StatusBadRequest
	} else {
		body, err = json.Marshal(trackings)
		if err != nil {
			status = http.StatusBadRequest
			body = nil
		}
	}

	w.WriteHeader(status)
	w.Write(body)
}

func createUser(w http.ResponseWriter, req *http.Request) {
	body := "Method is not allowed"
	status := http.StatusBadRequest
	w.Header().Set("Content-Type", "text/plain")
	if req.Method != "POST" {
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	err := req.ParseForm()
	if err != nil {
		body = "Failed to add user"
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	params := req.Form
	userId := params.Get("id")
	name := params.Get("name")
	connection := NewConnection()

	userIdStr, err := StringToUint(userId)
	if err != nil {
		log.Println("User", name, "is not added")
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}
	err = AddUser(connection, &User{ID: userIdStr, UserName: name})
	if err != nil {
		log.Println("User", name, "is not added")
		w.WriteHeader(status)
		w.Write([]byte(body))
		return
	}

	body = fmt.Sprintf("User %s is added", name)
	log.Println("User", name, "is added")
	status = http.StatusOK
	w.WriteHeader(status)
	w.Write([]byte(body))

}

func checkConstants() {
	if os.Getenv("STACKOVERFLOW_SERVER_PORT") == "" {
		log.Println("Os's constants have not been found")
		os.Exit(1)
	}
	CheckConstants()
}

func MakeServer() {
	checkConstants()
	http.HandleFunc("/question/track", trackQuestion)
	http.HandleFunc("/question/untrack", untrackQuestion)
	http.HandleFunc("/user/checking", checkUserExisting)
	http.HandleFunc("/user/add", createUser)
	http.HandleFunc("/question/list", getTrackingsListByUser)

	MakeDB()

	cannal := make(chan Tracking)

	go ListenTrackings(cannal)
	go SendInfoToBot(cannal)

	http.ListenAndServe(os.Getenv("STACKOVERFLOW_SERVER_PORT"), nil)
}
