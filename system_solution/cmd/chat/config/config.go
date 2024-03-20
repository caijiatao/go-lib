package config

type Config struct {
	ETCDEndpoints []string
	Env           *Env
	RPCServer     *RPCServer
}

type Env struct {
	Host string
	Port string
}

func (e *Env) GetTarget() string {
	return e.Host + ":" + e.Port
}

type RPCServer struct {
	Network string
}
