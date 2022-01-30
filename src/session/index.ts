import fetch from 'node-fetch';
import { Auth, Config, FlagOptions } from '../structs';
import { getConfig } from '../config/funcs';
import log from '../log';
import Spinner from '../log/spinner';

export default class Session {
    public config:       Config;
    public auth:         Auth;
    public spinner:      Spinner | null;
    public showDebugLog: boolean;
    public showHttpLog:  boolean;

    constructor(type: 'application' | 'client', options: FlagOptions) {
        this.config = getConfig();
        this.auth = this.config[type];
        this.spinner = null;
        this.showDebugLog = this.config.logs.showDebug;
        this.showHttpLog = this.config.logs.showHttpLog;

        this.setOptions(options);
    }

    private setOptions(options: FlagOptions) {
        if (options.silent) {
            this.showDebugLog = false;
            this.showHttpLog = false;
        } else {
            this.spinner = new Spinner()
                .setMessage(log.parse('%yfetching%R /application/users', 'info'))
                .onEnd(t => log.parse(`%gfetched%R /application/users (${t}ms taken)`, 'info'))
                .onError(t => log.parse(`%rfetch failed%R /application/users (${t}ms timeout)`, 'info'));
        }
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

    public async handleRequest(method: string, path: string, data?: object) {
        this.log('debug', 'Starting HTTP request');
        this.log('http', `Sending a request to '${this.auth.url + path}'`);
        this.spinner?.start();

        const res = await fetch(this.auth.url + path, {
            method,
            headers:{
                'Content-Type': 'application/json',
                'Accept': 'application/json',
                'Authorization': `Bearer ${this.auth.key}`,
                'User-Agent': `Soar Client v0.0.1`
            },
            body: data ? JSON.stringify(data) : null
        });

        this.log('http', `Received status: ${res.status}`);

        if (res.status === 204) {
            this.log('debug', 'Request ended with no response body');
            this.spinner?.stop(false);
            return Promise.resolve<void>(null);
        }
        if ([200, 201].includes(res.status)) {
            this.spinner?.stop(false);
            if (res.headers.get('content-type') === 'application/json')
                return await res.json();

            this.log('debug', 'Buffer response body received, attempting to resolve...');
            return await res.buffer();
        }

        this.spinner?.stop(true);
        if (res.status >= 400 && res.status < 500) return log.fromPtero(await res.json(), true);

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