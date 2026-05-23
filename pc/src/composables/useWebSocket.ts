import { ref, onUnmounted } from "vue"

export interface WsMessage {
  id?: number
  sender_id: number
  receiver_id: number
  content: string
  msg_type: string
  created_at?: string
}

export interface IncomingMessage {
  type: string
  data?: WsMessage
  sender_id?: number
}

export function useWebSocket(userId: number) {
  let ws: WebSocket | null = null
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let pingTimer: ReturnType<typeof setInterval> | null = null

  const connected = ref(false)
  const lastMessage = ref<IncomingMessage | null>(null)

  function connect() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"
    const host = window.location.host
    const token = localStorage.getItem("token") || ""
    const url = `${protocol}//${host}/api/v1/ws?token=${token}`

    ws = new WebSocket(url)

    ws.onopen = () => {
      connected.value = true
      pingTimer = setInterval(() => {
        ws?.send(JSON.stringify({ type: "ping" }))
      }, 25000)
    }

    ws.onclose = () => {
      connected.value = false
      if (pingTimer) clearInterval(pingTimer)
      reconnectTimer = setTimeout(connect, 3000)
    }

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data) as IncomingMessage
      lastMessage.value = data
    }

    ws.onerror = () => {
      connected.value = false
    }
  }

  function sendMessage(receiverId: number, content: string, msgType = "text") {
    if (!ws || ws.readyState !== WebSocket.OPEN) return
    ws.send(JSON.stringify({
      type: "message",
      receiver_id: receiverId,
      content,
      msg_type: msgType,
    }))
  }

  function sendTyping(receiverId: number) {
    if (!ws || ws.readyState !== WebSocket.OPEN) return
    ws.send(JSON.stringify({
      type: "typing",
      receiver_id: receiverId,
    }))
  }

  function disconnect() {
    if (reconnectTimer) clearTimeout(reconnectTimer)
    if (pingTimer) clearInterval(pingTimer)
    if (ws) {
      ws.onclose = null
      ws.close()
      ws = null
    }
    connected.value = false
  }

  onUnmounted(disconnect)

  return { connect, disconnect, sendMessage, sendTyping, connected, lastMessage }
}
