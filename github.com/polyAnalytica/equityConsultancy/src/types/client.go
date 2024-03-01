package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/polyAnalytica/equityConsultancy/src/db"
	"gorm.io/gorm"
)

type client struct {
	gorm.Model
	Id       uint   `gorm:"primarykey; autoIncrement:true; unique"`
	Name     string `gorm:"size:255;not null;unique" json:"name"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	PassHash string `gorm:"size:255;not null;" json:"passHash"`
}

func (cl client) TableName() string {
	return "clients"
}

func init() {
	cl := &client{}
	if err := db.Migrate(cl); err != nil {
		log.Fatal(err)
	}
}

func NewClient(clientData []byte) (*client, error) {
	newClient := &client{}
	err := json.Unmarshal(clientData, newClient)
	if err != nil {
		return nil, err
	}

	conn, err := db.Connect()
	if err != nil {
		return nil, err
	}
	conn.Create(newClient)
	fmt.Println(newClient)

	return newClient, nil
}

func LoadClient(email string, passHash string) (*client, error) {
	client := &client{}
	conn, err := db.Connect()
	if err != nil {
		return nil, err
	}

	conn.Where("email=?", email).Find(client)
	if client.PassHash != passHash {
		return nil, errors.New("invalid password")
	}

	return client, nil
}
