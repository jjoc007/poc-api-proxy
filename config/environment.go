package config

import "os"

const (
	BaseURLBackendServer = "http://localhost:8080"
	XSourceIPHeader      = "x-source-ip"
)

var (
	DataBaseDefinition = ConnectionData{
		Schema:    os.Getenv("SCHEMA"),
		Host:      os.Getenv("HOST"),
		ReadUser:  os.Getenv("READ_USER"),
		ReadPwd:   os.Getenv("READ_PASSWORD"),
		WriteUser: os.Getenv("WRITE_USER"),
		WritePwd:  os.Getenv("WRITE_PASSWORD"),
	}
)

type ConnectionData struct {
	Schema    string
	Host      string
	ReadUser  string
	ReadPwd   string
	WriteUser string
	WritePwd  string
}
