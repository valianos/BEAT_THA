module myFirstProject/web_server

go 1.15

replace myFirstProject/logger => ../logger

replace myFirstProject/httpUtil => ../httpUtil

require (
	myFirstProject/httpUtil v0.0.0-00010101000000-000000000000
	myFirstProject/logger v0.0.0-00010101000000-000000000000
)
