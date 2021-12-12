import {
    UserdataApi,
    UserdataApiFactory,
    AutomationsApi,
    AutomationsApiFactory,
    Configuration,
    GraphApi,
    GraphApiFactory,
    GroupsApi,
    GroupsApiFactory,
    LogsApi,
    LogsApiFactory,
    PlaybooksApi,
    PlaybooksApiFactory,
    RulesApi,
    RulesApiFactory,
    StatisticsApi,
    StatisticsApiFactory,
    TasksApi,
    TasksApiFactory,
    TemplatesApi,
    TemplatesApiFactory,
    TicketsApi,
    TicketsApiFactory,
    TickettypesApi,
    TickettypesApiFactory,
    UsersApi,
    UsersApiFactory,
    SettingsApi,
    SettingsApiFactory,
    JobsApi,
    JobsApiFactory,
} from "@/client";

const config = new Configuration({
    basePath:
        window.location.protocol +
        "//" +
        window.location.hostname +
        ":" +
        window.location.port +
        "/api"
});

export const API: TicketsApi &
    TemplatesApi &
    PlaybooksApi &
    RulesApi &
    AutomationsApi &
    UserdataApi &
    LogsApi &
    GraphApi &
    UsersApi &
    GroupsApi &
    StatisticsApi &
    SettingsApi &
    TickettypesApi &
    JobsApi &
    TasksApi = Object.assign(
    {},
    TicketsApiFactory(config),
    PlaybooksApiFactory(config),
    TemplatesApiFactory(config),
    RulesApiFactory(config),
    AutomationsApiFactory(config),
    SettingsApiFactory(config),
    LogsApiFactory(config),
    GraphApiFactory(config),
    UsersApiFactory(config),
    UserdataApiFactory(config),
    GroupsApiFactory(config),
    StatisticsApiFactory(config),
    SettingsApiFactory(config),
    TickettypesApiFactory(config),
    TasksApiFactory(config),
    SettingsApiFactory(config),
    JobsApiFactory(config)
);
