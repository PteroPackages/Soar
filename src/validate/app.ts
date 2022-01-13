import { FlagOptions } from '../structs';

export function parseUserGroup(args: object): FlagOptions {
    return {
        prompt: !!args['silent'],
        writeFile: (args['str'] && 'str') || (args['yaml'] && 'yaml') || 'json',
        responseType: args['output']
    } as FlagOptions;
}
