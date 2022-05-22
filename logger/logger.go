package logger

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/logging"
	mrpb "google.golang.org/genproto/googleapis/api/monitoredres"
)

type Logger interface {
	Debug(payload interface{})
	Info(payload interface{})
	Error(payload interface{})
	Emergency(payload interface{})
	Flush() error
	Close() error
}

type productionLogger struct {
	c *logging.Client
	l *logging.Logger
}

type developmentLogger struct{}

var Log Logger

func init() {
	ctx := context.Background()
	projectID := os.Getenv("GCP_PROJECT")
	functionName := os.Getenv("FUNCTION_NAME")
	region := os.Getenv("FUNCTION_REGION")

	log.Println(projectID, functionName, region)
	if projectID == "" || functionName == "" || region == "" {
		Log = &developmentLogger{}
		return
	}

	client, err := logging.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}
	logName := "cloudfunctions.googleapis.com%2Fcloud-functions"

	l := client.Logger(logName, logging.CommonResource(&mrpb.MonitoredResource{
		Type: "cloud_function",
		Labels: map[string]string{
			"function_name": functionName,
			"project_id":    projectID,
			"region":        region,
		},
	}))

	Log = &productionLogger{c: client, l: l}
}

func (pl *productionLogger) Debug(payload interface{}) {
	pl.l.Log(logging.Entry{Severity: logging.Debug, Payload: payload})
}
func (pl *productionLogger) Info(payload interface{}) {
	pl.l.Log(logging.Entry{Severity: logging.Info, Payload: payload})
}
func (pl *productionLogger) Error(payload interface{}) {
	pl.l.Log(logging.Entry{Severity: logging.Error, Payload: payload})
}
func (pl *productionLogger) Emergency(payload interface{}) {
	pl.l.Log(logging.Entry{Severity: logging.Emergency, Payload: payload})
}
func (pl *productionLogger) Flush() error {
	return pl.l.Flush()
}
func (pl *productionLogger) Close() error {
	return pl.c.Close()
}

func (pl *developmentLogger) Debug(payload interface{}) {
	log.Printf("[DEBUG]: %+v", payload)
}
func (pl *developmentLogger) Info(payload interface{}) {
	log.Printf("[INFO]: %+v", payload)
}
func (pl *developmentLogger) Error(payload interface{}) {
	log.Printf("[ERROR]: %+v", payload)
}
func (pl *developmentLogger) Emergency(payload interface{}) {
	log.Printf("[EMERGENCY]: %+v", payload)
}
func (pl *developmentLogger) Flush() error {
	return nil
}
func (pl *developmentLogger) Close() error {
	return nil
}
