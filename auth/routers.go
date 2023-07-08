package auth

import (
	"github.com/campus-fora/mail"
	"github.com/campus-fora/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController AuthController
}

func NewAuthRouteController(authController AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(mail_channel chan mail.Mail, rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/register", rc.authController.signUpHandler(mail_channel))
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/refresh", rc.authController.RefreshAccessToken)
	router.GET("/logout", middleware.Authenticator(), rc.authController.LogoutUser)
	router.GET("/verifyemail/:verificationCode", rc.authController.VerificationLinkHandler(mail_channel))

}

type UserRouteController struct {
	userController UserController
}

func NewRouteUserController(userController UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("/users")
	router.GET("/whoami", middleware.Authenticator(), whoamiHandler)
}
