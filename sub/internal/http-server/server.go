package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	"test/internal/cache"
	. "test/internal/model"
	pg "test/internal/pgconn"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Cache *cache.Cache
	DB    *pg.DB
}

func NewServer() *Server {
	return &Server{
		Cache: cache.InitCache(),
	}
}

func (s *Server) Start(ctx context.Context) {
	var err error

	s.DB, err = pg.InitPG(ctx, "postgresql://kirill:tsbgfsvwasdqw@localhost:5432/orders?sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}

	defer s.DB.Close()

	dir, _ := os.Executable()
	flpth := path.Join(path.Dir(dir), "/template/index.html")

	router := gin.Default()
	router.LoadHTMLGlob(flpth)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
		})
	})
	router.GET("/getOrder/:id", s.getOrder)
	router.POST("/order", s.saveOrder)
	router.Run("localhost:8100")
}

func (s *Server) getOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := s.DB.GetSingleMessage(id)
	if err != nil {
		c.IndentedJSON(400, "Order doesn't exist")
		return
	}
	c.IndentedJSON(200, order)
}

func (s *Server) saveOrder(c *gin.Context) {

	var order OrderMessage

	if err := c.BindJSON(&order); err != nil {
		c.String(400, "Something went wrong")
		return
	}
	if err := order.Validate(); err != nil {
		c.String(400, "Order is not valid")
		return
	}
	if err := s.DB.InsertMessage(order); err != nil {
		c.String(400, err.Error())
		return
	}
	s.Cache.Add(order)
	c.String(200, "Success")
}
