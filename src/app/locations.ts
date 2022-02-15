import { Command, Option } from 'commander';
import Session from '../session';
import parseDiffView, { highlight } from '../session/view';
import { buildLocation, parseFlagOptions } from '../validate';
import log from '../log';

const getLocationsCmd = new Command('get-locations')
    .description('Fetches node locations from the panel')
    .addHelpText('before', 'Fetches all node locations from the panel (can specify with flags)')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .option('--id <id>', 'The node location ID to fetch')
    .option('--nodes', 'Include nodes in the request', false)
    .option('--servers', 'Include servers in the request', false)
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        const data = await session.handleRequest('GET', buildLocation(args));
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('request results:\n');
            console.log(out);
        }
    });

const createLocationCmd = new Command('create-location')
    .description('Creates a new node location')
    .addHelpText('before', 'Creates a new node location')
    .option('--json', 'Send the response output as JSON', true)
    .option('--yaml', 'Send the response output as YAML', false)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .requiredOption('-d, --data <json>', 'The json data to create the location with')
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (args: object) => {
        const options = parseFlagOptions(args);

        let json: object;
        try {
            json = JSON.parse(args['data']);
        } catch (err) {
            log.error(
                'Argument Error',
                [
                    'couldn\'t parse json data argument:',
                    err.message
                ],
                true
            );
        }

        const missing: string[] = [];
        for (const key of ['short', 'long']) {
            if (key in json) continue;
            missing.push(key);
        }
        if (missing.length) log.error(
            'Argument Error',
            [
                `missing required key${missing.length > 1 ? 's' : ''}:`,
                missing.join(', ')
            ],
            true
        );

        const session = new Session('application', options);
        const data = await session.handleRequest('POST', buildLocation({}), json);
        const out = await session.handleClose(data, options);
        if (out) {
            if (!options.silent) log.success('location created! request result:\n');
            console.log(out);
        }
    });

const updateLocationCmd = new Command('update-location')
    .description('Updates a node location on the panel')
    .addHelpText('before', 'Updates a specified node location on the panel')
    .argument('<id>', 'The ID of the node location to update')
    .option('--json', 'Send the response output as JSON', false)
    .option('--yaml', 'Send the response output as YAML', true)
    .option('--text', 'Send the response output as formatted text', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request', false)
    .option('-s, --silent', 'Don\'t log request messages', false)
    .option('-o, --output [file]', 'Writes the output to a file')
    .requiredOption('-d, --data <json>', 'The json data to update the user with')
    .option('--no-diff', 'Don\'t show the properties changed in the request', false)
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (id: string, args: object) => {
        const options = parseFlagOptions(args);

        let json: object;
        try {
            json = JSON.parse(args['data']);
        } catch (err) {
            log.error(
                'Argument Error',
                [
                    'couldn\'t parse json data argument:',
                    err.message
                ],
                true
            );
        }
        if (!Object.entries(json).length) log.error(
            'Argument Error',
            'no json was provided to update.',
            true
        );

        const session = new Session('application', options);
        const location = await session.handleRequest('GET', buildLocation({ id }));
        if (!location) log.error('NOT_FOUND_LOCATION', null, true);

        json['short'] ||= location['short'];
        json['long'] ||= location['long'];

        const data = await session.handleRequest('PATCH', buildLocation({ id }), json);
        const out = await session.handleClose(data, options);

        if (out && args['diff']) {
            const view = parseDiffView(options.responseType, location, data);
            log.print(
                'success',
                `made %c${view.totalChanges}%R changes`+
                ` (%g+${view.additions}%R | %r-${view.subtractions}%R)`,
                false
            );
            console.log(
                '\n'+ (session.config.logs.useColour
                ? highlight(view.output)
                : view.output)
            );
        } else {
            log.success(`updated location: ${id}`);
        }
    });

const deleteLocationCmd = new Command('delete-location')
    .description('Deletes a node location from the panel')
    .addHelpText('before', 'Deletes a node location from the panel')
    .argument('<id>', 'The ID of the node location to delete')
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .addOption(new Option('--debug').default(false).hideHelp())
    .action(async (id: string, args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        await session.handleRequest('DELETE', buildLocation({ id }));
        if (!options.silent) log.success(`deleted location: ${id}`);
    });

export default [
    getLocationsCmd,
    createLocationCmd,
    updateLocationCmd,
    deleteLocationCmd
]
