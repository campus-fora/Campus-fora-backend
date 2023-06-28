package auth

import (
	"github.com/gin-gonic/gin"
	"campus-fora/middleware"
)

type AuthRouteController struct {
	authController AuthController
}

func NewAuthRouteController(authController AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", rc.authController.signUpHandler)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.Authenticator(), rc.authController.LogoutUser)
}

type UserRouteController struct {
	userController UserController
}

func NewRouteUserController(userController UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("/users")
	router.GET("/whoami",  middleware.Authenticator(), whoamiHandler)
}

