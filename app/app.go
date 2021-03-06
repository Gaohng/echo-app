package app

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/components"
	"github.com/gw123/echo-app/services"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var App *EchoApp


type EchoApp struct {
	areaSvc     echoapp.AreaService
	smsSvc      echoapp.SmsService
	UserSvr     echoapp.UserService
	dbPool      echoapp.DbPool
	resourceSvc echoapp.ResourceService
	goodsSvc    echoapp.GoodsService
	orderSvc    echoapp.OrderService
}

func init() {
	App = &EchoApp{}
}

func GetAreaService() (echoapp.AreaService, error) {
	if App.areaSvc != nil {
		return App.areaSvc, nil
	}
	areaSvc, err := services.NewAreaService(echoapp.ConfigOpts.Asset.AreaRoot)
	if err != nil {
		return nil, errors.Wrap(err, "GetAreaService")
	}
	App.areaSvc = areaSvc
	return areaSvc, nil
}

func MustGetAreaService() echoapp.AreaService {
	areaSvc, err := GetAreaService()
	if err != nil {
		panic(err)
	}
	return areaSvc
}

func GetSmsService() (echoapp.SmsService, error) {
	if App.smsSvc != nil {
		return App.smsSvc, nil
	}
	smsSvc := services.NewSmsService(echoapp.ConfigOpts.SmsOptionTokenMap)
	App.smsSvc = smsSvc
	return smsSvc, nil
}
func MustGetSmsService() echoapp.SmsService {
	areaSvc, err := GetSmsService()
	if err != nil {
		panic(err)
	}
	return areaSvc
}

func GetDb(dbname string) (*gorm.DB, error) {
	if App.dbPool != nil {
		return App.dbPool.Db(dbname)
	}
	dbPool := services.NewDbPool(echoapp.ConfigOpts.DBMap)
	App.dbPool = dbPool
	return dbPool.Db(dbname)
}

func MustGetDb(dbName string) *gorm.DB {
	db, err := GetDb(dbName)
	if err != nil {
		panic(errors.Wrap(err, "MustGetDb->GetDb"))
	}
	return db
}

func GetJwsHelper() (*components.JwsHelper, error) {
	jws, err := components.NewJwsHelper(echoapp.ConfigOpts.Jws)
	if err != nil {
		return nil, errors.Wrap(err, "NewJwsHelper")
	}
	return jws, nil
}

func MustGetJwsHelper() *components.JwsHelper {
	userSvr, err := GetJwsHelper()
	if err != nil {
		panic(errors.Wrap(err, "GetJwsHelper"))
	}
	return userSvr
}

func GetUserService() (echoapp.UserService, error) {
	if App.UserSvr != nil {
		return App.UserSvr, nil
	}
	userDb, err := GetDb("user")
	if err != nil {
		return nil, errors.Wrap(err, "GetUserSerevice->GetDb")
	}
	/*
	redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "GetRedis")
	}
*/
	jws, err := GetJwsHelper()
	if err != nil {
		return nil, errors.Wrap(err, "GetJws")
	}
	App.UserSvr = services.NewUserService(userDb, nil, jws)
	return App.UserSvr, nil
}

func MustUserService() echoapp.UserService {
	userSvr, err := GetUserService()
	if err != nil {
		panic(errors.Wrap(err, "MustuserSer-> GetUserSvr"))
	}
	return userSvr
}
func GetResourceService() (echoapp.ResourceService, error) {
	if App.resourceSvc != nil {
		return App.resourceSvc, nil
	}
	userDb, err := GetDb("user")
	if err != nil {
		return nil, errors.Wrap(err, "GetResourceService->GetDb")
	}
	App.resourceSvc = services.NewResourceService(userDb)
	return App.resourceSvc, nil
}
func MustGetResourceService() echoapp.ResourceService {
	resource, err := GetResourceService()
	if err != nil {
		panic(err)
	}
	return resource
}
func GetGoodsService() (echoapp.GoodsService, error) {
	if App.goodsSvc != nil {
		return App.goodsSvc, nil
	}
	db, err := GetDb("user")
	if err != nil {
		return nil, errors.Wrap(err, "GetGoodsService->GetDb")
	}
	App.goodsSvc = services.NewGoodsService(db)
	return App.goodsSvc, nil
}
func MustGetGoodsService() echoapp.GoodsService {
	goods, err := GetGoodsService()
	if err != nil {
		panic(err)
	}
	return goods
}
func GetOrderService() (echoapp.OrderService, error) {
	if App.orderSvc != nil {
		return App.orderSvc, nil
	}
	db, err := GetDb("user")
	if err != nil {
		return nil, errors.Wrap(err, "GetOrderService->GetDb")
	}
	App.orderSvc = services.NewOrderService(db)
	return App.orderSvc, nil
}
func MustGetOrderService() echoapp.OrderService {
	goods, err := GetOrderService()
	if err != nil {
		panic(err)
	}
	return goods
}
