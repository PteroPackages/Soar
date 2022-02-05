import yaml from 'yaml';
import { formatString } from './response';
import { COLOURS } from '../log';

export function viewAs(type: string, _old: object, _new: object): [string, string] {
    switch (type) {
        case 'json':{
            return [JSON.stringify(_old, null, 2), JSON.stringify(_new, null, 2)];
        }
        case 'yaml':{
            return [yaml.stringify(_old), yaml.stringify(_new)];
        }
        case 'text':{
            return [formatString(_old), formatString(_new)];
        }
    }
}

export interface View {
    output:       string;
    additions:    number;
    subtractions: number;
    totalChanges: number;
}

export default function parseDiffView(type: string, _old: object, _new: object): View {
    const [a, b] = viewAs(type, _old, _new);

    const res: string[] = [];
    const indexA: string[] = a.split('\n');
    const indexB: string[] = b.split('\n');

    const view: View = {
        output: '',
        additions: 0,
        subtractions: 0,
        totalChanges: 0
    };

    for (let i=0; i<indexA.length; i++) {
        if (indexA[i] !== indexB[i]) {
            if (indexA[i].length > indexB[i].length) {
                res.push('- '+ indexB[i]);
                view.subtractions++;
            } else {
                res.push('+ '+ indexB[i]);
                view.additions++;
            }
            view.totalChanges++;
        } else {
            res.push('  '+ indexA[i]);
        }
    }

    view.output = res.join('\n');
    return view;
}

export function highlight(input: string): string {
    const res: string[] = [];

    for (const line of input.split('\n')) {
        if (line.startsWith('+')) {
            res.push(COLOURS.GREEN + line + COLOURS.RESET);
        } else if (line.startsWith('-')) {
            res.push(COLOURS.RED + line + COLOURS.RESET);
        } else {
            res.push(line);
        }
    }

    return res.join('\n');
}
