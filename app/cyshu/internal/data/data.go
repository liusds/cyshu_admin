package data

import (
	"cyshu_admin/app/cyshu/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo, NewAdminRepo)

// Data .
type Data struct {
	db          *gorm.DB
	redisClient *redis.Client
	// TODO wrapped database client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	//mysql ..
	dsn := c.Database.Source + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.NewHelper(logger).Info("mysql server start...")

	//redis ..
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: "",
		DB:       0,
	})
	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.NewHelper(logger).Info("redis server start...")

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		db:          db,
		redisClient: client,
	}, cleanup, nil
}
