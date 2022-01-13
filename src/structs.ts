export interface LogConfig {
    strictMode:  boolean;
    showDebug:   boolean;
    showHttpLog: boolean;
    showWsLog:   boolean;
    errorOutDir: string;
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

    for (let [k, v] of Object.entries(data)) {
        if (
            typeof v === 'object' &&
            v !== undefined
        ) v = parseStruct<unknown>(v);
        res[camelCase(k)] = v;
    }

    return res;
}

export function jsonStruct<T>(data: T): object {
    const res = {};

    for (let [k, v] of Object.entries(data)) {
        if (
            typeof v === 'object' &&
            v !== undefined
        ) v = jsonStruct(v);
        res[snakeCase(k)] = v;
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

function snakeCase(str: string): string {
    let res = '';

    const isUpper = (c: string) => 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'.includes(c);

    for (const c of str.split('')) {
        if (isUpper(c)) res += '_';
        res += c.toLowerCase();
    }

    return res;
}
