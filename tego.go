package tego

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func SendMsg(msg string) error {
	token := os.Getenv("TEGO_TELEGRAM_BOT_TOKEN")
	if token == "" {
		return errors.New("env TEGO_TELEGRAM_BOT_TOKEN empty")
	}

	chatID := os.Getenv("TEGO_TELEGRAM_CHAT_ID")
	if chatID == "" {
		return errors.New("env TEGO_TELEGRAM_CHAT_ID empty")
	}

	topicID := os.Getenv("TEGO_TELEGRAM_TOPIC_ID")
	if topicID == "" {
		return errors.New("env TEGO_TELEGRAM_TOPIC_ID empty")
	}

	body := map[string]interface{}{
		"message_thread_id": topicID,
		"chat_id":           chatID,
		"text":              msg,
	}

	jsonByte, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error json marshal body: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonByte))
	if err != nil {
		return fmt.Errorf("error create new request: %w", err)
	}
	httpReq.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	httpRes, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sent request: %w", err)
	}
	defer func() {
		err := httpRes.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("error read response body: %w", err)
	}

	if httpRes.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code %d", httpRes.StatusCode)
	}

	return nil
}

func SendJSON(caption string, jsonMap map[string]interface{}) error {
	if len(jsonMap) == 0 {
		return errors.New("jsonmap empty")
	}

	token := os.Getenv("TEGO_TELEGRAM_BOT_TOKEN")
	if token == "" {
		return errors.New("env TEGO_TELEGRAM_BOT_TOKEN empty")
	}

	chatID := os.Getenv("TEGO_TELEGRAM_CHAT_ID")
	if chatID == "" {
		return errors.New("env TEGO_TELEGRAM_CHAT_ID empty")
	}

	topicID := os.Getenv("TEGO_TELEGRAM_TOPIC_ID")
	if topicID == "" {
		return errors.New("env TEGO_TELEGRAM_TOPIC_ID empty")
	}

	jsonByte, err := json.Marshal(jsonMap)
	if err != nil {
		return fmt.Errorf("error json marshal jsonMap: %w", err)
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	err = writer.WriteField("message_thread_id", topicID)
	if err != nil {
		return fmt.Errorf("error write field message_thread_id: %w", err)
	}

	err = writer.WriteField("chat_id", chatID)
	if err != nil {
		return fmt.Errorf("error write field chat_id: %w", err)
	}

	err = writer.WriteField("caption", caption)
	if err != nil {
		return fmt.Errorf("error write field caption: %w", err)
	}

	part, err := writer.CreateFormFile("document", "data.json")
	if err != nil {
		return fmt.Errorf("error write field caption: %w", err)
	}

	_, err = part.Write(jsonByte)
	if err != nil {
		return fmt.Errorf("error write: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error close writer: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token)

	httpReq, err := http.NewRequest(http.MethodPost, url, &buffer)
	if err != nil {
		return fmt.Errorf("error create new request: %w", err)
	}
	httpReq.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	httpRes, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error sent request: %w", err)
	}
	defer func() {
		err := httpRes.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.ReadAll(httpRes.Body)
	if err != nil {
		return fmt.Errorf("error read response body: %w", err)
	}

	if httpRes.StatusCode != http.StatusOK {
		return fmt.Errorf("response status code %d", httpRes.StatusCode)
	}

	return nil
}
