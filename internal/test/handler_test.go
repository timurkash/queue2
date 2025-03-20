package test

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	"github.com/timurkash/queue2/internal/biz"
	"github.com/timurkash/queue2/internal/handler"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

var (
	queue                  string
	queueCap               int
	apiTestPut, apiTestGet *apitest.APITest
)

func TestMain(m *testing.M) {
	//_ = godotenv.Load(".env.test")
	configMap, _ := godotenv.Read(".env.test")

	queue = configMap["QUEUE_NAME"]
	maxQueues, _ := strconv.Atoi(configMap["MAX_QUEUES"])
	queueCap, _ = strconv.Atoi(configMap["QUEUE_CAPACITY"])
	srv := biz.New(maxQueues, queueCap)
	theHandler := handler.New(srv, time.Second)
	apiTestPut = apitest.New().HandlerFunc(theHandler.PutQueue)
	apiTestGet = apitest.New().HandlerFunc(theHandler.GetQueue)
	code := m.Run()
	os.Exit(code)
}

func TestPutQueue_Success(t *testing.T) {
	apiTestPut.
		Put("/queue/" + queue).
		JSON(`{"message": "` + os.Getenv("TEST_MESSAGE") + `"}`).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestPutQueue_QueueFull(t *testing.T) {
	for i := range queueCap { // QUEUE_CAPACITY=3
		apiTestPut.
			Put("/queue/overflow").
			JSON(fmt.Sprintf(`{"message": "msg%d"}`, i)).
			Expect(t).
			Status(200).
			End()
	}
	apiTestPut.
		Put("/queue/overflow").
		JSON(fmt.Sprintf(`{"message": "msg%d"}`, 4)).
		Expect(t).
		Status(503).
		End()

}

func TestGetQueue_Success(t *testing.T) {
	apiTestPut.
		Put("/queue/" + queue).
		JSON(`{"message": "some"}`).
		Expect(t).
		Status(200).
		End()
	apiTestGet.
		Get("/queue/" + queue).
		Expect(t).
		//Body(`{"message": "` + os.Getenv("TEST_MESSAGE") + `"}`).
		Status(http.StatusOK).
		End()
}

func TestGetQueue_Timeout(t *testing.T) {
	apiTestGet.
		Get("/queue/empty").
		Query("timeout", "1").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestGetQueue_InvalidParams(t *testing.T) {
	apiTestGet.
		Get("/queue/abc").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
