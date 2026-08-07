export let socket;

export function connectWebsocket() {
  socket = new WebSocket(`ws://${window.location.host}/ws`);

  console.log("2- socket created:", socket);

  socket.onopen = () => {
    console.log("WebSocket is connected now");
    // socket.send("hello from browser");
  };

  socket.onmessage = (event) => {
    console.log("Message from server:", event.data);
  };

  socket.onclose = () => {
    console.log("Websockent close now");
  };

  socket.onerror = (error) => {
    console.error("WebSocket error :", error);
  };
}
