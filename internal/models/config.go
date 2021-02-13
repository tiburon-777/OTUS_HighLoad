package models

type Configuration struct {
	Server	Server
	DSN		DSN
}

type Server struct {
	Address string
	Port    string
}

type DSN struct {
	Master	string
	Slave1	string
	Slave2	string
	Port	string
	User	string
	Pass	string
	Base	string
}
