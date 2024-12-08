package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

func timestampLog(logFile *os.File, prefix string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s %v", timestamp, prefix, v)
	log.Println(logMessage)

	if logFile != nil {
		logFile.WriteString(logMessage + "\n")
		logFile.Sync()
	}
}

func main() {
	url := "wss://thingsboard.bda-itnovum.com/api/ws"

	logDir := "logs"
	os.MkdirAll(logDir, os.ModePerm)

	generalLogFile, err := os.Create(filepath.Join(logDir, "general.log"))
	if err != nil {
		log.Fatal("Fehler beim Erstellen der Allgemeinen Log-Datei:", err)
	}
	defer generalLogFile.Close()

	messageLogFile, err := os.Create(filepath.Join(logDir, "messages.log"))
	if err != nil {
		log.Fatal("Fehler beim Erstellen der Nachrichten Log-Datei:", err)
	}
	defer messageLogFile.Close()

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		timestampLog(generalLogFile, "FEHLER:", "Verbindungsfehler", err)
	}
	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	readData(conn, done, generalLogFile, messageLogFile)
	sendInitData(conn, generalLogFile)

	// go func() {
	// 	for {
	// 		message := []byte("Periodische Aktualisierung")
	// 		err := conn.WriteMessage(websocket.TextMessage, message)
	// 		if err != nil {
	// 			log.Println("Sendefehler:", err)
	// 			return
	// 		}
	// 		// time.Sleep(5 * time.Second)
	// 	}
	// }()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			timestampLog(generalLogFile, "INFO:", "Verbindung wird geschlossen")

			err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				timestampLog(generalLogFile, "FEHLER:", "Fehler beim Schließen", err)
				return
			}
			return
		}
	}
}

func sendInitData(conn *websocket.Conn, generalLogFile *os.File) {
	jsonData := `{"cmds":[{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"airTemperature"}]},"cmdId":1},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"windSpeed"}]},"cmdId":2},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"rainGauge"}]},"cmdId":3},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":4},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":5},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":6},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":7},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":8},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"windDirectionSensor"},{"type":"TIME_SERIES","key":"windSpeed"}]},"cmdId":9},{"type":"ALARM_DATA","query":{"entityFilter":{"type":"deviceType","resolveMultiple":true,"deviceTypes":["Dragino LDDS Water Level"],"deviceNameFilter":""},"pageLink":{"page":0,"pageSize":1024,"textSearch":null,"typeList":[],"severityList":[],"statusList":["ACTIVE"],"searchPropagatedAlarms":false,"assigneeId":null,"sortOrder":{"key":{"key":"createdTime","type":"ALARM_FIELD"},"direction":"DESC"},"timeWindow":604800000},"alarmFields":[{"type":"ALARM_FIELD","key":"originatorLabel"},{"type":"ALARM_FIELD","key":"createdTime"},{"type":"ALARM_FIELD","key":"type"},{"type":"ALARM_FIELD","key":"severity"}],"entityFields":[],"latestValues":[]},"cmdId":10},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"deviceType","resolveMultiple":true,"deviceTypes":["Dragino LDDS Water Level"],"deviceNameFilter":""},"pageLink":{"page":0,"pageSize":1024,"textSearch":null,"dynamic":true,"sortOrder":{"key":{"key":"displayName","type":"ATTRIBUTE"},"direction":"ASC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"ATTRIBUTE","key":"displayName"},{"type":"TIME_SERIES","key":"waterLevel"}]},"latestCmd":{"keys":[{"type":"ATTRIBUTE","key":"displayName"},{"type":"TIME_SERIES","key":"waterLevel"}]},"cmdId":11}],"authCmd":{"cmdId":0,"token":"eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJkNThiMThhMC0xNDQwLTExZWYtYWVmNC1hZjI4M2U1MDk0ZDkiLCJ1c2VySWQiOiIxMzgxNDAwMC0xZGQyLTExYjItODA4MC04MDgwODA4MDgwODAiLCJzY29wZXMiOlsiQ1VTVE9NRVJfVVNFUiJdLCJzZXNzaW9uSWQiOiI4MTRiZWE4ZC1lMWUzLTQ3YWYtYTFjOS04MTE0MDE1NTliOTUiLCJleHAiOjE3MzM2ODM2MDQsImlzcyI6InRoaW5nc2JvYXJkLmlvIiwiaWF0IjoxNzMzNjc0NjA0LCJmaXJzdE5hbWUiOiJQdWJsaWMiLCJsYXN0TmFtZSI6IlB1YmxpYyIsImVuYWJsZWQiOnRydWUsImlzUHVibGljIjp0cnVlLCJ0ZW5hbnRJZCI6ImNmYTE4ZjYwLWI5MDktMTFlZS1hYTE0LTYxMDRhYWQwMjBlYyIsImN1c3RvbWVySWQiOiJkNThiMThhMC0xNDQwLTExZWYtYWVmNC1hZjI4M2U1MDk0ZDkifQ.GsrVSShKKkv0GhI2HvzJdNHTvDw60XADcWT-x-C4MXXViGX1kp_u_fzVAgwrJGyaGInc_VFVA5cL8j0Qmexitw"}}`

	err := conn.WriteMessage(websocket.TextMessage, []byte(jsonData))
	if err != nil {
		timestampLog(generalLogFile, "FEHLER:", "Sendefehler", err)
		return
	}
	timestampLog(generalLogFile, "INFO:", "Initiale Nachricht gesendet")
}

func readData(conn *websocket.Conn, done chan struct{}, generalLogFile *os.File, messageLogFile *os.File) {
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				timestampLog(generalLogFile, "FEHLER:", "Lesefehler", err)
				return
			}
			timestampLog(messageLogFile, "NACHRICHT:", string(message))
		}
	}()
}
