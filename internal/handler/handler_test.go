package handler

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/steinfletcher/apitest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env")
	code := m.Run()
	os.Exit(code)
}

func TestPutQueue_Success(t *testing.T) {
	baseUrl := os.Getenv("BASE_URL")
	t.Log(baseUrl)
	apitest.New().
		Put(baseUrl + "/queue/" + os.Getenv("QUEUE_NAME")).
		JSON(`{"message": "` + os.Getenv("TEST_MESSAGE") + `"}`).
		Expect(t).
		Status(200).
		End()
}

func TestPutQueue_QueueFull(t *testing.T) {
	for i := range 5 { // QUEUE_CAPACITY=3
		apitest.New().
			Put(os.Getenv("BASE_URL") + "/queue/overflow").
			JSON(fmt.Sprintf(`{"message": "msg%d"}`, i)).
			Expect(t).
			Status(503).
			End()
	}
}

func TestGetQueue_Success(t *testing.T) {
	apitest.New().
		Get(os.Getenv("BASE_URL") + "/queue/" + os.Getenv("QUEUE_NAME")).
		Expect(t).
		Body(`{"message": "` + os.Getenv("TEST_MESSAGE") + `"}`).
		Status(200).
		End()
}

func TestGetQueue_Timeout(t *testing.T) {
	apitest.New().
		Get(os.Getenv("BASE_URL")+"/queue/empty").
		Query("timeout", "1").
		Expect(t).
		Status(404).
		End()
}

func TestGetQueue_InvalidParams(t *testing.T) {
	apitest.New().
		Get(os.Getenv("BASE_URL") + "/queue/").
		Expect(t).
		Status(404).
		End()
}
