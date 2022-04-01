package pkg

import (
	"flag"
	"log"
	"os"
	"strings"
)

type ParsedFlag struct {
	BasePath         string
	Domain           string
	UserName         string
	Password         string
	IdentityFilePath string
}

func ParseFlags() ParsedFlag {
	basePathFlag := "base_path"
	domainFlag := "domain"
	userNameFlag := "user_name"
	passwordFlag := "password"
	identityFilePathFlag := "identity_file_path"
	basePath := flag.String(basePathFlag, "", "Base folder path to store data")
	domain := flag.String(domainFlag, "", "Domain name to serve")
	user := flag.String(userNameFlag, "", "Username")
	password := flag.String(passwordFlag, "", "Password")
	identityFilePath := flag.String(identityFilePathFlag, "", "Identity (X25519) private key file path")

	flag.Parse()

	domainTrimmed := strings.TrimSpace(*domain)
	if len(domainTrimmed) == 0 {
		log.Fatalf("flag '%s' is required", domainFlag)
	}
	return ParsedFlag{
		BasePath:         validatePath(basePath, basePathFlag),
		Domain:           domainTrimmed,
		UserName:         *user,
		Password:         *password,
		IdentityFilePath: validatePath(identityFilePath, identityFilePathFlag),
	}
}

func validatePath(path *string, flag string) string {
	trimmed := strings.TrimSpace(*path)
	if len(trimmed) == 0 {
		log.Fatalf("flag '%s' is required", flag)
	}
	if _, err := os.Stat(trimmed); os.IsNotExist(err) {
		log.Fatalf("invalid '%s'. Path '%s' does not exist.", flag, trimmed)
	}
	return trimmed
}
