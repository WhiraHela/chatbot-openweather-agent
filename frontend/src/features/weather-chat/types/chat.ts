export type ChatMessage = {
    role: "user" | "assistant";
    content: string;
};

export type ChatApiSuccessResponse = {
    answer: string;
};

export type ChatApiErrorResponse = {
    error: true;
    code: string;
    message: string;
};

export type SendChatMessageResult =
    | {
        ok: true;
        answer: string;
    }
    | {
        ok: false;
        status: number;
        code: string;
        message: string;
    };