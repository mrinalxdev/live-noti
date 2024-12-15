class NotificationSystem {
  constructor() {
    this.ws = null;
    this.userId = uuid.v4();
    this.channelId = null;
  }

  connect() {
    this.ws = new WebSocket("ws://localhost:8080/ws");
    this.ws.onopen = () => {
      console.log("Websocket connected");
    };

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleNotification(message);
    };

    this.ws.onclose = () => {
      console.log("WebSocket disconnected");
      setTimeout(() => this.connect(), 5000);
    };

    this.ws.onerror = (error) => {
      console.error("Websocket error :", error);
    };
  }

  joinChannel(channelId) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      alert("WebSocket not connected");
      return;
    }
    this.channelId = channelId;
    const message = {
      type: "join_channel",
      channelId: channelId,
      userId: this.userId,
    };

    this.ws.send(JSON.stringify(message));
    this.showToast(`Joined channel : ${channelId}`);
  }

  performAction(action) {
    if (!this.channelId) {
      alert("Please join a channel first");
      return;
    }
    const message = {
      type: "user_action",
      channelId: this.channelId,
      userId: this.userId,
      action: action,
    };

    this.ws.send(JSON.stringify(message));
  }

  handleNotification(message) {
    if (message.userId !== this.userId) {
      this.addNotification(message);
      this.showToast(`User ${message.userId} peformed : ${message.action}`);
    }
  }

  addNotification(message) {
    const notifications = document.getElementById("notifications");
    const notificationElement = docuement.createElement("div");
    notificationElement.className = "p-4 bg-gray-100 rounded-lg";
    notificationElement.innerHTML = `<p class="text-sm">
                <span class="font-bold">User ${message.userId}</span>
                performed: ${message.action}
            </p>
            <p class="text-xs text-gray-500 mt-1">
                ${new Date().toLocaleTimeString()}
            </p>`;
    notifications.insertBefore(notificationElement, notifications.firstChild);
  }

  showToast(message) {
    const toast = document.createElement("div");
    toast.className =
      "bg-black bg-opacity-75 text-white px-6 py-3 rounded-lg transform transition-all duration-500 opacity-0";
    toast.textContent = message;

    const container = document.getElementById("toastContainer");
    container.appendChild(toast);

    setTimeout(() => {
      toast.classList.add("opacity-100");
    }, 10);

    //logic to remove the toast
    setTimeout(() => {
      toast.classList.remove("opacity-100");
      setTimeout(() => {
        container.removeChild(toast);
      }, 500);
    }, 3000);
  }
}

const notificationSystem = new NotificationSystem();
notificationSystem.connect();

function joinChannel() {
  const channelId = document.getElementById("channelId").value;
  if (channelId) {
    notificationSystem.joinChannel(channelId);
  }
}

function performAction() {
  const action = document.getElementById("action").value;
  if (action) {
    notificationSystem.performAction(action);
    document.getElementById("action").value = "";
  }
}
