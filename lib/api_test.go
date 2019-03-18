package lib

import (
	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"net/http"
	"testing"
	"time"
)

func TestStartApi(t *testing.T) {
	err := LoadConfig("./config.json")
	if err != nil {
		t.Fatal(err)
	}
	httpHandler := getRoutes()
	corseHandler := cors.New(httpHandler)
	logger := logger.New(corseHandler, Config.LogLevel)
	go func(){
		err = http.ListenAndServe(":", logger)
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
}
