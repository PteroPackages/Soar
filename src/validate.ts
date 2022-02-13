import { FlagOptions } from './structs';

export function parseFlagOptions(args: object): FlagOptions {
    const type = (args['text'] && 'text') || (args['yaml'] && 'yaml') || 'json';
    const silent = !process.stdout.isTTY || args['silent'];
    let file = '';

    if (args['output']) {
        if (typeof args['output'] === 'boolean') file = `soar_log_${Date.now()}`;
        else file = args['output'];
    }
    if (file.length && !file.endsWith('.'+ type)) file += '.'+ type;

    return {
        silent,
        prompt: args['prompt'],
        writeFile: file,
        responseType: type
    } as FlagOptions;
}

export function buildUser(args: object): string {
    let base = '/api/application/users';
    if (args['id']) return `${base}/${args['id']}`;
    if (args['email']) return `${base}?filter[email]=${args['email']}`;
    if (args['uuid']) return `${base}?filter[uuid]=${args['uuid']}`;
    if (args['username']) return `${base}?filter[username]=${args['username']}`;
    if (args['external']) return `${base}?filter[external_id]=${args['external']}`;
    return base;
}

export function buildServer(args: object): string {
    let base = '/api/application/servers';
    if (args['id']) {
        base = `/api/application/servers/${args['id']}`;
        if (args['suspend']) return base +'/suspend';
        if (args['unsuspend']) return base +'/unsuspend';
        if (args['reinstall']) return base +'/reinstall';
    }
    if (args['uuid']) return `${base}?filter[uuid]=${args['uuid']}`;
    if (args['name']) return `${base}?filter[name]=${args['name']}`;
    if (args['external']) return `${base}?filter[external_id]=${args['external']}`;
    if (args['image']) return `${base}?filter[image]=${args['image']}`;
    return base;
}

export function buildNode(args: object): string {
    let base = '/api/application/nodes';
    if (args['id']) base += `/${args['id']}`;
    if (args['config']) return base +'/configuration';
    return base;
}

function assertType(key: string, value: any, expected: any): void {
    if (typeof value !== typeof expected)
        throw `expected type ${typeof expected} for '${key}'; got ${typeof value}`;
}

function assertValue(key: string, value: any, expected: string): void {
    if (!value) throw `expected ${expected} for '${key}'`;
}

export function parseConfig(config: any): string {
    try {
        assertType('version', config.version, '');
        assertType('application.url', config.application.url, '');
        assertType('application.key', config.application.key, '');
        assertType('client.url', config.client.url, '');
        assertType('client.key', config.client.key, '');
        assertType('logs.showDebug', config.logs.show_debug, true);
        assertType('logs.showHttp', config.logs.show_http, true);
        assertType('logs.showWebsocket', config.logs.show_websocket, true);
        assertType('logs.useColour', config.logs.use_colour, true);
        assertType('http.saveRequests', config.http.save_requests, true);
        assertType('http.sendFullBody', config.http.send_full_body, true);
        assertType('http.retryRatelimit', config.http.retry_ratelimit, true);
        assertType('core.ignoreWarnings', config.core.ignore_warnings, true);
        assertType('core.stopAtSysError', config.core.stop_at_sys_error, true);
        assertType('core.saveErrorLogs', config.core.save_error_logs, true);

        assertValue(
            'version',
            /\d\.\d\.\d/g.test(config.version),
            'version to match semver spec'
        );
        assertValue(
            'application.url',
            /https?\:\/\//g.test(config.application.url),
            "url to start with 'http://' or 'https://'"
        );
        assertValue(
            'client.url',
            /https?\:\/\//g.test(config.application.url),
            "url to start with 'http://' or 'https://'"
        );
    } catch (err) {
        return err;
    }
}
