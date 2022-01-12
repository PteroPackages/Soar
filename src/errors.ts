export const ERRORS = {
    MISSING_ENV: "Environment variable 'SOAR_PATH' is not set. Please set this to the path of the Soar library.",
    INVALID_ENV: "Environment variable 'SOAR_PATH' path does not exist.",
    CANNOT_READ_ENV: 'Soar config file could not be opened. Please make sure Soar has the necessary permissions.'
}

export function get(key: string): string {
    if (!ERRORS[key]) {
        const lowerKeys = Object.keys(ERRORS).map(k => k.toLowerCase());
        if (!lowerKeys.includes(key)) throw new Error(`Invalid error key '${key}'.`);
    }
    return ERRORS[key];
}

export function tryGet(key: string): string | null {
    try {
        return get(key);
    } catch {
        return null;
    }
}
