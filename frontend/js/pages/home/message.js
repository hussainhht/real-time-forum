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

      <div id="messages-list" class="messages-list" data-user-id="${escapeHtml(user.id)}">
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
      <p id="chat-form-message" class="form-message"></p>
    </section>
  `;
}

export async function loadMessagesForUser(userId) {
  const messagesContainer = document.getElementById("messages-list");

  if (!messagesContainer) {
    console.error("Messages container not found");
    return;
  }

  messagesContainer.textContent = "Loading messages...";

  const response = await apiFetch(`/api/messages/${userId}`);

  if (!document.body.contains(messagesContainer)) return;

  if (!response.ok) {
    messagesContainer.innerHTML = `<p class="chat-error">${errorMessage(response, "Failed to load messages")}</p>`;
    return;
  }
  const messages = Array.isArray(response.data) ? response.data : [];

  renderMessages(messages, messagesContainer);
}

export function renderMessages(messages, messagesContainer) {
  if (messages.length === 0) {
    messagesContainer.innerHTML = `<p class="empty-chat">No messages found</p>`;
    return;
  }

  messagesContainer.innerHTML = messages
    .map((message) => buildMessageBubbleHtml(message))
    .join("");

  messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

function buildMessageBubbleHtml(message) {
  return `
    <div class="chat-message">
      <span class="chat-message-sender">
        ${escapeHtml(message.username)}
      </span>

      <span class="chat-message-content">
        ${escapeHtml(message.content)}
      </span>

      <span class="chat-message-timestamp">
        ${formatTimestamp(message.created_at)}
      </span>
    </div>
  `;
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

  form.addEventListener("submit", (event) => {
    handleMessageFormSubmit(event, userId);
  });

  await loadMessagesForUser(userId);
}

async function handleMessageFormSubmit(event, userId) {
  event.preventDefault();

  const input = document.getElementById("message-input");
  const message = document.getElementById("chat-form-message");

  const submitButton = event.target.querySelector("button[type='submit']");

  const content = input.value.trim();
  if (content === "") {
    return;
  }

  submitButton.disabled = true;

  const response = await apiFetch(`/api/messages`, {
    method: "POST",
    body: {
      receiver_id: userId,
      content: content,
    },
  });

  submitButton.disabled = false;

  if (!document.body.contains(input)) return;

  if (!response.ok) {
    if (message)
      message.textContent = errorMessage(response, "Failed to send message");
    return;
  }

  input.value = "";

  const messagesContainer = document.getElementById("messages-list");
  if (messagesContainer) {
    appendMessage(response.data);
  }
}


export function appendMessage(message) {
  const messagesContainer = document.getElementById("messages-list");
  if (!messagesContainer) return;

  const emptyState = messagesContainer.querySelector(".empty-chat");
  if (emptyState) emptyState.remove();

  messagesContainer.insertAdjacentHTML(
    "beforeend",
    buildMessageBubbleHtml(message),
  );

  messagesContainer.scrollTop = messagesContainer.scrollHeight;
}