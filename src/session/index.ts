import fetch, { Response } from 'node-fetch';
import { Auth, Config, FlagOptions } from '../structs';
import { getConfig } from '../config/funcs';
import log from '../log';
import Spinner from '../log/spinner';
import { createRequestLog } from '../logs/funcs';

export default class Session {
    public config:       Config;
    public auth:         Auth;
    public type:         string;
    public spinner:      Spinner | null;
    public showDebugLog: boolean;
    public showHttpLog:  boolean;

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

        this.showDebugLog = this.config.logs.showDebug;
        this.showHttpLog = this.config.logs.showHttpLog;
    }

    private setOptions(options: FlagOptions) {
        if (options.silent) {
            this.showDebugLog = false;
            this.showHttpLog = false;
        } else {
            this.spinner = new Spinner();
        }
    }

    private setLogs(message: string, success: string, error: string): void {
        this.spinner?.setMessage(message)
            .onEnd(t => success.replace('$', t.toString()))
            .onError(t => error.replace('$', t.toString()));
    }

    private log(type: string, message: string): void {
        if (this.spinner?.running) return;
        if (type === 'debug') {
            if (!this.showDebugLog) return;
            log.debug(message);
        } else {
            if (!this.showHttpLog) return;
            log.print(`%B${type}%R: ${message}`);
        }
    }

    private logHttp(method: string, path: string, res: Response): void {
        if (this.config.logs.logHttpRequests) createRequestLog({
            date: Date.now(),
            method,
            response: res.status,
            type: 'D',
            domain: this.auth.url,
            path
        });
    }

    public async handleRequest(method: string, path: string, data?: object) {
        await this.getConfig();
        this.log('debug', 'Starting HTTP request');
        this.log('http', `Sending a request to '${this.auth.url + path}'`);

        const base = path.slice(4).split('?')[0];
        this.setLogs(
            log.parse(`%yfetching%R ${base}`, 'info'),
            log.parse(`%gfetched%R ${base} ($ms taken)`, 'info'),
            log.parse(`%rfetch failed%R ${base} ($ms timeout)`, 'info')
        );
        this.spinner?.start();

        const res = await fetch(this.auth.url + path, {
            method,
            headers:{
                'Content-Type': 'application/json',
                'Accept': 'application/json',
                'Authorization': `Bearer ${this.auth.key}`,
                'User-Agent': `Soar Client v0.0.2`
            },
            body: data ? JSON.stringify(data) : null
        });

        this.log('http', `Received status: ${res.status}`);

        if (res.status === 204) {
            this.log('debug', 'Request ended with no response body');
            this.spinner?.stop(false);
            this.logHttp(method, path, res);
            return Promise.resolve<void>(null);
        }
        if ([200, 201].includes(res.status)) {
            this.spinner?.stop(false);
            this.logHttp(method, path, res);
            if (res.headers.get('content-type') === 'application/json')
                return await res.json();

            this.log('debug', 'Buffer response body received, attempting to resolve...');
            return await res.buffer();
        }

        this.spinner?.stop(true);
        this.logHttp(method, path, res);
        if (res.status >= 400 && res.status < 500)
            return log.fromPtero(await res.json(), true);

        log.error(
            'API Error',
            [
                `Status code ${res.status} receieved;`,
                'The API could not be contacted securely',
                'Please contact a system administrator to resolve.'
            ]
        );
    }
}