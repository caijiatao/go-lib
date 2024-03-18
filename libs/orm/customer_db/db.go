package customer_db

type DB interface {
	Decode()
	Encode()
	Discovery()
}
