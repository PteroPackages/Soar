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

function assert(key: string, value: any, expected: any): void {
    if (typeof value !== typeof expected)
        throw `expected type ${typeof value} for '${key}'; got ${typeof expected}`;
}

export function parseConfig(config: any): string {
    try {
        assert('version', config.version, '');
        assert('application.url', config.application[0].url, '');
        assert('application.key', config.application[0].key, '');
        assert('client.url', config.client[0].url, '');
        assert('client.key', config.client[0].key, '');
        assert('logs.strictMode', config.logs.strict_mode, true);
        assert('logs.showDebug', config.logs.show_debug, true);
        assert('logs.showHttpLog', config.logs.show_http_log, true);
        assert('logs.showWsLog', config.logs.show_ws_log, true);
        assert('logs.logHttpRequests', config.logs.log_http_requests, true);
    } catch (err) {
        return err;
    }
}
