package pkg

import (
	"flag"
	"github.com/mrnakumar/e2g_utils"
	"log"
	"os"
	"strconv"
	"strings"
)

type ParsedFlag struct {
	ScreenShotInterval uint16
	SyncInterval       uint16
	ServerUrl          string
	UserName           string
	Password           string
	HeaderKeyPath      string
	ShotKeyPath        string
	ShotsFolder        string
	StorageLimit       uint16
}

type flagInfo struct {
	name         string
	defaultValue string
	usage        string
	userSupplied *string
	validation   func(info *flagInfo)
	provide      func(info *flagInfo)
}

func ParseFlags() ParsedFlag {
	var screenShotInterval uint16 = 0
	var syncInterval uint16
	var storageLimit uint16
	var serverUrl string
	var userName string
	var password string
	var HeaderKeyPath string
	var ShotKeyPath string
	var shotsFolder string
	flagInfos := []*flagInfo{
		{
			name:         "cap_interval",
			usage:        "Capture interval in seconds. Must be greater than 10. Ex. 22",
			userSupplied: nil,
			validation:   checkUnsigned16,
			provide:      asUin16(&screenShotInterval),
		},
		{
			name:         "sync_interval",
			usage:        "Sync interval in seconds. Must be greater than 10.",
			userSupplied: nil,
			validation:   checkUnsigned16,
			provide:      asUin16(&syncInterval),
		},
		{
			name:         "storage_limit",
			usage:        "Screenshot storage limit in Megabyte.",
			userSupplied: nil,
			validation:   checkUnsigned16,
			provide:      asUin16(&storageLimit),
		},
		{
			name:         "user_name",
			usage:        "Sender user name (as expected by server)",
			userSupplied: nil,
			validation:   emptyCheck,
			provide:      asBase64Decode(&userName),
		},
		{
			name:         "server_url",
			usage:        "Server url to post shots to",
			userSupplied: nil,
			validation:   emptyCheck,
			provide:      asBase64Decode(&serverUrl),
		},
		{
			name:         "gate",
			usage:        "Sender's password",
			userSupplied: nil,
			validation:   emptyCheck,
			provide:      asPassword(&password),
		},
		{
			name:         "shot_key_path",
			usage:        "Recipient's public key path",
			userSupplied: nil,
			validation:   checkPathExists,
			provide:      asBase64Decode(&ShotKeyPath),
		},
		{
			name:         "header_key_path",
			usage:        "Recipient's public key path",
			userSupplied: nil,
			validation:   checkPathExists,
			provide:      asBase64Decode(&HeaderKeyPath),
		},
		{
			name:         "shots",
			usage:        "Shots folder path",
			userSupplied: nil,
			validation:   checkPathExists,
			provide:      asBase64Decode(&shotsFolder),
		},
	}
	for _, flagInfo := range flagInfos {
		flagInfo.userSupplied = flag.String(flagInfo.name, flagInfo.defaultValue, flagInfo.usage)
	}
	flag.Parse()
	validateAndProvide(flagInfos)

	return ParsedFlag{
		ScreenShotInterval: screenShotInterval,
		SyncInterval:       syncInterval,
		UserName:           userName,
		ServerUrl:          serverUrl,
		Password:           password,
		HeaderKeyPath:      HeaderKeyPath,
		ShotKeyPath:        ShotKeyPath,
		ShotsFolder:        shotsFolder,
		StorageLimit:       storageLimit,
	}
}

func validateAndProvide(flags []*flagInfo) {
	for _, flagInfo := range flags {
		flagInfo.validation(flagInfo)
		flagInfo.provide(flagInfo)
	}
}
func emptyCheck(input *flagInfo) {
	if len(strings.TrimSpace(*input.userSupplied)) == 0 {
		log.Fatalf("invalid '%s' '%s'. Must not be empty", input.name, *input.userSupplied)
	}
}

func checkUnsigned16(input *flagInfo) {
	interval, err := strconv.Atoi(*input.userSupplied)
	if err != nil || interval <= 10 || interval > 65535 {
		log.Fatalf("invalid '%s' '%s'. Allowed range (10, 65535]", input.name, *input.userSupplied)
	}
}

func checkPathExists(input *flagInfo) {
	decoded := e2g_utils.Base64DecodeWithKill(*input.userSupplied, input.name)
	if _, err := os.Stat(decoded); os.IsNotExist(err) {
		log.Fatalf("invalid '%s' '%s'. Path does not exist.", input.name, decoded)
	}
}

func asUin16(target *uint16) func(*flagInfo) {
	return func(input *flagInfo) {
		interval, _ := strconv.Atoi(*input.userSupplied)
		*target = uint16(interval)
	}
}

func asBase64Decode(target *string) func(*flagInfo) {
	return func(input *flagInfo) {
		decoded, err := e2g_utils.Base64Decode(*input.userSupplied)
		if err != nil {
			log.Fatalf("failed to decode '%s'", input.name)
		}
		*target = decoded
	}
}

func asPassword(target *string) func(info *flagInfo) {
	return func(input *flagInfo) {
		*target = e2g_utils.ParsePassword(*input.userSupplied, input.name)
	}
}
