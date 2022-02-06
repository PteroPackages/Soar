export const ERRORS = {
    MISSING_ENV: "environment variable 'SOAR_PATH' is not set. Please set this to the path of the soar library",
    INVALID_ENV: "environment variable 'SOAR_PATH' path does not exist",
    CANNOT_READ_ENV: 'soar config file could not be opened. Please make sure soar has the necessary permissions',

    MISSING_CONFIG: "could not find config file. Make sure the environment variable 'SOAR_PATH' is set, or create a '.soar-local.yml' local config",
    MISSING_AUTH_APPLICATION: 'the necessary application api url or key is missing from the config',
    MISSING_AUTH_CLIENT: 'the necessary client api url or key is missing from the config',

    NOT_FOUND_USER: 'an account with the associated identifiers was not found',
    NOT_FOUND_SERVER: 'a server with the associated identifiers was not found',
    NOT_FOUND_NODE: 'a node with the associated identifiers was not found',
    NOT_FOUND_LOCATION: 'a node location with the associated identifiers was not found'
}

export function get(key: string): string {
    return ERRORS[key.toUpperCase()];
}

export default {
    ERRORS,
    get
}
