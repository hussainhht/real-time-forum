import { appendMessage } from "./pages/home/message.js";
import { apiFetch } from "./api.js";

export let socket;

let reconnectTimeout = null;
let intentionalClose = false;

export function connectWebsocket() {
  if (
    socket &&
    (socket.readyState === WebSocket.OPEN ||
      socket.readyState === WebSocket.CONNECTING)
  ) {
    console.log(
      "[ws] connectWebsocket() skipped, existing socket readyState =",
      socket.readyState,
    );
    return;
  }

  intentionalClose = false;
  console.log("[ws] opening new socket...");
  socket = new WebSocket(`ws://${window.location.host}/ws`);

  socket.onopen = () => {
    console.log("[ws] OPEN at", new Date().toISOString());
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log("[ws] MESSAGE at", new Date().toISOString(), data);

    switch (data.type) {
      case "user_online":
        updateUserOnlineDot(data.content.user_id, true);
        break;
      case "user_offline":
        updateUserOnlineDot(data.content.user_id, false);
        break;
      case "new_message":
        handleIncomingMessage(data.content);
        break;
    }
  };

  socket.onclose = (event) => {
    console.log(
      "[ws] CLOSE at",
      new Date().toISOString(),
      "code:",
      event.code,
      "reason:",
      event.reason,
      "wasClean:",
      event.wasClean,
      "intentionalClose:",
      intentionalClose,
    );

    if (intentionalClose) return;

    // kickout. If the session is still valid, just reconnect.
    apiFetch("/api/session").then((result) => {
      if (result.ok && result.data && result.data.authenticated) {
        console.log("[ws] session still valid, reconnecting in 2s");
        reconnectTimeout = setTimeout(connectWebsocket, 2000);
      }
      // if not authenticated, apiFetch's global 401 handling only fires
      // for non-session endpoints, so we redirect explicitly here.
      else if (!result.ok) {
        console.log("[ws] session invalid, reloading page");
        window.location.reload();
      }
    });
  };

  socket.onerror = (error) => {
    console.error("[ws] ERROR at", new Date().toISOString(), error);
  };
}

export function disconnectWebsocket() {
  console.log(
    "disconnectWebsocket called. socket:",
    socket,
    "readyState:",
    socket ? socket.readyState : "no socket",
  );

  intentionalClose = true;
  if (reconnectTimeout) {
    clearTimeout(reconnectTimeout);
    reconnectTimeout = null;
  }
  if (socket) {
    socket.close();
  }
}

function updateUserOnlineDot(userId, isOnline) {
  const buttons = document.getElementsByClassName("chat-user-btn");
  console.log(
    "[ws] updateUserOnlineDot called for userId=",
    userId,
    "isOnline=",
    isOnline,
    "buttons found in DOM:",
    buttons.length,
  );

  let userButton = null;
  for (let i = 0; i < buttons.length; i++) {
    if (Number(buttons[i].dataset.userId) === Number(userId)) {
      userButton = buttons[i];
      break;
    }
  }
  if (!userButton) {
    console.warn("[ws] no button found for userId=", userId);
    return;
  }

  userButton.dataset.online = isOnline ? "true" : "false";

  const dots = userButton.getElementsByClassName("chat-user-circle");
  if (dots.length === 0) return;

  const dot = dots[0];

  if (isOnline) {
    dot.classList.remove("offline");
    dot.classList.add("online");
  } else {
    dot.classList.remove("online");
    dot.classList.add("offline");
  }

  const messagesList = document.getElementById("messages-list");
  if (messagesList && Number(messagesList.dataset.userId) === Number(userId)) {
    const statusLabel = document.querySelector(".message-user-status");
    if (statusLabel) {
      statusLabel.textContent = isOnline ? "Online" : "Offline";
    }
  }
}

function handleIncomingMessage(message) {
  const messagesList = document.getElementById("messages-list");
  if (!messagesList) return;

  const openChatUserId = Number(messagesList.dataset.userId);
  const senderId = Number(message.sender_id);

  if (openChatUserId !== senderId) return;

  appendMessage(message);
}
