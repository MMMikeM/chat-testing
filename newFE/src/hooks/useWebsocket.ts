import { useEffect, useState } from "react"
import { useParams } from "react-router-dom"

const useWebsocket = () => {
  const { conversationId, userId } = useParams()
  const [ws, setWs] = useState<WebSocket>(
    () => new WebSocket(`ws://localhost:3001/ws?conversation_id=${conversationId}`)
  )

  useEffect(() => {
    if (!ws || (conversationId && ws?.url.includes(conversationId))) {
      setWs(new WebSocket(`ws://localhost:3001/ws?conversation_id=${conversationId}`))
    }
    return () => {
      ws?.close()
    }
  }, [conversationId])

  return { ws, userId, conversationId }
}

export default useWebsocket
