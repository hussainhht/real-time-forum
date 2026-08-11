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

  console.log("2- socket created:", socket);

  socket.onopen = () => {
    console.log("WebSocket is connected now");
    // socket.send("hello from browser");
  };

  socket.onmessage = (event) => {
    const data = JSON.parse(event.data);

    console.log("Websocket event :", data);

    switch (data.type) {
      case "user_online":
        console.log("User online :", data.content);
        break;
      case "user_offline":
        console.log("User offline:", data.content);
        break;
      case "new_message":
        console.log("New message:", data.content);
        break;
    }

    console.log("Message from server:", data);
  };

  socket.onclose = () => {
    console.log("Websockent close now");
  };

  socket.onerror = (error) => {
    console.error("WebSocket error :", error);
  };
}