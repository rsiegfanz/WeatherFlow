using System.Text.Json.Serialization;

namespace WeatherFlow.Thingsboard;

// --- Request DTOs ---

public class ThingsBoardRequest
{
    [JsonPropertyName("cmds")] public Command[] Cmds { get; set; } = [];
    [JsonPropertyName("authCmd")] public AuthCommand AuthCmd { get; set; } = new();
}

public class Command
{
    [JsonPropertyName("type")] public string Type { get; set; } = "";
    [JsonPropertyName("query")] public Query Query { get; set; } = new();
    [JsonPropertyName("latestCmd")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public LatestCmd? LatestCmd { get; set; }
    [JsonPropertyName("cmdId")] public int CmdId { get; set; }
}

public class LatestCmd
{
    [JsonPropertyName("keys")] public LatestValue[] Keys { get; set; } = [];
}

public class Query
{
    [JsonPropertyName("entityFilter")] public EntityFilter EntityFilter { get; set; } = new();
    [JsonPropertyName("pageLink")] public PageLink PageLink { get; set; } = new();
    [JsonPropertyName("entityFields")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public EntityField[]? EntityFields { get; set; }
    [JsonPropertyName("latestValues")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public LatestValue[]? LatestValues { get; set; }
    [JsonPropertyName("alarmFields")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public EntityField[]? AlarmFields { get; set; }
}

public class EntityFilter
{
    [JsonPropertyName("type")] public string Type { get; set; } = "";
    [JsonPropertyName("singleEntity")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public SingleEntity? SingleEntity { get; set; }
    [JsonPropertyName("resolveMultiple")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingDefault)]
    public bool ResolveMultiple { get; set; }
    [JsonPropertyName("deviceTypes")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string[]? DeviceTypes { get; set; }
    [JsonPropertyName("deviceNameFilter")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string? DeviceNameFilter { get; set; }
}

public class SingleEntity
{
    [JsonPropertyName("id")] public string Id { get; set; } = "";
    [JsonPropertyName("entityType")] public string EntityType { get; set; } = "DEVICE";
}

public class PageLink
{
    [JsonPropertyName("pageSize")] public int PageSize { get; set; } = 1;
    [JsonPropertyName("page")] public int Page { get; set; }
    [JsonPropertyName("sortOrder")] public SortOrder SortOrder { get; set; } = new();
    [JsonPropertyName("dynamic")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingDefault)]
    public bool Dynamic { get; set; }
    [JsonPropertyName("statusList")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingNull)]
    public string[]? StatusList { get; set; }
    [JsonPropertyName("timeWindow")][JsonIgnore(Condition = JsonIgnoreCondition.WhenWritingDefault)]
    public int TimeWindow { get; set; }
}

public class SortOrder
{
    [JsonPropertyName("key")] public SortKey Key { get; set; } = new();
    [JsonPropertyName("direction")] public string Direction { get; set; } = "DESC";
}

public class SortKey
{
    [JsonPropertyName("type")] public string Type { get; set; } = "ENTITY_FIELD";
    [JsonPropertyName("key")] public string Key { get; set; } = "createdTime";
}

public class EntityField
{
    [JsonPropertyName("type")] public string Type { get; set; } = "";
    [JsonPropertyName("key")] public string Key { get; set; } = "";
}

public class LatestValue
{
    [JsonPropertyName("type")] public string Type { get; set; } = "";
    [JsonPropertyName("key")] public string Key { get; set; } = "";
}

public class AuthCommand
{
    [JsonPropertyName("cmdId")] public int CmdId { get; set; }
    [JsonPropertyName("token")] public string Token { get; set; } = "";
}

// --- Response DTOs ---

public class WsMessage
{
    [JsonPropertyName("cmdId")] public int CmdId { get; set; }
    [JsonPropertyName("errorCode")] public int ErrorCode { get; set; }
    [JsonPropertyName("errorMsg")] public string? ErrorMsg { get; set; }
    [JsonPropertyName("data")] public EntityDataPage? Data { get; set; }
    [JsonPropertyName("update")] public List<EntityDataEntry>? Update { get; set; }
}

public class EntityDataPage
{
    [JsonPropertyName("data")] public List<EntityDataEntry> Data { get; set; } = [];
    [JsonPropertyName("totalElements")] public int TotalElements { get; set; }
}

public class EntityDataEntry
{
    [JsonPropertyName("entityId")] public EntityId EntityId { get; set; } = new();
    [JsonPropertyName("latest")] public Dictionary<string, Dictionary<string, TsValue>> Latest { get; set; } = new();
}

public class EntityId
{
    [JsonPropertyName("entityType")] public string EntityType { get; set; } = "";
    [JsonPropertyName("id")] public string Id { get; set; } = "";
}

public class TsValue
{
    [JsonPropertyName("ts")] public long Ts { get; set; }
    [JsonPropertyName("value")] public string Value { get; set; } = "";
}
