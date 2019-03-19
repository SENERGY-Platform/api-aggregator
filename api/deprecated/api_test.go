package deprecated

import (
	"github.com/SmartEnergyPlatform/api-aggregator/lib"
	"github.com/SmartEnergyPlatform/util/http/cors"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"net/http"
	"testing"
	"time"
)

func TestStartApi(t *testing.T) {
	config, err := lib.LoadConfig("../../config.json")
	if err != nil {
		t.Fatal(err)
	}
	httpHandler := GetRoutes(lib.New(config))
	corseHandler := cors.New(httpHandler)
	logger := logger.New(corseHandler, config.LogLevel)
	go func() {
		err = http.ListenAndServe(":", logger)
		if err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
}
