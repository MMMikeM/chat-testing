import { useParams } from "react-router-dom"

const useWebsocket = () => {
  const params = useParams()
  const conversationId = params.conversationId ? parseInt(params.conversationId) : null
  const userId = params.userId ? parseInt(params.userId) : null

  const ws = new WebSocket(`ws://localhost:3001/ws?conversation_id=${conversationId}`)
  ws.onopen = function () {
    console.log("Connected")
  }

  return { ws, userId, conversationId }
}

export default useWebsocket
