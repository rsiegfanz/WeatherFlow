namespace WeatherFlow.Thingsboard;

public static class ThingsBoardRequestBuilder
{
    public static ThingsBoardRequest CreateWeatherDataRequest(string deviceId, string token)
    {
        return new ThingsBoardRequest
        {
            AuthCmd = new AuthCommand
            {
                CmdId = 0,
                Token = token
            },
            Cmds =
            [
                CreateSensorCommand(deviceId, "airTemperature", 1),
                CreateSensorCommand(deviceId, "rainGauge", 2),
                CreateSensorCommand(deviceId, "windSpeed", 3),
                CreateWindDetailsCommand(deviceId, 6)
            ]
        };
    }

    private static Command CreateSensorCommand(string deviceId, string sensorKey, int cmdId)
    {
        return new Command
        {
            Type = "ENTITY_DATA",
            CmdId = cmdId,
            Query = new Query
            {
                EntityFilter = new EntityFilter
                {
                    Type = "singleEntity",
                    SingleEntity = new SingleEntity
                    {
                        Id = deviceId,
                        EntityType = "DEVICE"
                    }
                },
                PageLink = new PageLink
                {
                    PageSize = 1,
                    Page = 0,
                    SortOrder = new SortOrder
                    {
                        Key = new SortKey
                        {
                            Type = "ENTITY_FIELD",
                            Key = "createdTime"
                        },
                        Direction = "DESC"
                    }
                },
                EntityFields =
                [
                    new EntityField { Type = "ENTITY_FIELD", Key = "name" },
                    new EntityField { Type = "ENTITY_FIELD", Key = "label" },
                    new EntityField { Type = "ENTITY_FIELD", Key = "additionalInfo" }
                ],
                LatestValues =
                [
                    new LatestValue { Type = "TIME_SERIES", Key = sensorKey }
                ]
            }
        };
    }

    private static Command CreateWindDetailsCommand(string deviceId, int cmdId)
    {
        return new Command
        {
            Type = "ENTITY_DATA",
            CmdId = cmdId,
            Query = new Query
            {
                EntityFilter = new EntityFilter
                {
                    Type = "singleEntity",
                    SingleEntity = new SingleEntity
                    {
                        Id = deviceId,
                        EntityType = "DEVICE"
                    }
                },
                PageLink = new PageLink
                {
                    PageSize = 1,
                    Page = 0,
                    SortOrder = new SortOrder
                    {
                        Key = new SortKey
                        {
                            Type = "ENTITY_FIELD",
                            Key = "createdTime"
                        },
                        Direction = "DESC"
                    }
                },
                EntityFields =
                [
                    new EntityField { Type = "ENTITY_FIELD", Key = "name" },
                    new EntityField { Type = "ENTITY_FIELD", Key = "label" },
                    new EntityField { Type = "ENTITY_FIELD", Key = "additionalInfo" }
                ],
                LatestValues =
                [
                    new LatestValue { Type = "TIME_SERIES", Key = "windDirectionSensor" },
                    new LatestValue { Type = "TIME_SERIES", Key = "windSpeed" }
                ]
            }
        };
    }
}