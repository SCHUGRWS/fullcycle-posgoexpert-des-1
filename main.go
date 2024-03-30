package main

import (
	"github.com/SCHUGRWS/fullcycle-posgoexpert-des-1/ratelimiter"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func initConfig() {
	viper.SetConfigName("config") // Nome do arquivo de configuração (sem a extensão)
	viper.SetConfigType("yaml")   // Tipo do arquivo de configuração
	viper.AddConfigPath(".")      // Caminho para olhar o arquivo de configuração
	err := viper.ReadInConfig()   // Ler o arquivo de configuração
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	initConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	rateLimiter := ratelimiter.NewRateLimiter(redisClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	wrappedMux := rateLimiter.Limit(mux)

	http.ListenAndServe(":8080", wrappedMux)
}
