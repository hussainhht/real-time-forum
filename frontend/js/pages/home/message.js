import { apiFetch, errorMessage, escapeHtml } from "../../api.js";

export function massagePage() {
  return `
  <section class="message-page">
      <header class="message-header">
        <div class="message-user-avatar"></div>

        <div>
          <h2>${escapeHtml(user.username)}</h2>
          <span class="message-user-status">Chat</span>
        </div>
      </header>

      <div id="messages-list" class="messages-list">
        <p class="messages-loading">Loading messages...</p>
      </div>

      <form id="message-form" class="message-form">
        <input
          id="message-input"
          type="text"
          placeholder="Write a message..."
          autocomplete="off"
          required
        >

        <button type="submit">
          Send
        </button>
      </form>
    </section>
    `;
}

export async function loadMessagesForUser(userId) {
  const messagesContainer = document.getElementById("messages-container");

  if (!messagesContainer) {
    console.error("Messages container not found");
    return;
  }

  messagesContainer.textContent = "Loading messages...";

  const response = await apiFetch(`/api/messages/${userId}`);

  if (!response.ok) {
    messagesContainer.innerHTML = `<p class="chat-error">${errorMessage(response, "Failed to load messages")}</p>`;
    return;
  }
  const messages = Array.isArray(response.data) ? response.data : [];

  renderMessages(messages, messagesContainer);
}


function renderMessages(messages, messagesContainer) {
  if (messages.length === 0) {
    messagesContainer.innerHTML = `<p class="empty-chat">No messages found</p>`;
    return;
  }

  messagesContainer.innerHTML = messages
    .map(
      (message) => `
            <div class="chat-message">
                <span class="chat-message-sender">${escapeHtml(message.sender_id)}</span>
                <span class="chat-message-content">${escapeHtml(message.content)}</span>
                <span class="chat-message-timestamp">${escapeHtml(message.created_at)}</span>
            </div>
        `,
    )
    .join("");
}
