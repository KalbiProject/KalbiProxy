package models

import (
	"fmt"
	"sync"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
  )

var lock = &sync.Mutex{}

var Database *gorm.DB

//NewConnection creates or returns *gorm.DB instance 
func NewConnection() *gorm.DB{
	lock.Lock()
	defer lock.Unlock()

	if Database == nil {
		db, err := gorm.Open(sqlite.Open("sipproxy.db"), &gorm.Config{})
		if err != nil{
			fmt.Println(err)
		}
		Database = db
	}

    return Database
  }