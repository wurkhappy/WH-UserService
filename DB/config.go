package DB

type config struct{
	DBName string
	DBURL string
}

Config := config{DBName: "UserDB", DBURL: "localhost:27017"}
