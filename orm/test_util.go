package orm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TestSuiteConfigOpt func(*TestSuiteConfig)

const (
	defaultTestDSN = iota + 1
)

var (
	testDSNMap = map[int]string{
		defaultTestDSN: "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	}
)

type TestSuiteConfig struct {
	Dsn int
}

func (config *TestSuiteConfig) GetDSN() string {
	dsn, ok := testDSNMap[config.Dsn]
	if !ok {
		return testDSNMap[defaultTestDSN]
	}
	return dsn
}

func TestWithDsnName(dsn int) TestSuiteConfigOpt {
	return func(config *TestSuiteConfig) {
		config.Dsn = dsn
	}
}

func InitTestSuite(opts ...TestSuiteConfigOpt) {
	config := &TestSuiteConfig{}
	for _, opt := range opts {
		opt(config)
	}
	dsn := config.GetDSN()
	err := NewOrmClient(&Config{
		Config: &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		},
		SourceConfig: &SourceDBConfig{},
		Dial:         postgres.Open(dsn),
	})
	if err != nil {
		panic(err)
	}
	if err = NewOrmClient(&Config{
		DBClientName: RecommendJobDBClientName,
		Config:       &gorm.Config{},
		SourceConfig: &SourceDBConfig{},
		Dial:         postgres.Open(dsn),
	}); err != nil {
		panic(err)
	}
}
