package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "chatapp"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	userConnections = make(map[string]*websocket.Conn)

	db *gorm.DB
)

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	if !db.Migrator().HasTable(&User{}) {
		db.AutoMigrate(&User{})
	}

	if !db.Migrator().HasTable(&Message{}) {
		db.AutoMigrate(&Message{})
	}

	if !db.Migrator().HasTable(&Group{}) {
		db.AutoMigrate(&Group{})
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/signup", SignUpHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")

	r.HandleFunc("/ws/{username}", WebSocketHandler)

	r.HandleFunc("/connections/{username}", AddConnection).Methods("POST")

	r.HandleFunc("/private_messages/{sender}/{receiver}", GetPrivateMessagesHandler).Methods("GET")
	r.HandleFunc("/private_messages/{sender}/{receiver}", CreatePrivateMessageHandler).Methods("POST")

	r.HandleFunc("/group/{username}", CreateGroupHandler).Methods("POST")
	r.HandleFunc("/group/{username}", GetUserGroupsHandler).Methods("GET")

	r.HandleFunc("/messages/{groupID}", GetGroupMessagesHandler).Methods("GET")
	r.HandleFunc("/messages/{groupID}", CreateGroupMessageHandler).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	http.Handle("/", r)

	fmt.Println("starting server at :8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Message model
type Message struct {
	gorm.Model
	senderName   string
	GroupID      uint   // if group message
	ReceiverName string // if private message
	Text         string
}

// User model
type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Password    string
	Groups      []Group `gorm:"many2many:group_members"`
	Connections string
}

// Group model
type Group struct {
	gorm.Model
	Name        string
	Members     []User `gorm:"many2many:group_members"`
	Messages    []Message
	CreatedByID uint
	CreatedBy   User `gorm:"foreignKey:CreatedByID;"`
}

// GroupRequest is used for decoding JSON when creating a new group
type GroupRequest struct {
	Name      string   `json:"name"`
	Usernames []string `json:"members"`
}

type AddConnectionRequest struct {
	Username string `json:"username"`
}

// SignUpHandler handles user registration.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	users := []User{}
	db.Where("username = ?", user.Username).Find(&users)
	if len(users) > 0 {
		RespondWithError(w, http.StatusConflict, "Username alerady exists")
		return
	}

	db.Create(&user)
	RespondWithJSON(w, http.StatusCreated, user)
}

// LoginHandler handles user login.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string
		Password string
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var user User
	db.Where("username = ?", credentials.Username).First(&user)
	if user.ID == 0 || user.Password != credentials.Password {
		RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}

// WebSocketHandler handles WebSocket connections.
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	fmt.Println("FromWebSocket>>>>>>>>>>>>>", username)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		switch messageType {
		case websocket.TextMessage:
			handleWebSocketTextMessage(p)
		}
	}
}

// GetGroupMessagesHandler handles retrieving group messages.
func GetGroupMessagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]

	var messages []Message
	db.Where("group_id = ?", groupID).Find(&messages)

	RespondWithJSON(w, http.StatusOK, messages)
}

// CreateGroupMessageHandler handles creating and broadcasting group messages.
func CreateGroupMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["groupID"]

	grpId, err := strconv.Atoi(groupID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid groupid")
		return
	}

	var message Message
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	message.GroupID = uint(grpId)
	fmt.Printf(">>>>>>>>>>>>%+v\n", message)
	db.Create(&message)

	broadcastGroupMessage(message.GroupID, message)

	RespondWithJSON(w, http.StatusCreated, message)
}

// GetPrivateMessagesHandler handles retrieving private messages between two users.
func GetPrivateMessagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderName := vars["sender"]

	var sender User
	result := db.Where("username = ?", senderName).First(&sender)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", senderName))
		return
	}

	receiverName := vars["receiver"]

	var receiver User
	result = db.Where("username = ?", receiverName).First(&receiver)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", receiverName))
		return
	}

	var messages []Message
	db.Where("(user_id = ? AND receiver_id = ?) OR (user_id = ? AND receiver_id = ?)", sender.ID, receiver.ID, receiver.ID, sender.ID).Find(&messages)

	RespondWithJSON(w, http.StatusOK, messages)
}

