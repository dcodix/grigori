package writer_redis

import (
	"github.com/dcodix/grigori/lib/domain/message"
	"github.com/dcodix/grigori/lib/domain/writer"
	"github.com/dcodix/grigori/lib/infrastructure/register_modules"
	"gopkg.in/redis.v5"
	"log"
	"time"
)

type WriterRedis struct {
	redis_host   string
	redis_port   string
	redis_key    string
	redis_client *redis.Client
}

// init registers with the program.
func init() {
	writer := new(WriterRedis)
	register_modules.RegisterWriter("redis", writer)
}

func (w *WriterRedis) Config(config map[string]interface{}) {
	if redis_host, ok := config["redis_host"]; ok {
		w.redis_host = redis_host.(string)
	} else {
		w.redis_host = "localhost"
	}
	if redis_port, ok := config["redis_port"]; ok {
		w.redis_port = redis_port.(string)
	} else {
		w.redis_port = "6379"
	}
	if redis_key, ok := config["redis_key"]; ok {
		w.redis_key = redis_key.(string)
	} else {
		w.redis_key = "grigori"
	}
	log.Println("Writing to REDIS: " + w.redis_key)
}

func (w *WriterRedis) Clone() writer.Writer {
	cloneWriter := new(WriterRedis)
	return cloneWriter
}

func (w *WriterRedis) connectToBackend() {
	var opt redis.Options
	opt.Addr = w.redis_host + ":" + w.redis_port
	opt.Network = "tcp"
	w.redis_client = redis.NewClient(&opt)
	pong, _ := w.redis_client.Ping().Result()
	if pong != "PONG" {
		log.Println("Unable to conect ro redis.")
	} else {
		log.Println("Redis conection successful.")
	}
}

func (w *WriterRedis) reconectToBackend() {
	log.Println("Connection to redis lost.(sleep)")
	time.Sleep(1 * time.Second)
	log.Println("Reconecting to redis...")
	w.connectToBackend()

}

func (w *WriterRedis) Write(messages chan message.Message) {
	debug := false

	log.Println("Connecting to redis...")
	w.connectToBackend()
	for message := range messages {
		//time.Sleep(1000 * time.Microsecond) //Testing monitoring
		redisaction, err := w.redis_client.RPush(w.redis_key, message.Message).Result()
		if err != nil {
			w.reconectToBackend()
		}
		if debug {
			log.Println(redisaction)
		}
	}
}
