package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"

	_ "beego-jwt-vergo/models"  // IMPORT IMPORTANT POUR LES MODELS
	_ "beego-jwt-vergo/routers" // IMPORT IMPORTANT POUR LES ROUTES
)

func init() {

	// Lecture config depuis app.conf
	host := beego.AppConfig.String("db_host")
	user := beego.AppConfig.String("db_user")
	pass := beego.AppConfig.String("db_pass")
	name := beego.AppConfig.String("db_name")
	port := beego.AppConfig.String("db_port")

	// DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, pass, host, port, name)

	fmt.Println("DSN utilisé :", dsn)

	// Driver & DB
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", dsn)

	// Sync DB – crée automatiquement la table users
	orm.RunSyncdb("default", false, true)
}

func main() {
	beego.Run()
}
