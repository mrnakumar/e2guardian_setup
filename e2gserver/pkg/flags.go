package pkg

import (
	"flag"
	"github.com/mrnakumar/e2g_utils"
	"log"
	"strconv"
	"strings"
)

type ParsedFlag struct {
	BasePath         string
	Domain           string
	UserName         string
	Password         string
	IdentityFilePath string
	DevelopMode      bool
	StorageLimit     uint16
}

func ParseFlags() ParsedFlag {
	basePathFlag := "base_path"
	domainFlag := "domain"
	userNameFlag := "user_name"
	passwordFlag := "password"
	identityFilePathFlag := "identity_file_path"
	developModeFlag := "develop_mode"
	storageLimitFlag := "storage_limit"
	basePath := flag.String(basePathFlag, "", "Base folder path to store data")
	domain := flag.String(domainFlag, "", "Domain name to serve")
	user := flag.String(userNameFlag, "", "Username")
	password := flag.String(passwordFlag, "", "Password")
	identityFilePath := flag.String(identityFilePathFlag, "", "Identity (X25519) private key file path")
	developMode := flag.String(developModeFlag, "true", "Run the server in develop mode.[true | false] Default is true")
	storageLimit := flag.String(storageLimitFlag, "", "Storage limit size in Megabyte")

	flag.Parse()

	domainTrimmed := strings.TrimSpace(*domain)
	if len(domainTrimmed) == 0 {
		log.Fatalf("flag '%s' is required", domainFlag)
	}
	storageLimitDecoded := e2g_utils.Base64DecodeWithKill(*storageLimit, storageLimitFlag)
	limit, err := strconv.Atoi(storageLimitDecoded)
	if err != nil || limit <= 10 || limit > 65535 {
		log.Fatalf("invalid '%s' '%s'. Allowed range (10, 65535]", storageLimitFlag, storageLimitDecoded)
	}
	return ParsedFlag{
		BasePath:         e2g_utils.ValidatePath(basePath, basePathFlag),
		Domain:           e2g_utils.Base64DecodeWithKill(domainTrimmed, domainFlag),
		UserName:         e2g_utils.Base64DecodeWithKill(*user, userNameFlag),
		Password:         e2g_utils.ParsePassword(*password, passwordFlag),
		IdentityFilePath: e2g_utils.ValidatePath(identityFilePath, identityFilePathFlag),
		DevelopMode:      *developMode == "true",
		StorageLimit:     uint16(limit),
	}
}
