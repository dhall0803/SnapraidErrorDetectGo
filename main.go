package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func sendTelegramNotification(message string) {
	token := os.Getenv("TELEGRAM_TOKEN")
	chatID := os.Getenv("TELEGRAM_CHATID")
	if token == "" {
		log.Println("WARNING: Could not send telegram notification")
		log.Println("WARNING: Could not find token, have you set it in the \"TELEGRAM_TOKEN\" environment variable?")
	}
	if chatID == "" {
		log.Println("WARNING: Could not send telegram notification")
		log.Println("WARNING: Could not find chatid, have you set it in the \"CHAT_ID\" environment variable?")
	}
	if chatID != "" && token != "" {
		log.Println("INFO: Sending notification...")
		resp, err := http.Get("https://api.telegram.org/bot" + token + "/sendMessage?chat_id=" + chatID + "&text=" + message)
		if err != nil {
			log.Println("WARNING: Error sending telegram notification")
			log.Println("WARNGING: " + err.Error())
		} else {
			log.Println("INFO: Notification sent")
			log.Println("INFO: HTTP status code" + resp.Status)
		}
	}
}

func main() {
	// Setup logging
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// Begin program
	log.Println("INFO: Starting program")

	// Run snapraid status

	output, err := exec.Command("snapraid", "status").CombinedOutput()

	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR: %s", err))
	}

	if !strings.Contains(string(output), "No error detected.") {
		message := "WARNING: Snapraid error detected!"
		log.Println(message)
		sendTelegramNotification(message)
	} else {
		log.Println("INFO: No snapraid errors detected")
	}
}
