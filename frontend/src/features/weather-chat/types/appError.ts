export type AppErrorCode =
    | "CITY_NOT_FOUND"
    | "INVALID_WEATHER_REQUEST"
    | "WEATHER_SERVICE_CONFIG_ERROR"
    | "WEATHER_SERVICE_UNAVAILABLE"
    | "BACKEND_UNAVAILABLE"
    | "OPENAI_API_KEY_MISSING"
    | "MESSAGE_REQUIRED"
    | "UNKNOWN_ERROR";

export type SerializedAppError = {
    __app_error: true;
    status: number;
    code: AppErrorCode | string;
    message: string;
    technicalMessage?: string;
};

export class AppError extends Error {
    code: string;
    status: number;
    userMessage: string;
    technicalMessage?: string;

    constructor(params: {
        code: string;
        status: number;
        userMessage: string;
        technicalMessage?: string;
    }) {
        super(params.technicalMessage || params.userMessage);

        this.name = "AppError";
        this.code = params.code;
        this.status = params.status;
        this.userMessage = params.userMessage;
        this.technicalMessage = params.technicalMessage;
    }
}

export function isAppError(error: unknown): error is AppError {
    return error instanceof AppError;
}

export function appErrorFromSerialized(error: SerializedAppError): AppError {
    return new AppError({
        code: error.code,
        status: error.status,
        userMessage: error.message,
        technicalMessage: error.technicalMessage,
    });
}

export function isSerializedAppError(value: unknown): value is SerializedAppError {
    if (!value || typeof value !== "object") {
        return false;
    }

    const possibleError = value as Partial<SerializedAppError>;

    return (
        possibleError.__app_error === true &&
        typeof possibleError.status === "number" &&
        typeof possibleError.code === "string" &&
        typeof possibleError.message === "string"
    );
}