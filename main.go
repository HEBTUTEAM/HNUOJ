package main

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/router"
	"github.com/HEBNUOJ/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func main() {
	utils.InitDbConfig()
	db := common.InitDB()
	client := common.InitRedis()
	defer db.Close()
	defer client.Close()
	r := router.CollectAuthorizeRoute(gin.Default())
	r = router.CollectVerifyRoute(r)
	r = router.CollectProblemRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())
}
