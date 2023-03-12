import { memo, useEffect, useMemo, useRef, useState } from "react"
import { useParams } from "react-router-dom"
import { Message } from "@/pages/Main"

const checkIfShouldReconnect = (ws: WebSocket | undefined, conversationId?: string) => {
  const isSameUrl = ws?.url.includes(conversationId ?? "")
  const isConnecting = ws?.CONNECTING
  const isOpen = ws?.OPEN

  if (!ws) return false
  if (!isSameUrl) {
    ws.close()
    return false
  }
  if (isSameUrl || isConnecting || isOpen) {
    return true
  }
}

const useWebsocket = () => {
  const { conversationId } = useParams()
  const activeWs = useRef<WebSocket>()
  const unmountedRef = useRef(false)
  const [latestMessage, setLatestMessage] = useState<Message>()

  const connect = () => {
    if (checkIfShouldReconnect(activeWs.current, conversationId)) return

    const ws = new WebSocket(`ws://localhost:3001/ws?conversation_id=${conversationId}`)

    ws.onerror = (error) => {
      if (unmountedRef.current) return
      console.log("error", error)
    }

    ws.onopen = () => {
      if (unmountedRef.current) return
      console.log("connected")
    }
    ws.onmessage = (e) => {
      if (unmountedRef.current) return
      const data: Message = JSON.parse(e.data)
      setLatestMessage(data)
      return data
    }
    ws.onclose = () => {
      if (unmountedRef.current) return
      console.log("disconnected")
      setTimeout(() => {
        connect()
      }, 3000)
    }

    activeWs.current = ws
  }

  const sendMessage: WebSocket["send"] = useMemo(
    () => (message) => {
      try {
        activeWs.current?.send(message)
      } catch (error) {
        throw new Error("WebSocket disconnected")
      }
    },
    [activeWs.current]
  )

  useEffect(() => {
    connect()
  }, [conversationId])

  useEffect(
    () => () => {
      // unmountedRef.current = true
      // checkIfShouldReconnect(activeWs.current, conversationId)
    },
    []
  )

  return { sendMessage, latestMessage }
}

export default useWebsocket
