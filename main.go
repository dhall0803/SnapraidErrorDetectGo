package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func parseEnvFile(path string) map[string]string {
	var outputMap map[string]string
	f, err := os.Open("path")
	if err != nil {
		log.Println("WARNING: Error opening .env file: \n" + err.Error())
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		envVarElements := strings.Split(scanner.Text(), "=")
		outputMap[envVarElements[0]] = envVarElements[1]
	}

	if err := scanner.Err(); err != nil {
		log.Println("Warning: Error reading from .env file: \n" + err.Error())
	}

	return outputMap
}

func loadTelegramEnvironmentVariables() (string, string) {
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

	return token, chatID
}

func sendTelegramNotification(token, chatID, message string) {
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

func isError(output string) bool {
	return !strings.Contains(output, "No error detected.")
}

func main() {
	// Telegram variables needed for sending notifications
	chatId := ""
	token := ""

	// Setup logging
	file, err := os.OpenFile("log.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	// Begin program
	log.Println("INFO: Starting program")

	// Load telegram variables
	if _, err := os.Stat(".env"); err != nil {
		log.Println("INFO: No .env file found, reading Telegram chatID and token from environment variables")
		token, chatId = loadTelegramEnvironmentVariables()
	} else {
		log.Println("INFO: .env file found, reading Telegram chatID and token from file")
		envVars := parseEnvFile(".env")
		chatId = envVars["TELEGRAM_CHATID"]
		token = envVars["TELEGRAM_TOKEN"]
	}

	if chatId == "" {
		log.Println("WARNING: Value for TELEGRAM_CHATID could not be found, notifications will not be sent")
	}

	if token == "" {
		log.Println("WARNING: Value for TELEGRAM_TOKEN could not be found, notifications will not be sent")
	}

	// Run snapraid status

	output, err := exec.Command("snapraid", "status").CombinedOutput()

	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR: %s", err))
	}

	// Send notification if there is an error
	if isError(string(output)) {
		message := "WARNING: Snapraid error detected!"
		log.Println(message)
		sendTelegramNotification(token, chatId, message)
	} else {
		log.Println("INFO: No snapraid errors detected")
	}
}
