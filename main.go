package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/coopernurse/gorp"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Id    int64
	Name  string
	Email sql.NullString
}

func initialize() (db *sql.DB, err error) {
	userName := "root"
	password := "password"
	databaseKind := "mysql"
	databaseName := "db1"
	databaseHost := "localhost"
	databasePort := 3306

	viper.SetConfigName("database")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	errDb := viper.ReadInConfig()
	if errDb == nil {
		userName = viper.GetString("user.name")
		password = viper.GetString("user.password")
		databaseName = viper.GetString("db.name")
		databaseHost = viper.GetString("db.hostname")
		databasePort = viper.GetInt("db.port")
	} else {
		log.Println(errDb)
	}

	dataSourceName := userName + ":" + password + "@tcp(" + databaseHost + ":" + strconv.Itoa(databasePort) + ")/" + databaseName
	log.Println(dataSourceName)

	db, err = sql.Open(databaseKind, dataSourceName)
	if err != nil {
		log.Printf("[initialize]: %s", err.Error())
	}
	return
}

func mapping(dbmap *gorp.DbMap) {
	tableName := "People"
	structKeyName := "Id"
	t1 := dbmap.AddTableWithName(Person{}, tableName).SetKeys(true, structKeyName)
	t1.ColMap(structKeyName).Rename("p_id")
	t1.ColMap("Name").Rename("name")
	t1.ColMap("Email")
}

func main() {

	db, err0 := initialize()
	if err0 != nil {
		log.Printf("[main 0]: %s", err0.Error())
		return
	}
	dbmap := &gorp.DbMap{Db: db,
		Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	mapping(dbmap)

	aliceEMail := sql.NullString{"alice@exemple.com", true}
	alice := &Person{1, "Alice", aliceEMail}
	log.Println(alice)
	err1 := dbmap.Insert(alice)
	if err1 != nil {
		log.Printf("[main 1]: %s", err1.Error())
		return
	}

	var people []Person
	_, err2 := dbmap.Select(&people, "select * from people order by p_id")
	if err2 != nil {
		log.Printf("[main 2]: %s", err2.Error())
		return
	}
	fmt.Println(people)
	for x, p := range people {
		log.Printf("    %d: %v\n", x, p)
	}

	return
}
