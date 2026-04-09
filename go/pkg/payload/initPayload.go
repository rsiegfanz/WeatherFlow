package payload

const deviceID = "26945210-05ec-11ef-ac80-dde635ebcdb2"

type InitPayload struct {
	Commands []Command `json:"cmds"`
	AuthCmd  AuthCmd   `json:"authCmd"`
}

type AuthCmd struct {
	CmdID int    `json:"cmdId"`
	Token string `json:"token"`
}

type Command struct {
	Type      string     `json:"type"`
	Query     Query      `json:"query"`
	LatestCmd *LatestCmd `json:"latestCmd,omitempty"`
	CmdID     int        `json:"cmdId"`
}

type LatestCmd struct {
	Keys []LatestValue `json:"keys"`
}

type Query struct {
	EntityFilter EntityFilter  `json:"entityFilter"`
	PageLink     PageLink      `json:"pageLink"`
	EntityFields []EntityField `json:"entityFields,omitempty"`
	LatestValues []LatestValue `json:"latestValues,omitempty"`
	AlarmFields  []EntityField `json:"alarmFields,omitempty"`
}

type EntityFilter struct {
	Type             string        `json:"type"`
	ResolveMultiple  bool          `json:"resolveMultiple,omitempty"`
	SingleEntity     *SingleEntity `json:"singleEntity,omitempty"`
	DeviceTypes      []string      `json:"deviceTypes,omitempty"`
	DeviceNameFilter string        `json:"deviceNameFilter,omitempty"`
}

type SingleEntity struct {
	ID         string `json:"id"`
	EntityType string `json:"entityType"`
}

type PageLink struct {
	PageSize               int         `json:"pageSize"`
	Page                   int         `json:"page"`
	SortOrder              SortOrder   `json:"sortOrder"`
	TextSearch             interface{} `json:"textSearch,omitempty"`
	TypeList               []string    `json:"typeList,omitempty"`
	SeverityList           []string    `json:"severityList,omitempty"`
	StatusList             []string    `json:"statusList,omitempty"`
	SearchPropagatedAlarms bool        `json:"searchPropagatedAlarms,omitempty"`
	AssigneeID             interface{} `json:"assigneeId,omitempty"`
	TimeWindow             int         `json:"timeWindow,omitempty"`
	Dynamic                bool        `json:"dynamic,omitempty"`
}

type SortOrder struct {
	Key struct {
		Type string `json:"type"`
		Key  string `json:"key"`
	} `json:"key"`
	Direction string `json:"direction"`
}

type EntityField struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

type LatestValue struct {
	Type string `json:"type"`
	Key  string `json:"key"`
}

func PrepareInitPayload(token string) InitPayload {
	return InitPayload{
		AuthCmd:  AuthCmd{Token: token},
		Commands: prepareCommands(),
	}
}

func prepareCommands() []Command {
	defaultEntityFields := []EntityField{
		{Type: "ENTITY_FIELD", Key: "name"},
		{Type: "ENTITY_FIELD", Key: "label"},
		{Type: "ENTITY_FIELD", Key: "additionalInfo"},
	}

	deviceFilter := EntityFilter{
		Type: "singleEntity",
		SingleEntity: &SingleEntity{
			ID:         deviceID,
			EntityType: "DEVICE",
		},
	}

	createdTimeSortOrder := SortOrder{
		Key: struct {
			Type string `json:"type"`
			Key  string `json:"key"`
		}{Type: "ENTITY_FIELD", Key: "createdTime"},
		Direction: "DESC",
	}

	cmds := []Command{
		// Command 1: All weather telemetry in one subscription
		{Type: "ENTITY_DATA", CmdID: 1, Query: Query{
			EntityFilter: deviceFilter,
			PageLink:     PageLink{PageSize: 1, SortOrder: createdTimeSortOrder},
			EntityFields: defaultEntityFields,
			LatestValues: []LatestValue{
				{Type: "TIME_SERIES", Key: "airTemperature"},
				{Type: "TIME_SERIES", Key: "airHumidity"},
				{Type: "TIME_SERIES", Key: "barometricPressure"},
				{Type: "TIME_SERIES", Key: "windSpeed"},
				{Type: "TIME_SERIES", Key: "windDirectionSensor"},
				{Type: "TIME_SERIES", Key: "rainGauge"},
				{Type: "TIME_SERIES", Key: "uvIndex"},
				{Type: "TIME_SERIES", Key: "lightIntensity"},
				{Type: "TIME_SERIES", Key: "battery"},
			},
		}},
		// Command 2: Device attributes (location, firmware, etc.)
		{Type: "ENTITY_DATA", CmdID: 2, Query: Query{
			EntityFilter: deviceFilter,
			PageLink:     PageLink{PageSize: 1, SortOrder: createdTimeSortOrder},
			EntityFields: defaultEntityFields,
			LatestValues: []LatestValue{
				{Type: "ATTRIBUTE", Key: "latitude"},
				{Type: "ATTRIBUTE", Key: "longitude"},
				{Type: "ATTRIBUTE", Key: "altitude"},
				{Type: "ATTRIBUTE", Key: "firmwareVersion"},
				{Type: "ATTRIBUTE", Key: "hardwareVersion"},
				{Type: "ATTRIBUTE", Key: "active"},
			},
		}},
	}

	cmds = append(cmds,
		// Command 10: Active alarms for water level sensors
		Command{Type: "ALARM_DATA", CmdID: 10, Query: Query{
			EntityFilter: EntityFilter{
				Type:             "deviceType",
				ResolveMultiple:  true,
				DeviceTypes:      []string{"Dragino LDDS Water Level"},
				DeviceNameFilter: "",
			},
			PageLink: PageLink{
				PageSize: 1024,
				SortOrder: SortOrder{
					Key: struct {
						Type string `json:"type"`
						Key  string `json:"key"`
					}{Type: "ALARM_FIELD", Key: "createdTime"},
					Direction: "DESC",
				},
				TypeList:               []string{},
				SeverityList:           []string{},
				StatusList:             []string{"ACTIVE"},
				SearchPropagatedAlarms: false,
				TimeWindow:             604800000,
			},
			AlarmFields: []EntityField{
				{Type: "ALARM_FIELD", Key: "originatorLabel"},
				{Type: "ALARM_FIELD", Key: "createdTime"},
				{Type: "ALARM_FIELD", Key: "type"},
				{Type: "ALARM_FIELD", Key: "severity"},
			},
			EntityFields: []EntityField{},
			LatestValues: []LatestValue{},
		}},
		// Command 11: Water level devices
		Command{Type: "ENTITY_DATA", CmdID: 11,
			LatestCmd: &LatestCmd{Keys: []LatestValue{
				{Type: "ATTRIBUTE", Key: "displayName"},
				{Type: "TIME_SERIES", Key: "waterLevel"},
			}},
			Query: Query{
				EntityFilter: EntityFilter{
					Type:             "deviceType",
					ResolveMultiple:  true,
					DeviceTypes:      []string{"Dragino LDDS Water Level"},
					DeviceNameFilter: "",
				},
				PageLink: PageLink{
					PageSize: 1024,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{Type: "ATTRIBUTE", Key: "displayName"},
						Direction: "ASC",
					},
					Dynamic: true,
				},
				EntityFields: defaultEntityFields,
				LatestValues: []LatestValue{
					{Type: "ATTRIBUTE", Key: "displayName"},
					{Type: "TIME_SERIES", Key: "waterLevel"},
				},
			},
		},
	)

	return cmds
}
