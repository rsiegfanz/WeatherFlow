package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/gorilla/websocket"
)

const url = "wss://thingsboard.bda-itnovum.com/api/ws"

type Client struct {
	Token string
}

func NewClient(token string) (*Client, error) {
	return &Client{
		Token: token,
	}, nil
}

func (c *Client) Close() {}

func (c *Client) Connect() error {

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	payload := `{"cmds":[{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"airTemperature"}]},"cmdId":1},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"windSpeed"}]},"cmdId":2},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"rainGauge"}]},"cmdId":3},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":4},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":5},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":6},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":7},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":8},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"windDirectionSensor"},{"type":"TIME_SERIES","key":"windSpeed"}]},"cmdId":9},{"type":"ALARM_DATA","query":{"entityFilter":{"type":"deviceType","resolveMultiple":true,"deviceTypes":["Dragino LDDS Water Level"],"deviceNameFilter":""},"pageLink":{"page":0,"pageSize":1024,"textSearch":null,"typeList":[],"severityList":[],"statusList":["ACTIVE"],"searchPropagatedAlarms":false,"assigneeId":null,"sortOrder":{"key":{"key":"createdTime","type":"ALARM_FIELD"},"direction":"DESC"},"timeWindow":604800000},"alarmFields":[{"type":"ALARM_FIELD","key":"originatorLabel"},{"type":"ALARM_FIELD","key":"createdTime"},{"type":"ALARM_FIELD","key":"type"},{"type":"ALARM_FIELD","key":"severity"}],"entityFields":[],"latestValues":[]},"cmdId":10},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"deviceType","resolveMultiple":true,"deviceTypes":["Dragino LDDS Water Level"],"deviceNameFilter":""},"pageLink":{"page":0,"pageSize":1024,"textSearch":null,"dynamic":true,"sortOrder":{"key":{"key":"displayName","type":"ATTRIBUTE"},"direction":"ASC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"ATTRIBUTE","key":"displayName"},{"type":"TIME_SERIES","key":"waterLevel"}]},"latestCmd":{"keys":[{"type":"ATTRIBUTE","key":"displayName"},{"type":"TIME_SERIES","key":"waterLevel"}]},"cmdId":11}],"authCmd":{"cmdId":0,"token":"XXTOKENXX"}}`

	log.Println(c.Token)
	strings.ReplaceAll(payload, "XXTOKENXX", c.Token)

	// payload := `{"cmds":[{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"airTemperature"}]},"cmdId":1},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"windSpeed"}]},"cmdId":2},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"rainGauge"}]},"cmdId":3},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":4},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":5},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":6},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":7},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1024,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[]},"cmdId":8},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"singleEntity","singleEntity":{"id":"26945210-05ec-11ef-ac80-dde635ebcdb2","entityType":"DEVICE"}},"pageLink":{"pageSize":1,"page":0,"sortOrder":{"key":{"type":"ENTITY_FIELD","key":"createdTime"},"direction":"DESC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"TIME_SERIES","key":"windDirectionSensor"},{"type":"TIME_SERIES","key":"windSpeed"}]},"cmdId":9},{"type":"ALARM_DATA","query":{"entityFilter":{"type":"deviceType","resolveMultiple":true,"deviceTypes":["Dragino LDDS Water Level"],"deviceNameFilter":""},"pageLink":{"page":0,"pageSize":1024,"textSearch":null,"typeList":[],"severityList":[],"statusList":["ACTIVE"],"searchPropagatedAlarms":false,"assigneeId":null,"sortOrder":{"key":{"key":"createdTime","type":"ALARM_FIELD"},"direction":"DESC"},"timeWindow":604800000},"alarmFields":[{"type":"ALARM_FIELD","key":"originatorLabel"},{"type":"ALARM_FIELD","key":"createdTime"},{"type":"ALARM_FIELD","key":"type"},{"type":"ALARM_FIELD","key":"severity"}],"entityFields":[],"latestValues":[]},"cmdId":10},{"type":"ENTITY_DATA","query":{"entityFilter":{"type":"deviceType","resolveMultiple":true,"deviceTypes":["Dragino LDDS Water Level"],"deviceNameFilter":""},"pageLink":{"page":0,"pageSize":1024,"textSearch":null,"dynamic":true,"sortOrder":{"key":{"key":"displayName","type":"ATTRIBUTE"},"direction":"ASC"}},"entityFields":[{"type":"ENTITY_FIELD","key":"name"},{"type":"ENTITY_FIELD","key":"label"},{"type":"ENTITY_FIELD","key":"additionalInfo"}],"latestValues":[{"type":"ATTRIBUTE","key":"displayName"},{"type":"TIME_SERIES","key":"waterLevel"}]},"latestCmd":{"keys":[{"type":"ATTRIBUTE","key":"displayName"},{"type":"TIME_SERIES","key":"waterLevel"}]},"cmdId":11}],"authCmd":{"cmdId":0,"token":"eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJkNThiMThhMC0xNDQwLTExZWYtYWVmNC1hZjI4M2U1MDk0ZDkiLCJ1c2VySWQiOiIxMzgxNDAwMC0xZGQyLTExYjItODA4MC04MDgwODA4MDgwODAiLCJzY29wZXMiOlsiQ1VTVE9NRVJfVVNFUiJdLCJzZXNzaW9uSWQiOiI4MTRiZWE4ZC1lMWUzLTQ3YWYtYTFjOS04MTE0MDE1NTliOTUiLCJleHAiOjE3MzM2ODM2MDQsImlzcyI6InRoaW5nc2JvYXJkLmlvIiwiaWF0IjoxNzMzNjc0NjA0LCJmaXJzdE5hbWUiOiJQdWJsaWMiLCJsYXN0TmFtZSI6IlB1YmxpYyIsImVuYWJsZWQiOnRydWUsImlzUHVibGljIjp0cnVlLCJ0ZW5hbnRJZCI6ImNmYTE4ZjYwLWI5MDktMTFlZS1hYTE0LTYxMDRhYWQwMjBlYyIsImN1c3RvbWVySWQiOiJkNThiMThhMC0xNDQwLTExZWYtYWVmNC1hZjI4M2U1MDk0ZDkifQ.GsrVSShKKkv0GhI2HvzJdNHTvDw60XADcWT-x-C4MXXViGX1kp_u_fzVAgwrJGyaGInc_VFVA5cL8j0Qmexitw"}}`

	// jsonData, err := json.Marshal(payload)
	// if err != nil {
	// 	return fmt.Errorf("JSON conversion error: %v", err)
	// }

	err = conn.WriteMessage(websocket.TextMessage, []byte(payload))
	if err != nil {
		return fmt.Errorf("send error: %v", err)
	}

	log.Println("Initial message sent")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})
	go c.readData(conn, done)

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Println("Closing connection")
			err := conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				return fmt.Errorf("error closing: %v", err)
			}
			return nil
		}
	}
}

func (c *Client) readData(conn *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(fmt.Errorf("Read error", err))
			return
		}
		log.Println("New message")
		log.Println(message)
	}
}
