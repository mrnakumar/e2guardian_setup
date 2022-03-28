package pkg

import (
	"flag"
	"log"
	"os"
	"strings"
)

type ParsedFlag struct {
	BasePath string
	Domain   string
	UserName string
	Password string
}

func ParseFlags() ParsedFlag {
	basePathFlag := "base_path"
	domainFlag := "domain"
	userNameFlag := "user_name"
	passwordFlag := "password"
	basePath := flag.String(basePathFlag, "", "Base folder path to store data")
	domain := flag.String(domainFlag, "", "Domain name to serve")
	user := flag.String(userNameFlag, "", "Username")
	password := flag.String(passwordFlag, "", "Password")

	flag.Parse()

	trimmed := strings.TrimSpace(*basePath)
	if len(trimmed) == 0 {
		log.Fatalf("flag '%s' is required", basePathFlag)
	}
	if _, err := os.Stat(trimmed); os.IsNotExist(err) {
		log.Fatalf("invalid '%s'. Path '%s' does not exist.", basePathFlag, trimmed)
	}
	domainTrimmed := strings.TrimSpace(*domain)
	if len(domainTrimmed) == 0 {
		log.Fatalf("flag '%s' is required", domainFlag)
	}
	return ParsedFlag{
		BasePath: trimmed,
		Domain:   domainTrimmed,
		UserName: *user,
		Password: *password,
	}
}
