let unauthorizedHandler = null;

export function setUnauthorizedHandler(handler) {
    unauthorizedHandler = handler;
}

// Kept in sync with the backend's own limits (post/comment content <10,000
// chars, message content <1,000 chars) so the frontend rejects oversized
// input before it hits the network instead of relying on the server 400.
const SIZE_LIMITS = {
    request: {
        auth: 2_000,
        post: 15_000,
        comment: 15_000,
        message: 2_000,
        default: 10_000,
    },
    response: {
        single: 100_000,
        feed: 1_000_000,
        default: 500_000,
    },
    warning: 50_000,
};

function byteSize(str) {
    return new TextEncoder().encode(str).length;
}

function formatBytes(bytes) {
    if (bytes === 0) return "0 B";
    const units = ["B", "KB", "MB", "GB"];
    const exponent = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
    return `${(bytes / Math.pow(1024, exponent)).toFixed(exponent === 0 ? 0 : 2)} ${units[exponent]}`;
}

function requestLimitFor(url, method) {
    if (method !== "POST") return SIZE_LIMITS.request.default;
    if (url.includes("/comments")) return SIZE_LIMITS.request.comment;
    if (url.includes("/messages")) return SIZE_LIMITS.request.message;
    if (url === "/register" || url === "/login") return SIZE_LIMITS.request.auth;
    if (url.startsWith("/posts")) return SIZE_LIMITS.request.post;
    return SIZE_LIMITS.request.default;
}

function responseLimitFor(url, method) {
    if (method === "GET") {
        if (url === "/posts") return SIZE_LIMITS.response.feed;
        if (/^\/posts\/[^/]+\/comments/.test(url)) return SIZE_LIMITS.response.feed;
        if (url.startsWith("/api/messages/")) return SIZE_LIMITS.response.feed;
        if (url === "/api/chat-users") return SIZE_LIMITS.response.feed;
    }
    return SIZE_LIMITS.response.single;
}

function sizeError(message) {
    return { ok: false, status: 413, error: "Payload too large", data: { error: message } };
}

function checkRequestSize(url, method, body) {
    if (!body) return null;
    const size = byteSize(body);
    const limit = requestLimitFor(url, method);
    if (size <= limit) return null;
    console.warn(`[SIZE] blocked request to ${url}: ${formatBytes(size)} exceeds limit ${formatBytes(limit)}`);
    return sizeError(`Your request is too large (${formatBytes(size)}). Maximum allowed is ${formatBytes(limit)}.`);
}

function checkResponseSize(url, method, text) {
    if (!text) return null;
    const size = byteSize(text);
    const limit = responseLimitFor(url, method);
    if (size > SIZE_LIMITS.warning) {
        console.warn(`[SIZE] large response from ${url}: ${formatBytes(size)}`);
    }
    if (size <= limit) return null;
    console.error(`[SIZE] blocked response from ${url}: ${formatBytes(size)} exceeds limit ${formatBytes(limit)}`);
    return sizeError(`The server returned too much data (${formatBytes(size)}). Please try again.`);
}

export async function apiFetch(url, options = {}) {
    let { method = "GET", headers = {}, body } = options;
    if (body && typeof body === "object") {
        body = JSON.stringify(body);
        headers["Content-Type"] = "application/json";
    }

    const requestSizeError = checkRequestSize(url, method, body);
    if (requestSizeError) return requestSizeError;

    let response;
    try {
        response = await fetch(url, {
            method,
            credentials: "same-origin",
            headers,
            body
        });
    } catch (error) {
        return { ok: false, status: 0, error: "Network error", data: null };
    }

    let text = "";
    try {
        text = await response.text();
    } catch (error) {
        text = "";
    }

    const responseSizeError = checkResponseSize(url, method, text);
    if (responseSizeError) return responseSizeError;

    let data = null;
    try {
        data = text ? JSON.parse(text) : null;
    } catch (error) {
        data = null;
    }

    // /api/session is expected to return 401 for a logged-out visitor on
    // first load, that's not a "you got kicked out" event, so it's excluded.
    if (response.status === 401 && url !== "/api/session" && unauthorizedHandler) {
        unauthorizedHandler();
    }

    return { ok: response.ok, status: response.status, error: null, data };
}


export function errorMessage(result, fallbackMessage = "An error occurred") {
    if (result.error === "Network error") {
        fallbackMessage = "Network error. Please check your connection.";
    }
    if (result.data && typeof result.data.error === "string") {
        fallbackMessage = result.data.error;
    }
    return fallbackMessage;
}


export function escapeHtml(value) {
    const div = document.createElement("div");
    div.textContent = value ? value : "";
    return div.innerHTML;
}
//todo: check it when we run the code in the browser, it should be tested to ensure that it correctly escapes HTML and prevents XSS attacks.