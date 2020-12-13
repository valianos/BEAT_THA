package protocol

const (
	SERVICE_A     PROVIDER = "ETAServiceA"
	SERVICE_A_URL string   = "http://localhost:8001/eta/calculate"
	SERVICE_B     PROVIDER = "ETAServiceB"
	SERVICE_B_URL string   = "http://localhost:8002/calculateETA"
	UNSPECIFIED   PROVIDER = ""
)

const (
	GET  METHOD = "GET"
	POST METHOD = "POST"
)

const (
	BAD_REQUEST        int = 400
	SERVER_ERROR       int = 500
	METHOD_NOT_ALLOWED int = 405
)
