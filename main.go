package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

// Configurations exported
type Configurations struct {
	Database DatabaseConfigurations
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBName     string
	DBHost     string
	DBUser     string
	DBPort     int
	DBPassword string
}

type UserModel struct{
	Id int `gorm:"primary_key";"AUTO_INCREMENT"`
	Name string `gorm:"size:255"`
	Address string `gorm:"type:varchar(100)"`
   }
var CIPHER_KEY []byte

func main() {
	var configuration Configurations
	CIPHER_KEY = []byte("0123456789012345")
	var db *gorm.DB
	var DbTypeMysql = "mysql"
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}



	db, err = setupMysqlDb(DbTypeMysql, configuration.Database)

	if err != nil {
		log.Panic(err)
	}
	log.Println("Connection Established")
	
	if db.HasTable(&UserModel{}) == false {
	 	createusertable(db)
	}
	
	user := UserModel{Name:"John",Address:"New York"}
	
	encryptinsert(db,user)

	reaaduserdata(db,user)

	db.Close()

}


func setupMysqlDb(DbTypeMysql string, dbinfo DatabaseConfigurations) (*gorm.DB, error) {

	connStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True", dbinfo.DBUser, dbinfo.DBPassword, dbinfo.DBHost, dbinfo.DBPort, dbinfo.DBName)
	return gorm.Open(DbTypeMysql, connStr)
}

func createusertable(db *gorm.DB) {
	db.DropTableIfExists(&UserModel{})
 	db.AutoMigrate(&UserModel{})

}

func encryptinsert(db *gorm.DB, newUser UserModel){


	fmt.Println(newUser)
		// Reading variables using the model

	usrName, err := Encrypt(CIPHER_KEY, newUser.Name)
	newUser.Name=usrName
	if err != nil {
		log.Println("failed to encrypt uaasecret", err)
		
	}

	 //Also we can use save that will return primary key
	 db.Create(&newUser)
}

func reaaduserdata(db *gorm.DB,user UserModel) {

	users := []UserModel{}
	db.Find(&users)
	
	fmt.Println("There are", len(users), "user records in the table.")

	for _,user := range users {
		usrName, err := Decrypt(CIPHER_KEY, user.Name)
		if err != nil {
			log.Println("failed to Decrypt Name", err)
			
		}
		 user.Name = usrName;
		 log.Println(user)
	}
}


