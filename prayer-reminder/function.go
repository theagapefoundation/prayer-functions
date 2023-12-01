// Package p contains an HTTP Cloud Function.
package prayerreminders

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func init() {
  functions.HTTP("NotifyUserPrayerReminders", notifyUserPrayerReminders)
}

func splitIntoBatches(input []string, batchSize int) [][]string {
  var batches [][]string
  
  for start := 0; start < len(input); start += batchSize {
    end := start + batchSize
    if end > len(input) {
      end = len(input)
    }
    batches = append(batches, input[start:end])
  }

  return batches
}

func notifyUserPrayerReminders(w http.ResponseWriter, r *http.Request) {
  conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
  if err != nil {
    log.Fatalf("Error while connecting to a database: %v\n", err)
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }
  defer conn.Close(context.Background())
  firebaseCtx := context.Background()
  firebaseApp, err := firebase.NewApp(firebaseCtx, nil)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    log.Panicf("Error while initializing a firebase: %v\n", err)
    return
  }
  messagingClient, err := firebaseApp.Messaging(firebaseCtx)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    log.Fatalf("Error while getting messaging client: %v\n", err)
    return
  }
  
  rows, err := conn.Query(context.Background(), `
    SELECT ARRAY_AGG(user_fcm_tokens.value), reminders.value, corporate_prayers.id, corporate_prayers.title
    FROM user_fcm_tokens
    INNER JOIN prayers ON prayers.user_id = user_fcm_tokens.user_id
    INNER JOIN corporate_prayers ON corporate_prayers.id = prayers.corporate_id
    INNER JOIN reminders ON reminders.id = corporate_prayers.reminder_id
    WHERE (
      corporate_prayers.started_at IS NULL
      OR DATE_TRUNC(
        'day', 
        NOW() AT TIME ZONE 'UTC' 
        + INTERVAL '1 hour' * EXTRACT(TIMEZONE_HOUR FROM reminders.time)
        + INTERVAL '1 minute' * EXTRACT(TIMEZONE_MINUTE FROM reminders.time)
      ) >= DATE_TRUNC('day', corporate_prayers.started_at)
    ) AND (
      corporate_prayers.ended_at IS NULL
      OR DATE_TRUNC(
        'day', 
        NOW() AT TIME ZONE 'UTC' 
        + INTERVAL '1 hour' * EXTRACT(TIMEZONE_HOUR FROM reminders.time)
        + INTERVAL '1 minute' * EXTRACT(TIMEZONE_MINUTE FROM reminders.time)
      ) <= DATE_TRUNC('day', corporate_prayers.ended_at)
    ) AND POSITION(
      EXTRACT(DOW FROM DATE_TRUNC(
          'day', 
          NOW() AT TIME ZONE 'UTC' 
          + INTERVAL '1 hour' * EXTRACT(TIMEZONE_HOUR FROM reminders.time)
          + INTERVAL '1 minute' * EXTRACT(TIMEZONE_MINUTE FROM reminders.time)
      ))::text in reminders.days
    ) > 0 AND 
    DATE_TRUNC('minute', reminders.time::time) = DATE_TRUNC(
      'minute', 
      (NOW() AT TIME ZONE 'UTC' 
      + INTERVAL '1 hour' * EXTRACT(TIMEZONE_HOUR FROM reminders.time)
      + INTERVAL '1 minute' * EXTRACT(TIMEZONE_MINUTE FROM reminders.time))::time
    )
    GROUP BY reminders.id, corporate_prayers.id
  `)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    log.Panicf("Error while running query: %v\n", err)
    return
  }
  var tokens pgtype.Array[string]
  var message string
  var corporateId string
  var title string
  _, err = pgx.ForEachRow(rows, []any{&tokens, &message, &corporateId, &title}, func() error {
    for _, partiion := range splitIntoBatches(tokens.Elements, 500) {
      message := &messaging.MulticastMessage{
        Notification: &messaging.Notification{
          Title: title,
          Body: message,
        },
        Data: map[string]string{
          corporateId: corporateId,
        },
        Tokens: partiion,
      }
      _, err = messagingClient.SendEachForMulticast(context.Background(), message)
      if err != nil {
        log.Printf("Error while running query: %v\n", err)
      }
    }
    return nil
  })
  if err != nil {
    log.Printf("Error while iterating rows")
  }
  fmt.Fprintf(w, "success\n")
}