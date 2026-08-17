import { apiFetch, errorMessage, escapeHtml } from "../../api.js";
import { currentUser } from "../../auth.js";
import { renderErrorPage } from "../error.js";
import { app } from "../../router.js";

let oldestMessageId = 0;
let isLoadingOlder = false;
let hasMoreMessages = true;

export function renderMessagePage(user) {
  return `
    <section class="message-page">
      <header class="message-header">
        <div>
          <h2>${escapeHtml(user.username)}</h2>
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

export async function loadMessagesForUser(userId, beforeId = 0) {
  const messagesContainer = document.getElementById("messages-list");

  if (!messagesContainer) {
    console.error("Messages container not found");
    return;
  }

  if (beforeId === 0) {
    messagesContainer.textContent = "Loading messages...";
  }
  const response = await apiFetch(
    `/api/messages/${userId}?before_id=${beforeId}`,
  );

  if (!document.body.contains(messagesContainer)) return;

  if (response.status === 500) {
    renderErrorPage(
      app,
      "500 Internal Server Error",
      "Something went wrong on our end. Please try again later."
    );
    return;
  }

  if (!response.ok) {
    messagesContainer.innerHTML = `<p class="chat-error">${errorMessage(response, "Failed to load messages")}</p>`;
    return;
  }

  const messages = Array.isArray(response.data) ? response.data : [];
  //response.data should contain the []structures.Messages between these two users
  
  if (beforeId === 0) {
    renderMessages(messages, messagesContainer);
    if (messages.length > 0) {
      oldestMessageId = messages[0].id;
    }
    hasMoreMessages = messages.length === 10;
    return;
  }

  prependMessages(messages, messagesContainer);

  if (messages.length > 0) {
    oldestMessageId = messages[0].id;
  }

  if (messages.length < 10) {
    hasMoreMessages = false;
  }
}

async function loadOlderMessages(userId) {
  if (isLoadingOlder) return;
  if (!hasMoreMessages) return;
  if (oldestMessageId === 0) return;

  isLoadingOlder = true;

  try {
    await loadMessagesForUser(userId, oldestMessageId);
  } finally {
    isLoadingOlder = false;
  }
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
  const isOwn = currentUser && message.sender_id === currentUser.id;

  return `
    <div class="chat-message" data-own="${isOwn}">
      <span class="chat-message-sender">
        ${escapeHtml(message.username)}
      </span>

      <p class="chat-message-content">
        ${escapeHtml(message.content)}
        <span class="chat-message-timestamp">
          ${formatTimestamp(message.created_at)}
        </span>
      </p>
    </div>
  `;
}

function formatTimestamp(timestamp) {
  const date = new Date(timestamp);
  return date.toLocaleString();
}

export async function setupMessageForm(userId) {
  const form = document.getElementById("message-form");
  const messagesContainer = document.getElementById("messages-list");

  if (!form || !messagesContainer) {
    console.error("Message form or messages container not found");
    return;
  }

  oldestMessageId = 0;
  isLoadingOlder = false;
  hasMoreMessages = true;

  form.addEventListener("submit", (event) => {
    handleMessageFormSubmit(event, userId);
  });

  await loadMessagesForUser(userId);

  messagesContainer.addEventListener(
    "scroll",
    throttle(() => {
      if (messagesContainer.scrollTop <= 50) {
        loadOlderMessages(userId);
      }
    }, 200),
  );
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

  const MAX_MESSAGE_SIZE = 10 * 1024;
  const body = {
    receiver_id: userId,
    content: content,
  };

  if (JSON.stringify(body).length > MAX_MESSAGE_SIZE) {
    if (message) message.textContent = "Message is too large (max 10 KB)";
    return;
  }

  submitButton.disabled = true;

  const response = await apiFetch(`/api/messages`, {
    method: "POST",
    body: body,
  });

  submitButton.disabled = false;

  if (!document.body.contains(input)) return;

  if (response.status === 500) {
    renderErrorPage(
      app,
      "500 Internal Server Error",
      "Something went wrong on our end. Please try again later."
    );
    return;
  }

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

function prependMessages(messages, messagesContainer) {
  if (messages.length === 0) {
    hasMoreMessages = false;
    return;
  }

  const oldScrollHeight = messagesContainer.scrollHeight;
  const oldScrollTop = messagesContainer.scrollTop;

  const oldMessagesHTML = messages
    .map((message) => buildMessageBubbleHtml(message))
    .join("");

  messagesContainer.insertAdjacentHTML("afterbegin", oldMessagesHTML);

  const newScrollHeight = messagesContainer.scrollHeight;

  messagesContainer.scrollTop =
    oldScrollTop + (newScrollHeight - oldScrollHeight);
}

function throttle(callback, delay) {
  let waiting = false;

  return (...args) => {
    if (waiting) return;

    callback(...args);

    waiting = true;

    setTimeout(() => {
      waiting = false;
    }, delay);
  };
}
