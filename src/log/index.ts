export const COLOURS = {
    BASE: '\x1b[',
    RESET: '\x1b[0m',
    RED: '\x1b[31m',
    YELLOW: '\x1b[33m',
    GREEN: '\x1b[32m',
    CYAN: '\x1b[36m',
    BLUE: '\x1b[34m',
    MAGENTA: '\x1b[35m'
}

const BASE = `[${COLOURS.CYAN}soar${COLOURS.RESET}]`;

export function info(message: string): void {
    console.log(`${BASE} ${COLOURS.BLUE}info${COLOURS.RESET}: ${message}`);
}

export function warn(type: string, message: string[]): void {
    const border = `${BASE} ${COLOURS.YELLOW}warning${COLOURS.RESET}: `;
    let fmt = border + type +'\n';
    fmt += message.map(m => border + m).join('\n');
    console.log(fmt);
}

export function error(type: string, message: string | string[], exit?: boolean): void | never {
    const border = `${BASE} ${COLOURS.RED}error${COLOURS.RESET}: `;
    let fmt = '';
    if (type.length) fmt = border + type + '\n';

    if (Array.isArray(message)) {
        fmt += message.map(m => border + m).join('\n');
    } else {
        fmt += border + message;
    }

    console.log(fmt);
    if (exit) process.exit(1);
}

export function fromError(_error: Error, exit?: boolean): void | never {
    return error('', _error.stack, exit);
}
