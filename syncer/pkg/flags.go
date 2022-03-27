package pkg

import (
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

type ParsedFlag struct {
	ScreenShotInterval uint16
	FromEmail          string
	Password           string
	ToEmail            string
	KeyPath            string
	ShotsFolder        string
}

func ParseFlags() ParsedFlag {
	screenShotIntervalFlag := "cap_interval"
	fromEmailFlag := "from_email"
	gateFlag := "gate"
	toEmailFlag := "to_email"
	keyPathFlag := "keyPath"
	shotsFolderFlag := "shots"
	screenShotInterval := flag.String(screenShotIntervalFlag, "", "Capture interval in seconds. Must be greater than 10. Ex. 22")
	fromEmail := flag.String(fromEmailFlag, "", "Sender email address")
	pwd := flag.String(gateFlag, "", "Sender's password")
	toEmail := flag.String(toEmailFlag, "", "Receiver email address")
	keyPath := flag.String(keyPathFlag, "", "Recipient's public key path")
	shotsFolder := flag.String(shotsFolderFlag, "", "Shots folder path")

	flag.Parse()

	if len(strings.TrimSpace(*fromEmail)) == 0 {
		log.Fatalf("invalid '%s' '%s'", fromEmailFlag, *fromEmail)
	}
	if len(strings.TrimSpace(*toEmail)) == 0 {
		log.Fatalf("invalid '%s' '%s'", toEmailFlag, *toEmail)
	}
	if len(strings.TrimSpace(*pwd)) == 0 {
		log.Fatalf("invalid '%s' '%s'", gateFlag, *pwd)
	}
	if len(strings.TrimSpace(*shotsFolder)) == 0 {
		log.Fatalf("invalid '%s' '%s'", shotsFolderFlag, *shotsFolder)
	}
	if len(strings.TrimSpace(*keyPath)) == 0 {
		log.Fatalf("invalid '%s' '%s'", keyPathFlag, *keyPath)
	}

	interval, err := strconv.Atoi(*screenShotInterval)
	if err != nil || interval <= 10 {
		log.Fatalf("invalid '%s' '%s'", screenShotIntervalFlag, *screenShotInterval)
	}
	if _, err := os.Stat(*shotsFolder); os.IsNotExist(err) {
		log.Fatalf("invalid '%s' '%s'", screenShotIntervalFlag, *screenShotInterval)
	}

	// todo: add base64 or similar decoding
	return ParsedFlag{
		ScreenShotInterval: uint16(interval),
		FromEmail:          *fromEmail,
		Password:           *pwd,
		ToEmail:            *toEmail,
		KeyPath:            *keyPath,
		ShotsFolder:        *shotsFolder,
	}
}
