package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/linzb0123/chatgpt-demo/chatgpt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	apiKey = ""

	logDir = "output/log"
)

func main() {
	initLog()

	flag.StringVar(&apiKey, "apiKey", "", "chatgpt api key")
	flag.Parse()

	if len(apiKey) == 0 {
		panic("api key invalid")
	}

	log.Printf("apiKey:%s", apiKey)
	chat := chatgpt.New(apiKey, "")
	color.Yellow("init chat-get success,model(%s);you can input content to chat", chat.GetModel())

	scanner := bufio.NewScanner(os.Stdin)

	var input string
	for {
		if scanner.Scan() {
			input = scanner.Text()
			switch input {
			case "/reset":
				chat.Reset()
			case "/exit":
				os.Exit(0)
			default:
				reply, err := chat.Chat(input)
				if err != nil {
					color.Red("**Error:%s", err.Error())
					continue
				}

				color.Green("%s", reply)
				color.Blue("===================================================================")
				color.Blue("===================================================================")
			}

		}

	}
}

func initLog() {
	logFileName := fmt.Sprintf("chat-gpt-%s.log", time.Now().Format("20060102"))
	logFilePath := filepath.Join(logDir, logFileName)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			log.Fatalf("init log config.mkdir err:%v", err)
		}
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("open log file err:%v", err)
	}

	absLogPath, err := filepath.Abs(logFilePath)
	if err != nil {
		log.Fatalf("get abs log file path err:%v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	color.Yellow("init log success.log path:%s", absLogPath)

}
