package impl

import (
	"context"
	"golib/libs/id_generate/service/impl/snowflake_generator"
	"golib/libs/test_libs"
	"math/rand"
	"testing"
	"time"
)

var (
	testIdGenerateService *IdGenerateService
	testCtx               context.Context
)

func initEnv() {
	test_libs.SetDebugMode()
	rand.Seed(time.Now().Unix())
}

func TestMain(m *testing.M) {
	initEnv()
	testIdGenerateService = NewIdGenerateService(snowflake_generator.NewSnowflakeGenerateService())
	m.Run()
}
