export const ERRORS = {
    MISSING_ENV: "Environment variable 'SOAR_PATH' is not set. Please set this to the path of the Soar library.",
    INVALID_ENV: "Environment variable 'SOAR_PATH' path does not exist.",
    CANNOT_READ_ENV: 'Soar config file could not be opened. Please make sure Soar has the necessary permissions.',

    MISSING_CONFIG: "Could not find config file. Make sure the environment variable 'SOAR_PATH' is set, or create a '.soar-local.yml' local config.",
    MISSING_AUTH_APPLICATION: 'The necessary application API URL or key is missing from the config.',
    MISSING_AUTH_CLIENT: 'The necessary client API URL or key is missing from the config.',

    NOT_FOUND_USER: 'An account with the associated identifiers was not found.',
    NOT_FOUND_SERVER: 'A server with the associated identifiers was not found.',
    NOT_FOUND_NODE: 'A node with the associated identifiers was not found.',
    NOT_FOUND_LOCATION: 'A node location with the associated identifiers was not found.'
}

export function get(key: string): string {
    return ERRORS[key.toUpperCase()];
}

export default {
    ERRORS,
    get
}
