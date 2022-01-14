import fetch from 'node-fetch';
import * as log from '../log';
import { getConfig } from '../config/funcs';

export async function handleRequest(method: string, path: string, data?: object): Promise<any |void> {
    if (['head', 'options', 'trace'].includes(method))
        throw new Error(`Unsupported Pterodactyl API request method '${method}'.`);

    const config = getConfig();
    const auth = path.includes('client')
        ? config.client
        : config.application;
    path = auth.url + path;

    const res = await fetch(path, {
        method,
        headers:{
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'Authorization': `Bearer ${auth.key}`,
            'User-Agent': `Soar Client v0.0.1`
        },
        body: data ? JSON.stringify(data) : null
    });

    if (res.status === 204) return Promise.resolve<void>(null);
    if ([200, 201].includes(res.status)) return await res.json();
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
