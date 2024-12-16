package server

import (
	"c_cache/consistent_hash"
	"c_cache/lruk"
	"github.com/gin-gonic/gin"
	"golib/libs/naming"
	"golib/libs/net_helper"
	"net/http"
	"strings"
	"sync"
)

type Server struct {
	engine *gin.Engine
	addr   string

	discSvc     *naming.ServiceDiscovery
	registerSvc *naming.ServiceRegister
	clusters    *consistent_hash.ConsistentHash

	cache *lruk.LRUK

	mu sync.Mutex
}

func NewServer(addr string) *Server {
	addr = net_helper.GetFigureOutListenOn(addr)

	cache, err := lruk.NewLRUK()
	if err != nil {
		panic(err)
	}

	svc := &Server{
		engine: gin.Default(),
		addr:   addr,

		discSvc:     naming.NewServiceDiscovery([]string{"127.0.0.1:2379"}, "cache"),
		registerSvc: naming.NewServiceRegister([]string{"127.0.0.1:2379"}, "cache", addr),
		clusters:    consistent_hash.NewConsistentHash(),

		cache: cache,

		mu: sync.Mutex{},
	}

	svc.registerRouter()

	go svc.watchClusters()

	return svc
}

func (s *Server) watchClusters() {
	for _ = range s.discSvc.Watch() {
		clusters := consistent_hash.NewConsistentHash()

		services := s.discSvc.GetServices()
		for _, service := range services {
			clusters.AddNode(service)
		}

		s.mu.Lock()
		s.clusters = clusters
		s.mu.Unlock()
	}
}

func (s *Server) registerRouter() {
	s.engine.GET("/ping", s.ping)

	s.engine.GET("/get/:key", s.get)

	s.engine.DELETE("/del/:key", s.del)

	s.engine.PUT("/put/:key", s.put)
}

func (s *Server) Run() error {
	err := s.engine.Run(s.addr)
	if err != nil {
		panic(err)
	}
	return nil
}

func (s *Server) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (s *Server) get(c *gin.Context) {
	key := c.Param("key")

	if value, ok := s.cache.Get(String(key)); ok {
		c.JSON(200, gin.H{
			"key":   key,
			"value": value,
		})
		return
	}

	node := s.clusters.GetNode(key)
	if node == "" || node == s.addr {
		c.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}

	resp, err := http.Get(node + "/get/" + key)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal error",
		})
		return
	}
	c.JSON(200, gin.H{
		"key":   key,
		"value": resp.Body,
	})
}

func (s *Server) del(c *gin.Context) {
	key := c.Param("key")

	node := s.clusters.GetNode(key)
	if node == s.addr {
		s.cache.Remove(String(key))
		c.JSON(200, gin.H{
			"message": "deleted",
		})
		return
	}

	resp, err := http.Post(node+"/del/"+key, "", nil)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": resp.Body,
	})

}

func (s *Server) put(c *gin.Context) {
	key := c.Param("key")
	value := c.PostForm("value")

	node := s.clusters.GetNode(key)
	if node == "" {
		c.JSON(404, gin.H{
			"message": "not found",
		})
		return
	}
	if node == s.addr {
		s.cache.Add(String(key), String(value))
		c.JSON(201, gin.H{
			"message": "created",
		})
		return
	}

	resp, err := http.Post(node+"/put/"+key, "application/x-www-form-urlencoded", strings.NewReader("value="+value))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "internal error",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": resp.Body,
	})
}
