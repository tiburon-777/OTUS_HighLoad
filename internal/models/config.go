package models

type Configuration struct {
	Server Server
	DSN    DSN
}

type Server struct {
	Address string
	Port    string
}

type DSN struct {
	Host string
	Port string
	User string
	Pass string
	Base string
}
