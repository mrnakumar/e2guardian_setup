package pkg

import (
	"flag"
	"github.com/mrnakumar/e2g_utils"
	"log"
	"strings"
)

type ParsedFlag struct {
	BasePath         string
	Domain           string
	UserName         string
	Password         string
	IdentityFilePath string
	DevelopMode      bool
}

func ParseFlags() ParsedFlag {
	basePathFlag := "base_path"
	domainFlag := "domain"
	userNameFlag := "user_name"
	passwordFlag := "password"
	identityFilePathFlag := "identity_file_path"
	developModeFlag := "develop_mode"
	basePath := flag.String(basePathFlag, "", "Base folder path to store data")
	domain := flag.String(domainFlag, "", "Domain name to serve")
	user := flag.String(userNameFlag, "", "Username")
	password := flag.String(passwordFlag, "", "Password")
	identityFilePath := flag.String(identityFilePathFlag, "", "Identity (X25519) private key file path")
	developMode := flag.String(developModeFlag, "true", "Run the server in develop mode.[true | false] Default is true")

	flag.Parse()

	domainTrimmed := strings.TrimSpace(*domain)
	if len(domainTrimmed) == 0 {
		log.Fatalf("flag '%s' is required", domainFlag)
	}
	return ParsedFlag{
		BasePath:         e2g_utils.ValidatePath(basePath, basePathFlag),
		Domain:           e2g_utils.Base64DecodeWithKill(domainTrimmed, domainFlag),
		UserName:         e2g_utils.Base64DecodeWithKill(*user, userNameFlag),
		Password:         e2g_utils.ParsePassword(*password, passwordFlag),
		IdentityFilePath: e2g_utils.ValidatePath(identityFilePath, identityFilePathFlag),
		DevelopMode:      *developMode == "true",
	}
}
