export default class ProgressWaiter {
    private startAt: number;
    private interval: NodeJS.Timer;
    public message: string;
    public endFunc: (time: number) => string | void;
    public count: number;

    constructor(message: string) {
        this.startAt = 0;
        this.message = message;
        this.count = 0;
    }

    public onEnd(func: (time: number) => string | void): this {
        this.endFunc = func;
        return this;
    }

    public start() {
        if (!this.message) throw new Error('No message provided.');
        this.startAt = Date.now();
        this.write(true);
        this.interval = setInterval<void[]>(() => this.handle(), 200).unref();
    }

    public stop() {
        clearInterval(this.interval);
        const res = this.endFunc(Date.now() - this.startAt);
        process.stdout.clearLine(0);
        process.stdout.cursorTo(0);
        process.stdout.write(res || '\n');
    }

    private write(_new: boolean = false) {
        if (_new) process.stdout.write('\n');
        process.stdout.clearLine(0);
        process.stdout.cursorTo(0);
        process.stdout.write(this.message);
        process.stdout.cursorTo(-1);
    }

    private handle() {
        if (this.count === 0) {
            this.count++;
            this.message += '.';
        } else if (this.count === 3) {
            this.message = this.message.slice(0, -3);
            this.count = 0;
        } else {
            this.message = this.message.slice(0, -this.count);
            this.count++;
            this.message += '.'.repeat(this.count);
        }
        this.write();
    }
}
