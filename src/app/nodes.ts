import { Command } from 'commander';
import Session from '../session';
import { buildNode, parseFlagOptions } from '../validate';
import log from '../log';

const getNodesCmd = new Command('get-nodes')
    .addHelpText('before', 'Fetches all nodes from the panel (can specify with flags).')
    .option('--json', 'Send the response output as JSON.', true)
    .option('--yaml', 'Send the response output as YAML.', false)
    .option('--text', 'Send the response output as formatted text.', false)
    .option('-n, --no-prompt', 'Don\'t prompt for user response after the request.', false)
    .option('-s, --silent', 'Don\'t log request messages.', false)
    .option('-o, --output [file]', 'Writes the output to a file.')
    .option('--id <id>', 'The node ID to fetch.')
    .option('--config', 'Fetch the node configuration setup.', false)
    .action(async (args: object) => {
        const options = parseFlagOptions(args);
        const session = new Session('application', options);

        const data = await session.handleRequest(
            'GET',
            buildNode({ config: args['config'], ...args })
        );
        if (!options.silent) log.success('request result:\n');

        const out = await session.handleClose(data, options);
        if (out) console.log(out);
    });

export default [
    getNodesCmd
]
