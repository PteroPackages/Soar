import errors from '../errors';

let useColour = true;

export function disableColour(): void {
    useColour = false;
}

export const COLOURS = {
    RESET: '\x1b[0m',
    BOLD: '\x1b[1m',
    RED: '\x1b[31m',
    YELLOW: '\x1b[33m',
    GREEN: '\x1b[32m',
    CYAN: '\x1b[96m',
    BLUE: '\x1b[34m',
    PURPLE: '\x1b[35m',
    MAGENTA: '\x1b[95m'
}

export const CODE_MAP = {
    '%R': COLOURS.RESET,
    '%B': COLOURS.BOLD,
    '%r': COLOURS.RED,
    '%y': COLOURS.YELLOW,
    '%g': COLOURS.GREEN,
    '%c': COLOURS.CYAN,
    '%b': COLOURS.BLUE,
    '%p': COLOURS.PURPLE,
    '%m': COLOURS.MAGENTA
}

const BASE = `${COLOURS.RESET}[soar]${COLOURS.RESET}`;

export function parse(message: string, type?: string): string {
    if (
        !['info', 'success', 'notice', 'debug', 'http', 'warn', 'error', undefined]
        .includes(type)
    ) throw new Error('Invalid log type');

    if (useColour) {
        for (const [k, c] of Object.entries<string>(CODE_MAP)) message = message.replaceAll(k, c);
    } else {
        message = message.replaceAll(/%\w/gi, '');
    }

    let fmt = '';
    if (type) {
        switch (type) {
            case 'info': fmt = useColour ? `${COLOURS.BLUE}info${COLOURS.RESET}` : 'info'; break;
            case 'success': fmt = useColour ? `${COLOURS.GREEN}success${COLOURS.RESET}` : 'success'; break;
            case 'notice': fmt = useColour ? `${COLOURS.CYAN}notice${COLOURS.RESET}` : 'notice'; break;
            case 'debug': fmt = useColour ? `${COLOURS.BOLD}debug${COLOURS.RESET}` : 'debug'; break;
            case 'http': fmt = useColour ? `${COLOURS.MAGENTA}http${COLOURS.RESET}` : 'http'; break;
            case 'warn': fmt = useColour ? `${COLOURS.YELLOW}warning${COLOURS.RESET}` : 'warning'; break;
            case 'error': fmt = useColour ? `${COLOURS.RED}error${COLOURS.RESET}` : 'error'; break;
        }
    }

    return `${BASE} ${fmt}: ${message}`;
}

export function print(
    type: string,
    message: string | string[],
    _return: boolean
): string | void {
    const fmt = Array.isArray(message) ? message : [message];
    let res: string[] = [];
    for (const m of fmt) res.push(parse(m, type));
    if (_return) return res.join('\n');
    console.log(res.join('\n'));
}

export function info(message: string | string[], _return: boolean = false): void {
    print('info', message, _return);
}

export function notice(message: string | string[], _return: boolean = false): void {
    print('notice', message, _return);
}

export function debug(message: string | string[], _return: boolean = false): void {
    print('debug', message, _return);
}

export function http(message: string | string[], _return: boolean = false): void {
    print('http', message, _return);
}

export function success(message: string | string[], _return: boolean = false): void {
    print('success', message, _return);
}

export function warn(message: string | string[], _return: boolean = false): void {
    print('warn', message, _return);
}

export function error(name: string, message?: string | string[], exit?: boolean): void | never {
    const border = parse('', 'error');
    let fmt: string[] = [name];

    if (!message && errors.get(name)) {
        fmt.push(errors.get(name));
    } else {
        if (Array.isArray(message)) {
            fmt = fmt.concat(message);
        } else {
            if (message) fmt.push(message);
        }
    }

    console.log(fmt.map(m => border + m).join('\n'));
    if (exit) process.exit(1);
}

export function fromError(_error: Error, exit?: boolean): void | never {
    return error('Internal Error', _error.stack.split('\n'), exit);
}

interface pteroError {
    errors:{
        code:   string;
        status: string;
        detail: string;
    }[];
}

export function fromPtero(data: pteroError, exit?: boolean): void | never {
    error(
        'API Request Error',
        `pterodactyl panel returned ${data.errors.length} error${data.errors.length > 1 ? 's' : ''}.`
    );

    for (const err of data.errors) {
        error('');
        error(err.code, err.detail || '[no details received]');
    }

    if (data.errors.some(e => e.status === '403')) {
        notice(
            'please ensure that your api key has the necessary'+
            ' read/write permissions before making requests.'
        );
    }

    if (exit) process.exit(1);
}

export default {
    parse,
    print,
    info,
    notice,
    debug,
    http,
    success,
    warn,
    error,
    fromError,
    fromPtero,
    disableColour
}
