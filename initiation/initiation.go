package initiation

import (
	"flag"
	"log"
)

type InitParams struct {
	Token      string
	AdminLogin string
}

func InitiateParams() InitParams {
	i := InitParams{}
	adminLogin := adminLoginSet()
	token := tokenSet()

	flag.Parse()

	i.AdminLogin = *adminLogin
	i.Token = *token

	if i.Token == "" {
		log.Fatal("Token not found")
	}
	if i.AdminLogin == "" {
		log.Fatal("AdminLogin not found")
	}
	return i
}

func tokenSet() *string {
	//botAuc -tgTokenBot 'my token'
	token := flag.String(
		"tgTokenBot",
		"",
		"Token for Tg send-bot",
	)
	return token
}

func adminLoginSet() *string {
	//botAuc -tgTokenBot 'my token'
	AdminLogin := flag.String(
		"tgAdminLogin",
		"",
		"userID, who can have extended rights",
	)

	return AdminLogin
}
