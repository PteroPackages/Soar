import { FlagOptions } from '../structs';

export function parseUserGroup(args: object): FlagOptions {
    const type = (args['text'] && 'text') || (args['yaml'] && 'yaml') || 'json';
    let file = '';

    if (args['output']) {
        if (typeof args['output'] === 'boolean') file = `soar_log_${Date.now()}`;
        else file = args['output'];
    }
    if (!file.endsWith('.'+ type)) file += '.'+ type;

    return {
        prompt: args['silent'],
        writeFile: file,
        responseType: type
    } as FlagOptions;
}
