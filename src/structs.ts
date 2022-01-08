export interface LogConfig {
    strictMode:       boolean;
    showDebug:        boolean;
    showHTTPLog:      boolean;
    showWebSocketLog: boolean;
    errorOutDir:      string;
}

export interface Auth {
    url: string;
    key: string;
}

export interface Config {
    version:     string;
    application: Auth;
    client:      Auth;
    logs:        LogConfig;
}

export interface AppUser {
    id:         number;
    uuid:       string;
    externalId: string | null;
    username:   string;
    email:      string;
    firstName:  string;
    lastName:   string;
    language:   string;
}

export function parseStruct<T>(data: any): T {
    const res = {} as unknown as T;
    for (const [k, v] of Object.entries(data)) {
        res[camelCase(k)] = v;
    }
    return res;
}

function camelCase(str: string): string {
    let res = '';
    let next = false;

    for (const c of str.split('')) {
        if (next) {
            next = false;
            res += c.toUpperCase();
        } else if (c === '_') {
            next = true;
        } else {
            res += c;
        }
    }
    return res;
}
