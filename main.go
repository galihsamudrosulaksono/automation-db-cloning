package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	_ "github.com/go-sql-driver/mysql"
	//use sqlx
	// pgx driver
)

const (
	username = "root"
	password = ""
	hostname = "127.0.0.1"
	dbname   = "smk7"
)

type PostgresConfig struct {
	PostgresqlHost         string
	PostgresqlPort         string
	PostgresqlUser         string
	PostgresqlPassword     string
	PostgresqlDbname       string
	PostgresqlSSLMode      bool
	PgDriver               string
	PostgresqlReplicaHosts []string
}

//create struct user
type User struct {
	ID       int    `db:"id"`
	Name     string `db:"nama"`
	Password string `db:"password"`
	Username string `db:"username"`
	Role     string `db:"role"`
}

// create user phone
type UserPhone struct {
	ID     int `db:"id"`
	UserID int `db:"id_user"`
	Phone  int `db:"nowa"`
}

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func main() {
	db, err := sql.Open("mysql", dsn("smk7"))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelfunc()

	// //get users
	// usersV2 := []User{}
	usersPhonesV2 := []UserPhone{}

	// just get user that not a student
	rows, err := db.QueryContext(ctx, "SELECT id, id_user, nowa FROM bukutamu_nowas ORDER BY id ASC")
	if err != nil {
		log.Printf("Error %s when getting users", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var id_user int
		var nowa *int

		//
		err = rows.Scan(&id,
			&id_user,
			&nowa)
		if err != nil {
			log.Printf("Error %s when scanning users", err)
			return
		}
		// //append to users
		// usersV2 = append(usersV2, User{
		// 	ID:       id,
		// 	Name:     *nama,
		// 	Password: *password,
		// 	Username: *username,
		// 	Role:     *role,
		// })
		// append to users phone
		usersPhonesV2 = append(usersPhonesV2, UserPhone{
			ID:     id,
			UserID: id_user,
			Phone:  *nowa,
		})

	}
	fmt.Println(usersPhonesV2, " user : ", len(usersPhonesV2))

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// postgresql

	dataSourceName := fmt.Sprintf("host=db-postgresql-do-user-12344224-0.b.db.ondigitalocean.com port=25060 user=doadmin dbname=letter-sim-db-smkn7smg sslmode=require sslrootcert=config/ca-certificate.crt password=AVNS_FaGscLRzgGbXRuffZRf")

	dbp, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	//check connection
	err = dbp.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	//insert usersv2 to postgresql
	for _, userPhone := range usersPhonesV2 {
		fmt.Println(userPhone.UserID, " ", userPhone.Phone, "di update")
		//update users phone
		_, err = dbp.Exec("UPDATE users SET phone = $1 WHERE id_v2 = $2", userPhone.Phone, userPhone.UserID)
		if err != nil {
			log.Printf("Error %s when updating users phone", err)
			return
		}
	}

	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// //get users
	// rowsP, err := dbp.QueryxContext(ctx, "SELECT * FROM users ORDER BY id ASC LIMIT 3")
	// if err != nil {
	// 	log.Printf("Error %s when getting users", err)
	// 	return
	// }

	// defer rowsP.Close()

	// for rowsP.Next() {
	// 	var id uuid.UUID
	// 	var name string
	// 	var email string
	// 	var avatar string
	// 	var email_verified_at time.Time
	// 	var password string
	// 	var role string
	// 	var remember_token *string
	// 	var created_at time.Time
	// 	var updated_at time.Time
	// 	var deleted_at *time.Time
	// 	var username string

	// 	err = rowsP.Scan(&id, &name, &email, &avatar, &email_verified_at, &password, &role, &remember_token, &created_at, &updated_at, &deleted_at, &username)
	// 	if err != nil {
	// 		log.Printf("Error %s when scanning users", err)
	// 		return
	// 	}

	// 	fmt.Println(id, name, email, avatar, email_verified_at, password, role, remember_token, created_at, updated_at, deleted_at, username)
	// }

	log.Printf("Connected to DB %s successfully\n", dbname)
}
