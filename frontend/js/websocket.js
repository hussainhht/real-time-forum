import { appendMessage } from "./pages/home/message.js";

export let socket;

export function connectWebsocket() {
  if (
    socket &&
    (socket.readyState === WebSocket.OPEN ||
      socket.readyState === WebSocket.CONNECTING)
  ) {
    return;
  }

  socket = new WebSocket(`ws://${window.location.host}/ws`);

  socket.onopen = () => {
    console.log("WebSocket is connected now");
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

    console.log("Websocket event :", data);

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

  socket.onclose = () => {
    console.log("WebSocket closed now");
  };

  socket.onerror = (error) => {
    console.error("WebSocket error :", error);
  };
}

function updateUserOnlineDot(userId, isOnline) {
  const buttons = document.getElementsByClassName("chat-user-btn");

  let userButton = null;
  for (let i = 0; i < buttons.length; i++) {
    if (Number(buttons[i].dataset.userId) === Number(userId)) {
      userButton = buttons[i];
      break;
    }
  }
  if (!userButton) return;

  const dots = userButton.getElementsByClassName("online-user-circle");
  if (dots.length === 0) return;

  const dot = dots[0];

  if (isOnline) {
    dot.classList.remove("offline");
    dot.classList.add("online");
  } else {
    dot.classList.remove("online");
    dot.classList.add("offline");
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