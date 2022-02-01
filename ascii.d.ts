/**
 * This is a hand-typed custom declaration file for the `ascii-table`
 * module as it does not contain one or reference to any pregenerated
 * or importable declaration files. PteroPackages has no affiliation
 * with the package or author with the package, and as such this file
 * is not copyrighted and can be used/modified by anyone.
 */

export interface TableOptions {
    prefix?: string;
}

export class AsciiTable {
    public options: TableOptions;

    public static VERSION: string;
    public static LEFT:    number;
    public static CENTER:  number;
    public static RIGHT:   number;

    constructor(name: string, options?: TableOptions);
    public static factory(name: string, options?: TableOptions): AsciiTable;

    public __name:         string;
    public __nameAlign:    number;
    public __rows:         any[];
    public __maxCells:     number;
    public __aligns:       any[];
    public __colMaxes:     any[];
    public __spacing:      number;
    public __heading:      any;
    public __headingAlign: number;
    public __edge:         string;
    public __fill:         string;
    public __top:          string;
    public __bottom:       string;
    public __border:       boolean;
    public __justify:      boolean;

    public align(dir: number, str?: string, len?: number, pad?: string): string;
    public alignLeft(str?: string, len?: number, pad?: string): string;
    public alignCenter(str?: string, len?: number, pad?: string): string;
    public alignRight(str?: string, len?: number, pad?: string): string;
    public alignAuto(str?: string, len?: number, pad?: string): string;

    public arrayFill(len: number, fill: any): any[];
    public arrayFill<T>(len: number, fill: T): T[];

    public reset(name: string | any): this;

    public setBorder(edge?: string, vertical?: string, top?: string, bottom?: string): this;
    public removeBorder(): this;
    public setAlign(idx: number, dir: number): this;
    public setTitle(name: string): this;
    public getTitle(): string;
    public setTitleAlign(dir: number): this;
    public sort(method: (a: any, b: any) => number): this;
    public sort<T>(method: (a: T, b: T) => number): this;
    public sortColumn(idx: number, method: (a: any, b: any) => number): this;
    public sortColumn<T>(idx: number, method: (a: T, b: T) => number): this;
    public setHeading(row: string | any, ...args: any[]): this;
    public setHeading<T>(row: string | any, ...args: T[]): this;
    public getHeading(): any[];
    public getHeading<T>(): T[];
    public setHeadingAlign(dir: number): this;
    public addRow(row: string | any, ...args: any[]): this;
    public addRow<T>(row: string | any, ...args: T[]): this;
    public getRows(): any[];
    public getRows<T>(): T[];
    public addRowMatrix(rows: any): this;
    public addRowMatrix<T>(rows: T): this;
    public addData(data: any[], rowCallback: (data: any) => any, asMatrix?: boolean): this;
    public addData<T>(data: T[], rowCallback: (data: T) => any, asMatrix?: boolean): this;
    public addData<T, U>(data: T[], rowCallback: (data: U) => any, asMatrix?: boolean): this;
    public clearRows(): this;
    public setJustify(val: boolean, ...args: any[]): this;
    public setJustify<T>(val: boolean, ...args: T[]): this;
    public toJSON(): object;
    public parse(obj: object): this;
    public fromJSON(obj: object): this;
    public render(): string;
    public valueOf(): string;
    public toString(): string;

    private _seperator(len: number, sep?: string): string;
    private _rowSeperator(): string;
    private _renderTitle(len: number): string;
    private _renderRow(row: any[], str: string, align: number): string;
}

export default AsciiTable;
