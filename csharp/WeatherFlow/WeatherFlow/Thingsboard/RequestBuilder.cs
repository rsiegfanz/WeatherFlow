namespace WeatherFlow.Thingsboard;

public static class RequestBuilder
{
    private const string DeviceId = "26945210-05ec-11ef-ac80-dde635ebcdb2";

    private static readonly EntityField[] DefaultEntityFields =
    [
        new() { Type = "ENTITY_FIELD", Key = "name" },
        new() { Type = "ENTITY_FIELD", Key = "label" },
        new() { Type = "ENTITY_FIELD", Key = "additionalInfo" }
    ];

    private static readonly LatestValue[] WeatherKeys =
    [
        new() { Type = "TIME_SERIES", Key = "airTemperature" },
        new() { Type = "TIME_SERIES", Key = "airHumidity" },
        new() { Type = "TIME_SERIES", Key = "barometricPressure" },
        new() { Type = "TIME_SERIES", Key = "windSpeed" },
        new() { Type = "TIME_SERIES", Key = "windDirectionSensor" },
        new() { Type = "TIME_SERIES", Key = "rainGauge" },
        new() { Type = "TIME_SERIES", Key = "uvIndex" },
        new() { Type = "TIME_SERIES", Key = "lightIntensity" },
        new() { Type = "TIME_SERIES", Key = "battery" }
    ];

    private static readonly LatestValue[] AttributeKeys =
    [
        new() { Type = "ATTRIBUTE", Key = "latitude" },
        new() { Type = "ATTRIBUTE", Key = "longitude" },
        new() { Type = "ATTRIBUTE", Key = "altitude" },
        new() { Type = "ATTRIBUTE", Key = "firmwareVersion" },
        new() { Type = "ATTRIBUTE", Key = "hardwareVersion" },
        new() { Type = "ATTRIBUTE", Key = "active" }
    ];

    private static readonly LatestValue[] WaterLevelKeys =
    [
        new() { Type = "ATTRIBUTE", Key = "displayName" },
        new() { Type = "TIME_SERIES", Key = "waterLevel" }
    ];

    public static ThingsBoardRequest Build(string token) => new()
    {
        AuthCmd = new AuthCommand { Token = token },
        Cmds =
        [
            // Cmd 1: All weather telemetry with live subscription
            new Command
            {
                Type = "ENTITY_DATA", CmdId = 1,
                LatestCmd = new LatestCmd { Keys = WeatherKeys },
                Query = new Query
                {
                    EntityFilter = DeviceFilter(),
                    PageLink = DefaultPageLink(),
                    EntityFields = DefaultEntityFields,
                    LatestValues = WeatherKeys
                }
            },
            // Cmd 2: Device attributes with live subscription
            new Command
            {
                Type = "ENTITY_DATA", CmdId = 2,
                LatestCmd = new LatestCmd { Keys = AttributeKeys },
                Query = new Query
                {
                    EntityFilter = DeviceFilter(),
                    PageLink = DefaultPageLink(),
                    EntityFields = DefaultEntityFields,
                    LatestValues = AttributeKeys
                }
            },
            // Cmd 10: Active alarms for water level sensors
            new Command
            {
                Type = "ALARM_DATA", CmdId = 10,
                Query = new Query
                {
                    EntityFilter = WaterLevelFilter(),
                    PageLink = new PageLink
                    {
                        PageSize = 1024,
                        SortOrder = new SortOrder
                        {
                            Key = new SortKey { Type = "ALARM_FIELD", Key = "createdTime" },
                            Direction = "DESC"
                        },
                        StatusList = ["ACTIVE"],
                        TimeWindow = 604800000
                    },
                    AlarmFields =
                    [
                        new EntityField { Type = "ALARM_FIELD", Key = "originatorLabel" },
                        new EntityField { Type = "ALARM_FIELD", Key = "createdTime" },
                        new EntityField { Type = "ALARM_FIELD", Key = "type" },
                        new EntityField { Type = "ALARM_FIELD", Key = "severity" }
                    ],
                    EntityFields = [],
                    LatestValues = []
                }
            },
            // Cmd 11: Water level devices with live subscription
            new Command
            {
                Type = "ENTITY_DATA", CmdId = 11,
                LatestCmd = new LatestCmd { Keys = WaterLevelKeys },
                Query = new Query
                {
                    EntityFilter = WaterLevelFilter(),
                    PageLink = new PageLink
                    {
                        PageSize = 1024,
                        SortOrder = new SortOrder
                        {
                            Key = new SortKey { Type = "ATTRIBUTE", Key = "displayName" },
                            Direction = "ASC"
                        },
                        Dynamic = true
                    },
                    EntityFields = DefaultEntityFields,
                    LatestValues = WaterLevelKeys
                }
            }
        ]
    };

    private static EntityFilter DeviceFilter() => new()
    {
        Type = "singleEntity",
        SingleEntity = new SingleEntity { Id = DeviceId }
    };

    private static EntityFilter WaterLevelFilter() => new()
    {
        Type = "deviceType",
        ResolveMultiple = true,
        DeviceTypes = ["Dragino LDDS Water Level"],
        DeviceNameFilter = ""
    };

    private static PageLink DefaultPageLink() => new()
    {
        PageSize = 1,
        SortOrder = new SortOrder
        {
            Key = new SortKey { Type = "ENTITY_FIELD", Key = "createdTime" },
            Direction = "DESC"
        }
    };
}
