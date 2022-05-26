package main

import (
	"fmt"
	"strconv"

	cart "github.com/yyystation/cart/proto"

	"github.com/go-micro/plugins/v4/registry/consul"
	ratelimit "github.com/go-micro/plugins/v4/wrapper/ratelimiter/uber"
	opentracing2 "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/yyystation/cart/domain/repository"
	service2 "github.com/yyystation/cart/domain/service"
	"github.com/yyystation/cart/handler"
	"github.com/yyystation/common"
	"go-micro.dev/v4"
	log "go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

var (
	service = "cart"
	version = "latest"
	address = "127.0.0.1:8087"
)

var QPS = 100

func main() {
	//配置中心
	consulConfig, err := common.GetConsulConfig("10.10.50.59", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//注册中心
	consul := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"10.10.50.59:8500",
		}
	})
	//链路追踪
	t, io, err := common.NewTracer("github.com/yyystation/cart", "10.10.50.59:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//获取mysql的配置
	mysqlInfo := common.GetMysqlFormConsul(consulConfig, "mysql")

	//初始化数据库连接
	dns := mysqlInfo.User + ":" + mysqlInfo.Pwd + "@tcp(" + mysqlInfo.Host + ":" + strconv.FormatInt(mysqlInfo.Port, 10) + ")/" + mysqlInfo.Database + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open("mysql", dns)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	//禁止副表
	db.SingularTable(true)

	//初始化
	// repository.NewCartRepository(db).InitTable()

	// Create service
	srv := micro.NewService(
		micro.Name(service),
		micro.Version(version),
		micro.Address(address),
		//注册中心
		micro.Registry(consul),
		//链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)
	srv.Init()

	cartDataService := service2.NewCartDataService(repository.NewCartRepository(db))

	// service2.RegisterCartHandler(service.Server())

	// Register handler
	cart.RegisterCartHandler(srv.Server(), &handler.Cart{CartDataService: cartDataService})

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
