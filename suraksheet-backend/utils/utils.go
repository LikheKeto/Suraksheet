package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LikheKeto/Suraksheet/config"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
)

var Validate = validator.New()

type ExtractionArgs struct {
	DocID     int    `json:"documentID"`
	UserID    int    `json:"userID"`
	FileKey   string `json:"fileKey"`
	Extension string `json:"extension"`
	Language  string `json:"language"`
}

func QueueForExtraction(ch *amqp.Channel, q amqp.Queue, args ExtractionArgs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	jsonb, err := json.Marshal(map[string]any{
		"documentID": args.DocID,
		"fileKey":    args.FileKey,
		"bucket":     config.Envs.MinioBucketName,
		"extension":  args.Extension,
		"language":   args.Language,
		"userID":     args.UserID,
	})
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         jsonb,
	})
	return err
}

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func HashString(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}
