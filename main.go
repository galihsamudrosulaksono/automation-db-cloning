package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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
	dbname   = "sditqan"
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

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func main() {
	db, err := sql.Open("mysql", dsn("sditqan"))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelfunc()

	// //get users
	usersV2 := []User{}

	rows, err := db.QueryContext(ctx, "SELECT * FROM data_user ORDER BY id ASC")
	if err != nil {
		log.Printf("Error %s when getting users", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email *string
		var username *string
		var password *string
		var passcode *string
		var username_new *string
		var password_new *string
		var modified_date_userpass *time.Time
		var approve_user *int
		var username_ori *string
		var password_ori *string
		var expire *[]uint8
		var avatar *string
		var alias *string
		var nama *string
		var gender *string
		var tentang *string
		var role *string
		var upd_email_time *int
		var upd_email_new *string
		var upd_email_code *string
		var xdat *string
		var login_last *[]uint8
		var modified *[]uint8
		var modified_id *int
		var session_id *string
		var deviceid *string
		var absensi_sid *string
		var perpus_sid *string
		var temp1 *string
		var temp2 *string
		var google_id *string
		var google_email *string
		var google_name *string
		var google_picture *string
		var google_modified *[]uint8
		var setting_presensi_ip *int
		var setting_presensi_qrcode *int
		var setting_presensi_radius *int
		var setting_presensi_foto *int
		var verification *int
		var verification_modified *[]uint8
		var verificator *int

		//
		err = rows.Scan(&id,
			&email,
			&username,
			&password,
			&passcode,
			&username_new,
			&password_new,
			&modified_date_userpass,
			&approve_user,
			&username_ori,
			&password_ori,
			&expire,
			&avatar,
			&alias,
			&nama,
			&gender,
			&tentang,
			&role,
			&upd_email_time,
			&upd_email_new,
			&upd_email_code,
			&xdat,
			&login_last,
			&modified,
			&modified_id,
			&session_id,
			&deviceid,
			&absensi_sid,
			&perpus_sid,
			&temp1,
			&temp2,
			&google_id,
			&google_email, &google_name, &google_picture, &google_modified, &setting_presensi_ip, &setting_presensi_qrcode, &setting_presensi_radius, &setting_presensi_foto, &verification, &verification_modified, &verificator)
		if err != nil {
			log.Printf("Error %s when scanning users", err)
			return
		}
		//append to users
		usersV2 = append(usersV2, User{
			ID:       id,
			Name:     *nama,
			Password: *password,
			Username: *username,
			Role:     *role,
		})
	}
	fmt.Println(usersV2, " user : ", len(usersV2))

	/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// postgresql

	dataSourceName := fmt.Sprintf("host=db-postgresql-do-user-12344224-0.b.db.ondigitalocean.com port=25060 user=doadmin dbname=itqan-sim-db sslmode=require sslrootcert=config/ca-certificate.crt password=AVNS_FaGscLRzgGbXRuffZRf")

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
	for _, user := range usersV2 {
		//generate uuid for user
		newId, err := uuid.NewUUID()
		if err != nil {
			log.Printf("Error %s when generating uuid", err)
			return
		}

		//role
		if user.Role == "sdm" {
			user.Role = "teacher"
		} else if user.Role == "siswa" {
			user.Role = "student"
		}

		//generate email from username
		email := user.Username + "@afresto.com"

		//check if username already exist in postgresql dont insert
		var count int
		err = dbp.QueryRow("SELECT COUNT(id) FROM users WHERE username = $1", user.Username).Scan(&count)
		if err != nil {
			log.Fatalln(err)
		}
		// print count
		fmt.Println(count, " count : ", user.Username)
		if count == 0 {
			// print count
			fmt.Println(user.Username + " dibuat")
			//insert to postgresql
			_, err = dbp.Exec("INSERT INTO users (id, name, password, username, email, role, id_v2, created_at, updated_at, avatar) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", newId, user.Name, user.Password, user.Username, email, user.Role, user.ID, time.Now(), time.Now(), "user-avatar/5.png")
			if err != nil {
				log.Fatalln(err)
			}
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
