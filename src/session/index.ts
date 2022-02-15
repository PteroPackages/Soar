import fetch, { Response } from 'node-fetch';
import { createInterface, Interface } from 'readline';
import yaml from 'yaml';
import { version } from '../../package.json';
import { Auth, Config, FlagOptions } from '../structs';
import { getConfig } from '../config/funcs';
import log from '../log';
import Spinner from '../log/spinner';
import { createRequestLog } from '../logs/funcs';
import * as response from './response';

export default class Session {
    public config:    Config;
    public auth:      Auth;
    public type:      string;
    public reader:    Interface;
    public spinner:   Spinner | null;
    public showDebug: boolean;

    constructor(type: 'application' | 'client', options: FlagOptions) {
        this.type = type;
        this.spinner = null;

        this.setOptions(options);
    }

    private async getConfig() {
        this.config = await getConfig(true);
        this.auth = this.config[this.type];
        if (!this.auth?.url || !this.auth?.key)
            log.error('MISSING_AUTH_APPLICATION', null, true);

        this.showDebug ||= this.config.logs.showDebug;
        if (!this.config.logs.useColour) log.disableColour();
    }

    private setOptions(options: FlagOptions) {
        if (options.silent) {
            this.showDebug = false;
        } else {
            this.spinner = new Spinner();
            this.showDebug = options.debugMode;
        }
    }

    private setLogs(message: string, success: string, error: string): void {
        this.spinner?.setMessage(message)
            .onEnd(t => success.replace('$', t.toString()))
            .onError(t => error.replace('$', t.toString()));
    }

    private log(message: string | string[]): void {
        if (!this.showDebug) return;
        log.debug(message);
    }

    private logHttp(method: string, path: string, res: Response): void {
        if (!this.config.http.saveRequests) return this.log('request not saved');
        this.log('attempting to save request');
        createRequestLog(
            {
                date: Date.now(),
                method,
                response: res.status,
                type: 'D',
                domain: this.auth.url,
                path
            },
            this.config.core.ignoreWarnings
        );
    }

    public async handleRequest(method: string, path: string, data?: object) {
        await this.getConfig();
        this.log([
            'starting http request',
            `url: '${this.auth.url + path}'`,
            `method: ${method}`,
            `payload: ${data ? getByteSize(data) : '0'} bytes`
        ]);

        const base = path.slice(4).split('?')[0];
        this.setLogs(
            log.parse(`%y${method}%R ${base}`, 'http'),
            log.parse(`%g${method}%R ${base} ($ms taken)`, 'http'),
            log.parse(`%r${method}%R ${base} ($ms timeout)`, 'http')
        );
        this.spinner?.start();

        const res = await fetch(this.auth.url + path, {
            method,
            headers:{
                'Content-Type': 'application/json',
                'Accept': 'application/json',
                'Authorization': `Bearer ${this.auth.key}`,
                'User-Agent': `Soar Client v${version}`
            },
            body: data ? JSON.stringify(data) : null
        });

        if (res.status === 204) {
            this.spinner?.stop(false);
            this.log('received status: 204; request ended with no response body');
            this.logHttp(method, path, res);
            return Promise.resolve<void>(null);
        }

        if ([200, 201].includes(res.status)) {
            this.spinner?.stop(false);
            this.log(`received status: ${res.status}`);
            this.logHttp(method, path, res);

            if (res.headers.get('content-type') === 'application/json') {
                const json = await res.json();
                this.log([
                    'json response received',
                    `body: ${getByteSize(json)} bytes`
                ]);
                if (this.config.http.sendFullBody) return json;
                return json['data'] || json['attributes'];
            }

            this.log('buffer response body received, attempting to resolve...');
            const buf = await res.buffer();
            this.log(`buffer size: ${getByteSize(buf)} bytes`);
            return buf;
        }

        this.spinner?.stop(true);
        this.log(`received status: ${res.status}`);
        this.logHttp(method, path, res);
        if (res.status === 429 && this.config.http.retryRatelimit) {
            if (!this.config.core.ignoreWarnings) log.warn('ratelimit received, retrying...');
            this.log("attempting new request from 'http.retryRatelimit'");
            return this.handleRequest(method, path, data);
        }

        if (res.status >= 400 && res.status < 500)
            return log.fromPtero(await res.json(), true);

        log.error(
            'API Error',
            [
                `status code ${res.status} receieved`,
                'the api could not be contacted securely',
                'please contact a system administrator to resolve'
            ]
        );
    }

    public async handleClose(data: object, options: FlagOptions) {
        let parsed: string;
        this.setOptions(options);

        switch (options.responseType) {
            case 'text': parsed = response.formatString(data); break;
            case 'yaml': parsed = yaml.stringify(data); break;
            default: parsed = JSON.stringify(data); break;
        }

        if (options.writeFile.length) {
            this.log(`writing response to: '${options.writeFile}'`);
            response.writeFileResponse(options.writeFile, parsed, !options.silent);
        } else if (options.prompt && !options.silent) {
            this.reader ??= createInterface(
                process.stdin,
                process.stdout
            );
            this.log('input reader created');

            const res = await response.getBoolInput(
                this.reader, 'should this request be saved? (y/n)'
            );

            if (res) {
                const fp = await response.getStringInput(
                    this.reader,
                    'enter the file path to save to, leave empty for default path',
                    true
                );

                if (fp) await response.writeFileResponse(
                    `soar_log_${Date.now()}.${options.responseType}`,
                    parsed, !options.silent
                );
            }
            this.reader.close();
            this.log('input reader closed');
        }

        return parsed;
    }
}

function getByteSize(o: any): number {
    let size = 0;

    switch (typeof o) {
        case 'number': case 'bigint': size += 8; break;
        case 'string': size += o.length * 2; break;
        case 'boolean': size += 4; break;
        case 'object':{
            if (o === null || o === undefined) break;
            if (Array.isArray(o)) {
                size += o.reduce<any>((a, b) => getByteSize(b) + a, 0);
            } else {
                size += Object.values<any>(o).reduce((a, b) => getByteSize(b) + a, 0);
            }
        }
    }

    return size;
}
