package main

import (
    "fmt"
    "log"

    "github.com/astaxie/beego"
    "github.com/astaxie/beego/orm"
    _ "github.com/go-sql-driver/mysql"

    _ "beego-jwt-vergo/models"  // register models via init()
    _ "beego-jwt-vergo/routers" // register routes via init()
)

func init() {
    // Read DB configuration from app.conf
    dbUser := beego.AppConfig.String("db_user")
    dbPass := beego.AppConfig.String("db_pass")
    dbHost := beego.AppConfig.String("db_host")
    dbName := beego.AppConfig.String("db_name")
    dbPort := beego.AppConfig.String("db_port")

    if dbUser == "" || dbName == "" {
        log.Fatal("database configuration is missing in conf/app.conf")
    }

    if dbHost == "" {
        dbHost = "localhost"
    }
    if dbPort == "" {
        dbPort = "3306"
    }

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
        dbUser, dbPass, dbHost, dbPort, dbName)

    fmt.Println("Using DSN:", dsn)

    orm.RegisterDriver("mysql", orm.DRMySQL)
    if err := orm.RegisterDataBase("default", "mysql", dsn); err != nil {
        log.Fatalf("failed to register database: %v", err)
    }

    // Synchronise the schema (only in dev)
    if err := orm.RunSyncdb("default", false, true); err != nil {
        log.Fatalf("failed to sync database schema: %v", err)
    }
}

func main() {
    beego.Run()
}
