export interface Auth {
    url: string;
    key: string;
}

export interface LogConfig {
    showDebug:     boolean;
    showHttp:      boolean;
    showWebsocket: boolean;
    useColour:     boolean;
}

export interface HttpConfig {
    saveRequests:   boolean;
    sendFullBody:   boolean;
    retryRatelimit: boolean;
}

export interface CoreConfig {
    ignoreWarnings: boolean;
    stopAtSysError: boolean;
    saveErrorLogs:  boolean;
}

export interface Config {
    version:     string;
    application: Auth;
    client:      Auth;
    logs:        LogConfig;
    http:        HttpConfig;
    core:        CoreConfig;
}

export interface ReqLog {
    date:     number;
    method:   string;
    response: number;
    type:     string;
    domain:   string;
    path:     string;
    ref?:     string;
}

export interface FlagOptions {
    writeFile:    string;
    responseType: string;
    prompt:       boolean;
    silent:       boolean;
    debugMode:    boolean;
}
