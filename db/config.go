package db

var (
	pool *connection
)

// Connect connects to database.
func Connect() {
	pool = newConnection()
}

// Disconnect closes the database connection.
func Disconnect() {
	pool.disconnect()
}