// CreatePrivateMessageHandler handles creating and sending private messages.
func CreatePrivateMessageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	senderName := vars["sender"]

	var sender User
	result := db.Where("username = ?", senderName).First(&sender)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", senderName))
		return
	}

	receiverName := vars["receiver"]

	var receiver User
	result = db.Where("username = ?", receiverName).First(&receiver)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", receiverName))
		return
	}

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	message.senderName = sender.Username
	message.ReceiverName = receiver.Username
	fmt.Printf(">>>>>>>>>>>>Private: %+v\n", message)
	db.Create(&message)

	sendPrivateMessage(strconv.Itoa(int(receiver.ID)), message)

	RespondWithJSON(w, http.StatusCreated, message)
}

// CreateGroupHandler handles creating a new group.
func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", username))
		return
	}

	var groupReq GroupRequest
	err := json.NewDecoder(r.Body).Decode(&groupReq)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	members := []User{user}
	for _, username = range groupReq.Usernames {
		result := db.Where("username = ?", username).First(&user)
		if result.Error != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", username))
			return
		}

		members = append(members, user)
	}

	// Create a new group
	newGroup := Group{
		Name:        groupReq.Name,
		Members:     members,
		CreatedByID: user.ID,
		CreatedBy:   user,
	}

	db.Create(&newGroup)

	fmt.Printf("all memebers when creating the group %s: %+v\n", groupReq.Name, members)

	RespondWithJSON(w, http.StatusCreated, newGroup)
}

// GetUserGroupsHandler handles getting all user groups.
func GetUserGroupsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	// Check if the user exists
	var user User
	result := db.Where("username = ?", username).Preload("Groups").First(&user)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", username))
		return
	}

	RespondWithJSON(w, http.StatusOK, user.Groups)
}

// AddConnection handles creating a new user connection.
func AddConnection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	// Check if the user exists
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Username %s not found", username))
		return
	}

	var connctionReq AddConnectionRequest
	err := json.NewDecoder(r.Body).Decode(&connctionReq)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var connection User
	result = db.Where("username = ?", connctionReq.Username).First(&connection)
	if result.Error != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("username %s not found", connctionReq.Username))
		return
	}

	connections := map[string][]string{}

	if user.Connections != "" {
		err = json.Unmarshal([]byte(user.Connections), &connections)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "processing error")
			return
		}
	}

	connectionData, ok := connections["connectionData"]
	if ok {
		for _, cnctn := range connectionData {
			if cnctn == connctionReq.Username {
				RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Connection with username %s already exists", cnctn))
				return
			}
		}
	}

	connectionData = append(connectionData, connctionReq.Username)
	connections["connectionData"] = connectionData

	bytes, err := json.Marshal(connections)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "processing error")
		return
	}

	user.Connections = string(bytes)

	db.Save(&user)
	RespondWithJSON(w, http.StatusOK, user)

}

// handleWebSocketTextMessage handles text messages received over WebSocket.
func handleWebSocketTextMessage(data []byte) error {
	// Handle the text message (e.g., store it in the database)
	var msg Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	db.Create(&msg)

	// Broadcast the message to all connected clients in the group
	broadcastGroupMessage(msg.GroupID, msg)

	return nil
}

// sendPrivateMessage sends a private message to the specified user.
func sendPrivateMessage(receiverID string, message Message) {
	// Check if the receiver is connected
	conn, ok := userConnections[receiverID]
	if !ok {
		log.Println("User not connected:", receiverID)
		return
	}

	// Convert the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message:", err)
		return
	}

	// Send the message to the receiver
	err = conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		log.Println("Error sending private message:", err)
	}
}

// broadcastGroupMessage broadcasts a group message to all connected users in the group.
func broadcastGroupMessage(groupID uint, message Message) {
	// Get all users in the group
	var group Group
	result := db.Preload("Members").Where("id = ?", groupID).First(&group)
	if result.Error != nil {
		log.Println("Error retrieving group:", result.Error)
		return
	}

	// Convert the message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message:", err)
		return
	}

	for _, member := range group.Members {
		conn, ok := userConnections[member.Username]
		if ok {
			err := conn.WriteMessage(websocket.TextMessage, messageJSON)
			if err != nil {
				log.Println("Error broadcasting group message:", err)
			}
		}
	}
}

// RespondWithError sends an error response with the specified status code and message.
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON sends a JSON response with the specified status code and data.
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
