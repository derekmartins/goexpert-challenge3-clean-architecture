package main

import (
	"database/sql"
	"fmt"
	"github.com/derekmartins/goexpert-challenge3-clean-architecture/configs"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	create_database_query := "CREATE TABLE orders (id varchar(255), price float NOT NULL, tax float NOT NULL, final_price float, PRIMARY KEY (ID))"

	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(create_database_query)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO orders (id, price, tax, final_price) VALUES ("uuid1", "40.0", "5.0", "45.0")`)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`INSERT INTO orders (id, price, tax, final_price) VALUES ("uuid2", "45.0", "5.0", "50.0")`)

	if err != nil {
		panic(err)
	}

}
