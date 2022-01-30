export default class Spinner {
    private interval: NodeJS.Timer;
    public message:   string;
    public running:   boolean;
    public startedAt: number;
    public endFunc:   (time: number) => string;
    public errFunc:   (time: number) => string;

    public CHAR_REF: { [key: string]: string } = {
        '⠇': '⠏',
        '⠏': '⠋',
        '⠋': '⠙',
        '⠙': '⠹',
        '⠹': '⠸',
        '⠸': '⠼',
        '⠼': '⠴',
        '⠴': '⠦',
        '⠦': '⠇'
    }

    constructor() {
        this.message = null;
        this.running = false;
        this.startedAt = 0;
        this.endFunc = null;
        this.errFunc = null;
    }

    public setMessage(message: string): this {
        this.message = message.trim() + ' ⠇    ';
        return this;
    }

    public onEnd(func: (time: number) => string): this {
        this.endFunc = func;
        return this;
    }

    public onError(func: (time: number) => string): this {
        this.errFunc = func;
        return this;
    }

    public start(): this {
        if (this.running) return this;
        this.startedAt = Date.now();
        this.running = true;
        this.interval = setInterval(() => this.handle(), 100).unref();
        return this;
    }

    public stop(error: boolean): void {
        if (!this.running) return;
        this.running = false;
        clearInterval(this.interval);
        this.clear();
        let res: string;
        if (error) {
            res = this.errFunc?.(Date.now() - this.startedAt);
        } else {
            res = this.endFunc?.(Date.now() - this.startedAt);
        }
        process.stdout.write(res +'\n');
    }

    private handle() {
        const char = this.message.slice(-5).slice(0, -4);
        this.message = this.message.slice(0, -5) + this.CHAR_REF[char] + '    ';
        this.clear();
        process.stdout.write(this.message);
    }

    private clear(): this {
        process.stdout.clearLine(0);
        process.stdout.cursorTo(0);
        return this;
    }
}
