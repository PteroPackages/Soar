import { Command } from 'commander';
import ascii from 'ascii-table';
import { fetchLogs } from './funcs';
import log from '../log';

const logsGetCmd = new Command('fetch')
    .description('Fetches Soar HTTP logs')
    .addHelpText('before', 'Fetches Soar HTTP logs with additional filters')
    .option('-s, --silent', 'Doesn\'t log response messages', false)
    .option('--from <date>', 'Gets logs from a specific date')
    .option('--method <method>', 'Gets logs using a specific HTTP method')
    .option('--app', 'Gets logs that used the application API', false)
    .option('--client', 'Gets logs that used the client API', false)
    .option('--raw', 'Gets logs that used the raw request method', false)
    .action((args: object) => {
        let date: Date;
        let method: string;

        if (args['from']) {
            if (/[a-z\/\-\:]+/g.test(args['from'])) {
                date = new Date(args['from']);
            } else {
                date = new Date(Number(args['from']));
            }

            if (date.toString() === 'Invalid Date') log.error(
                'Argument Error',
                [
                    'invalid date provided, valid formats are:',
                    '1643677200000, 01/02/2022, 2022-02-01T01:00:00.000Z'
                ],
                true
            );
        }

        if (args['method']) {
            if (
                !['GET', 'POST', 'PATCH', 'PUT', 'DELETE'].includes(
                    args['method']!.toUpperCase()
                )
            ) log.error(
                'Argument Error',
                [
                    'invalid http method specified',
                    'methods: GET, POST, PATCH, PUT, DELETE'
                ],
                true
            );
            method = args['method']!.toUpperCase();
        }

        let logs = fetchLogs();

        if (args['from']) logs = logs.filter(l => l.date >= date.getTime());
        if (args['method']) logs = logs.filter(l => l.method === method);
        if (args['app']) logs = logs.filter(l => l.path.includes('app'));
        if (args['client']) logs = logs.filter(l => l.path.includes('client'));
        if (args['raw']) logs = logs.filter(l => l.type !== 'D');

        if (!args['silent']) log.info([
            log.parse(`found %c${logs.length}%R logs`),
            `filter options: ${Object.values(args).filter(Boolean).length || 'none'}`,
            'results:\n'
        ]);

        const sortCode = (c: number) => {
            if (c < 200 || (c >= 300 && c < 400)) return log.parse(`%y${c}$R`);
            if (c >= 400) return log.parse(`%r${c}$R`);
            return log.parse(`%g${c}%R`);
        }

        const table = new ascii('Request Logs')
            .setHeading('Date', 'Method-Status', 'Domain', 'Path');

        logs.forEach(l => table.addRow(
            new Date(l.date).toLocaleString(),
            l.method +' '+ sortCode(l.response),
            l.domain,
            l.path
        ));

        console.log(table.render());
    });

export default [
    logsGetCmd
]
