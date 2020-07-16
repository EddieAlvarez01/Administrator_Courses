package main

import (
	"context"
	"fmt"

	"github.com/EddieAlvarez01/administrator_courses/dao/mongodb"
)

func main() {
	db := mongodb.GetConnection()
	defer db.Disconnect(context.TODO())
	fmt.Println("Database is connected")
}
