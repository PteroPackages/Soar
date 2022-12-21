package app

var getUsersHelp = "Gets a list of user accounts from the panel (supports the --id flag)."

var createUserHelp = "Creates a user account on the panel using the data provided from one of the following options:\n" +
	"'--data source' - takes a set of key-value pairs for arguments (e.g. \"username=example email=test@example.com\")\n" +
	"'--data-file file' - takes a file path to a JSON file with the data fields\n" +
	"'--data-json source' - takes a raw JSON data input\n\n" +
	"The username, email, first_name, last_name and root_admin fields are required.\n" +
	"The external_id and password fields are optional and are omitted by default."

var getServersHelp = "Gets a list of servers from the panel (supports the --id flag)."

var getNodesHelp = "Gets a list of nodes from the panel (supports the --id flag)."

var getLocationsHelp = "Gets a list of node locations from the panel (supports --id flag)."

var getNestsHelp = "Gets a list of nests from the panel (supports the --id flag)."

var getNestEggsHelp = "Gets a list of eggs for a specified nest from the panel (supports the --id flag)."
