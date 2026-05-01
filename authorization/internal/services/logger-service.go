package services

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

func NewLoggger(name string, logger **slog.Logger) {
	err := os.MkdirAll("/logs", os.ModePerm)
	if err != nil {
		fmt.Println("[-] logger.mkdir:", err.Error())
		return
	}

	filename := fmt.Sprintf("/logs/%s-%s.log", name, time.Now().Format(time.DateOnly))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("[-] logger.open:", err.Error())
		return
	}

	fmt.Println("[+] logger.new:", name)
	*logger = slog.New(slog.NewTextHandler(
		file, nil,
	))
}
