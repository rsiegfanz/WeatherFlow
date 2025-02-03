namespace WeatherFlow.Thingsboard;

public class ThingsBoardRequest
{
    public Command[] Cmds { get; set; }
    public AuthCommand AuthCmd { get; set; }
}

public class Command
{
    public string Type { get; set; }
    public Query Query { get; set; }
    public int CmdId { get; set; }
}

public class Query
{
    public EntityFilter EntityFilter { get; set; }
    public PageLink PageLink { get; set; }
    public EntityField[] EntityFields { get; set; }
    public LatestValue[] LatestValues { get; set; }
}

public class EntityFilter
{
    public string Type { get; set; }
    public SingleEntity SingleEntity { get; set; }
}

public class SingleEntity
{
    public string Id { get; set; }
    public string EntityType { get; set; }
}

public class PageLink
{
    public int PageSize { get; set; }
    public int Page { get; set; }
    public SortOrder SortOrder { get; set; }
}

public class SortOrder
{
    public SortKey Key { get; set; }
    public string Direction { get; set; }
}

public class SortKey
{
    public string Type { get; set; }
    public string Key { get; set; }
}

public class EntityField
{
    public string Type { get; set; }
    public string Key { get; set; }
}

public class LatestValue
{
    public string Type { get; set; }
    public string Key { get; set; }
}

public class AuthCommand
{
    public int CmdId { get; set; }
    public string Token { get; set; }
}