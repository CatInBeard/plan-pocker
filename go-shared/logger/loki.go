package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/google/uuid"
)

func SendToLoki(entry LogEntry) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	errorCode := extractErrorCode(entry.Header)

	additionalInfo := fmt.Sprintf("File: %s, line: %d", entry.File, entry.Line)

	streams := map[string]string{
		"container": entry.Container,
		"level":     fmt.Sprintf("%s", entry.Level),
		"host":      hostname,
		"errorCode": errorCode,
		"file":      entry.File,
	}

	entry.AdditionalLabels = append(entry.AdditionalLabels, Label{Key: "LogUID", Value: uuid.New().String()})

	for _, label := range entry.AdditionalLabels {
		streams[label.Key] = label.Value
	}

	logData := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": streams,
				"values": [][]interface{}{
					{
						fmt.Sprintf("%d", entry.Time.UnixNano()),
						entry.Header + ": " + entry.Body + "\nAdditionalInfo: " + additionalInfo,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", lokiURL+"/loki/api/v1/push", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return errors.New("failed to send log to Loki, status code: " + resp.Status)
	}

	return nil
}

func extractErrorCode(header string) string {

	re := regexp.MustCompile(`\[(.*)\]`)
	matches := re.FindStringSubmatch(header)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
