package receiver

import (
	"testing"
	"os"
	"strconv"
)

func TestNewReceiver(t *testing.T) {
	port,_ := strconv.ParseInt(os.Getenv("REDIS_PORT"),10,32)
	db,_ := strconv.ParseInt(os.Getenv("REDIS_DB"),10,32)
	config := ReceiverConfig{
		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: int(port),
		RedisPwd: os.Getenv("REDIS_PASSWORD"),
		RedisDB: int(db),
	}

	NewReceiver(config)

}
