package middleware

import (
	"fmt"
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/model"
	"github.com/HEBNUOJ/response"
	"github.com/HEBNUOJ/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取jwtToken和refresToken
		jwtToken := ctx.GetHeader("Authorization")
		refreshToken := ctx.GetHeader("RefreshToken")

		// 验证格式, 判断token是否以"Bearer "为前缀
		if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "jwtToken格式错误",
			})
			ctx.Abort() // 抛弃该次请求
			return
		}

		// 验证jwtToken和refreshToken是否有效
		jwtToken = jwtToken[7:]
		token, claims, err := common.ParseToken(jwtToken)
		flag, _ := common.GetRedisClient().Get(refreshToken).Result()
		blackToken, _ := common.GetRedisClient().Get(jwtToken).Result() // jwt是否在黑名单中

		if !token.Valid && len(flag) == 0 { //   如果jwtToken不合法，并且refreshToken也不在redis中
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// 验证通过后获取的claim中的userId
		userID := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userID)

		if user.Id == 0 || user.Id != userID {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "用户权限不足",
			})
			ctx.Abort()
			return
		}

		// 用户存在，将user信息写入上下文
		ctx.Set("user", user)
		fmt.Sprintf("%s", err)
		if len(flag) > 0 && (len(blackToken) == 0 || strings.Contains(err.Error(), "expired")) {
			ctx.Set("renewal", "true")
		} else {
			ctx.Set("renewal", "false")
		}

		ctx.Next()
	}
}

func RenewalTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next() // 先执行auth鉴权中间件，如果鉴权成功会继续往下执行
		// 如果不用续签，则直接退出即可
		flag, _ := ctx.Get("renewal")
		if flag == "false" {
			ctx.Next()
		}

		// 取出鉴权成功但需要续签的user对象
		tmp, _ := ctx.Get("user")
		user, _ := tmp.(model.User)

		jwtToken := ctx.GetHeader("Authorization")
		common.GetRedisClient().Set(jwtToken, 1, 10*time.Minute)
		// 续签token
		token, err := common.ReleaseToken(user)
		if err != nil {
			response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")
			utils.Log("token.log", 5).Println(err) // 记录错误日志
			return
		}
		ctx.Writer.Header().Set("token", token)
		ctx.Next()
	}
}
