import { apiFetch, errorMessage, escapeHtml } from "../../api.js";

export function renderMessagePage(user) {
  return `
    <section class="message-page">
      <header class="message-header">
        <div>
          <h2>${escapeHtml(user.username)}</h2>
          <span class="message-user-status">Chat</span>
        </div>
      </header>

      <div id="messages-list" class="messages-list">
        <p class="messages-loading">
          Loading messages...
        </p>
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
  const message_list = document.getElementById("messages-list");

  if (!message_list) {
    console.error("Messages container not found");
    return;
  }

  message_list.textContent = "Loading messages...";

  const response = await apiFetch(`/api/messages/${userId}`);

  if (!response.ok) {
    message_list.innerHTML = `<p class="chat-error">${errorMessage(response, "Failed to load messages")}</p>`;
    return;
  }
  const messages = Array.isArray(response.data) ? response.data : [];

  renderMessages(messages, message_list);
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
                <span class="chat-message-sender">${escapeHtml(message.username)}</span>
                <span class="chat-message-content">${escapeHtml(message.content)}</span>
                <span class="chat-message-timestamp">${formatTimestamp(message.created_at)}</span>
            </div>
        `,
    )
    .join("");
}

function formatTimestamp(timestamp) {
  const date = new Date(timestamp);
  return date.toLocaleString();
}

export async function setupMessageForm(userId) {
  const form = document.getElementById("message-form");

  if (!form) {
    console.error("Message form not found");
    return;
  }

  form.addEventListener("submit", async (event) => {
    handleMessageFormSubmit(event, userId);
  });

  await loadMessagesForUser(userId);
}

async function handleMessageFormSubmit(event, userId) {
  event.preventDefault();

  const input = document.getElementById("message-input");

  const submitButton = event.target.querySelector("button[type='submit']");

  const content = input.value.trim();
  if (content === "") {
    errorMessage("Message content cannot be empty");
    return;
  }

  submitButton.disabled = true;

  const response = await apiFetch(`/api/messages/${userId}`, {
    method: "POST",
    body: {
      receiver_id: receiver_id,
      content: content,
    },
  });

  submitButton.disabled = false;

  if (!response.ok) {
    errorMessage(response, "Failed to send message");
    return;
  }

  input.value = "";

  await loadMessagesForUser(userId);
}
