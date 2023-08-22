package mail

import (
	_ "github.com/campus-fora/config"
	"github.com/spf13/viper"
)

var (
	user string
	pass string
	host string
	port string
	//webteam string
	// batch   int
	sender string
)

func init() {

	user = viper.GetString("MAIL.LOGIN")
	sender = user

	pass = viper.GetString("MAIL.PSWD")
	host = viper.GetString("MAIL.HOST")
	port = viper.GetString("MAIL.PORT")
	//webteam = viper.GetString("MAIL.WEBTEAM")

	// batch = viper.GetInt("MAIL.BATCH")
}
