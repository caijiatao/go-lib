package impl

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"golib/common"
	"golib/id_generate/service/impl/snowflake_generator"
)

var (
	testIdGenerateService *IdGenerateService
	testCtx               context.Context
)

func initEnv() {
	common.SetDebugMode()
	rand.Seed(time.Now().Unix())
}

func TestMain(m *testing.M) {
	initEnv()
	testIdGenerateService = NewIdGenerateService(snowflake_generator.NewSnowflakeGenerateService())
	m.Run()
}
