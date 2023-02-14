package database

import (
	"fmt"
	"golang/constant"
	"golang/errs"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetCPDbClient() (*sqlx.DB, *errs.AppErrorvalidation) {

	cpDBHost := constant.CP_DB_ADDR
	cpDBUser := constant.CP_DB_USER
	cpDbPassword := constant.CP_DB_PASS
	cpDbReaderName := constant.CP_DB_NAME

	dataSourse := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cpDBHost, cpDBUser, cpDbPassword, cpDbReaderName)

	client, err1 := sqlx.Open("postgres", dataSourse)
	if err1 != nil {
		panic(err1)
	}

	// Confirm a successful connection.
	if err := client.Ping(); err != nil {
		// log.Fatal(err)
		fmt.Println("error in db connection")
		return nil, errs.ValidateResponse(nil, 500, "Unexpected error from database")
		panic(err)
	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client, nil

}

func BTMSDBConnection() (*sqlx.DB, *errs.AppErrorvalidation) {

	dbUser := constant.BC_DB_USER // DB username
	dbPass := constant.BC_DB_PASS // DB Password
	dbHost := constant.BC_DB_HOST // DB Hostname/IP
	dbName := constant.BC_DB_NAME // Database name

	datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", dbUser, dbPass, dbHost, dbName)

	log.Println("datasource :" + datasource)

	client, err1 := sqlx.Open("mysql", datasource)

	if err1 != nil {
		log.Println("error client", err1)
		panic(err1)
	}
	fmt.Println("confirm db connection")
	// Confirm a successful connection.
	if err := client.Ping(); err != nil {
		// log.Fatal(err)
		log.Println("error client ping", err)

		fmt.Println("error in db connection")
		return nil, errs.ValidateResponse(nil, 500, "Error while trying to connect to BTMS database")

		//panic(err)

	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	log.Println("exited successfully")

	return client, nil
}

func BTMSDBWriterConnection() (*sqlx.DB, *errs.AppError) {

	dbUser := constant.BC_DB_USER // DB username
	dbPass := constant.BC_DB_PASS // DB Password
	dbHost := constant.BC_DB_HOST // DB Hostname/IP
	dbName := constant.BC_DB_NAME // Database name

	datasource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", dbUser, dbPass, dbHost, dbName)

	log.Println("datasource :" + datasource)

	client, err1 := sqlx.Open("mysql", datasource)

	if err1 != nil {
		log.Println("error client", err1)
		panic(err1)
	}
	fmt.Println("confirm db connection")
	// Confirm a successful connection.
	if err := client.Ping(); err != nil {
		// log.Fatal(err)
		log.Println("error client ping", err)

		fmt.Println("error in db connection")
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		//panic(err)

	}

	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	log.Println("exited successfully")

	return client, nil
}

func GetCpdbClient() (*sqlx.DB, *errs.AppErrorvalidation) {

	cpDBHost := constant.CP_DB_ADDR
	cpDBUser := constant.CP_DB_USER
	cpDbPassword := constant.CP_DB_PASS
	cpDbReaderName := constant.CP_DB_NAME

	dataSourse := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", cpDBHost, cpDBUser, cpDbPassword, cpDbReaderName)
	fmt.Println(dataSourse)
	client, err1 := sqlx.Open("postgres", dataSourse)
	// db, err = sql.Open("mysql", connString)
	if err1 != nil {
		panic(err1)
	}
	fmt.Println("confirm db connection")
	// Confirm a successful connection.
	if err := client.Ping(); err != nil {
		// log.Fatal(err)
		fmt.Println("error in db connection")
		return nil, errs.ValidateResponse(nil, 500, "Error while trying to connect to Customer Portal database")
		panic(err)
	}
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client, nil
}

func GetSharedDbClient() (*sqlx.DB, *errs.AppError) {

	dbUser := constant.SH_DB_USER
	dbAddr := constant.SH_DB_ADDR

	dbPass := constant.SH_DB_PASS
	dbName := constant.SH_DB_NAME

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbAddr,
		dbUser, dbPass, dbName)
	fmt.Println(dataSource)
	client, err1 := sqlx.Open("postgres", dataSource)
	// db, err = sql.Open("mysql", connString)

	if err1 != nil {
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
	}
	fmt.Println("confirm db connection")
	// Confirm a successful connection.
	if err := client.Ping(); err != nil {
		// log.Fatal(err)
		fmt.Println("error in db connection")
		return nil, errs.NewUnexpectedError("Unexpected error from database", "500")
		//panic(err)
	}
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)
	return client, nil
}
