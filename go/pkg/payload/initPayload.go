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
	Type         string        `json:"type"`
	Query        Query         `json:"query"`
	EntityFields []EntityField `json:"entityFields,omitempty"`
	LatestValues []LatestValue `json:"latestValues,omitempty"`
	CmdID        int           `json:"cmdId"`
}

type Query struct {
	EntityFilter EntityFilter `json:"entityFilter"`
	PageLink     PageLink     `json:"pageLink"`
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
	payload := InitPayload{
		AuthCmd: AuthCmd{
			CmdID: 0,
			Token: token,
		},
		Commands: prepareCommands(),
	}

	return payload
}

func prepareCommands() []Command {
	return []Command{
		// Command 1: Air Temperature
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{
				{Type: "TIME_SERIES", Key: "airTemperature"},
			},
			CmdID: 1,
		},
		// Command 2: Wind Speed
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{
				{Type: "TIME_SERIES", Key: "windSpeed"},
			},
			CmdID: 2,
		},
		// Command 3: Rain Gauge
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{
				{Type: "TIME_SERIES", Key: "rainGauge"},
			},
			CmdID: 3,
		},
		// Command 4-8: Empty Latest Values
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1024,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{},
			CmdID:        4,
		},
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1024,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{},
			CmdID:        5,
		},
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1024,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{},
			CmdID:        6,
		},
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1024,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{},
			CmdID:        7,
		},
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1024,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{},
			CmdID:        8,
		},
		// Command 9: Wind Direction and Wind Speed
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type: "singleEntity",
					SingleEntity: &SingleEntity{
						ID:         generateDeviceId(),
						EntityType: "DEVICE",
					},
				},
				PageLink: PageLink{
					PageSize: 1,
					Page:     0,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ENTITY_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{
				{Type: "TIME_SERIES", Key: "windDirectionSensor"},
				{Type: "TIME_SERIES", Key: "windSpeed"},
			},
			CmdID: 9,
		},
		// Command 10: Alarm Data
		{
			Type: "ALARM_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type:             "deviceType",
					ResolveMultiple:  true,
					DeviceTypes:      []string{"Dragino LDDS Water Level"},
					DeviceNameFilter: "",
				},
				PageLink: PageLink{
					Page:     0,
					PageSize: 1024,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ALARM_FIELD",
							Key:  "createdTime",
						},
						Direction: "DESC",
					},
					TextSearch:             nil,
					TypeList:               []string{},
					SeverityList:           []string{},
					StatusList:             []string{"ACTIVE"},
					SearchPropagatedAlarms: false,
					AssigneeID:             nil,
					TimeWindow:             604800000,
				},
			},
			EntityFields: []EntityField{},
			LatestValues: []LatestValue{},
			CmdID:        10,
		},
		// Command 11: Entity Data for Water Level
		{
			Type: "ENTITY_DATA",
			Query: Query{
				EntityFilter: EntityFilter{
					Type:             "deviceType",
					ResolveMultiple:  true,
					DeviceTypes:      []string{"Dragino LDDS Water Level"},
					DeviceNameFilter: "",
				},
				PageLink: PageLink{
					Page:     0,
					PageSize: 1024,
					SortOrder: SortOrder{
						Key: struct {
							Type string `json:"type"`
							Key  string `json:"key"`
						}{
							Type: "ATTRIBUTE",
							Key:  "displayName",
						},
						Direction: "ASC",
					},
					TextSearch: nil,
					Dynamic:    true,
				},
			},
			EntityFields: []EntityField{
				{Type: "ENTITY_FIELD", Key: "name"},
				{Type: "ENTITY_FIELD", Key: "label"},
				{Type: "ENTITY_FIELD", Key: "additionalInfo"},
			},
			LatestValues: []LatestValue{
				{Type: "ATTRIBUTE", Key: "displayName"},
				{Type: "TIME_SERIES", Key: "waterLevel"},
			},
			CmdID: 11,
		},
	}
}

func generateDeviceId() string {
	return deviceID
}
