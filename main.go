package main

import (
	"github.com/SCHUGRWS/fullcycle-posgoexpert-des-1/ratelimiter"
	redisStore "github.com/SCHUGRWS/fullcycle-posgoexpert-des-1/ratelimiter/store/redis"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	initConfig()

	store := redisStore.NewRedisStore(nil)

	rateLimiter := ratelimiter.NewRateLimiter(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world from Rate Limeterlandia!"))
	})

	wrappedMux := rateLimiter.Limit(mux)

	http.ListenAndServe(":8080", wrappedMux)
}
