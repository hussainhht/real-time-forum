export async function apiFetch(url, options = {}) {
    let { method = "GET", headers = {}, body } = options;
    if (body && typeof body === "object") {
        body = JSON.stringify(body);
        headers["Content-Type"] = "application/json";
    }
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

    let data;
    try {
        const text = await response.text();
        data = text ? JSON.parse(text) : null;
    } catch (error) {
        data = null;
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
//todo check it when we run the code in the browser, it should be tested to ensure that it correctly escapes HTML and prevents XSS attacks.


