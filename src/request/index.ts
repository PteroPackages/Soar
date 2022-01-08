import fetch from 'node-fetch';
import { getConfig } from '../config/funcs';

export async function handleRequest(method: string, path: string, data?: object) {
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

    if (res.status === 201) return;
    if (res.status === 200) return await res.json();
    // TODO: handle 400-500 here
    return;
}
