import { useState } from "react"
import { Box, Button, Flex, FormControl, Input, Spacer, Text } from "@chakra-ui/react"
import useWebsocket from "../hooks/useWebsocket"
import UserForm from "./UserForm"
import ConversationForm from "./ConversationForm"

type Messages = {
  id: number
  body: string
  from_user_id: number
  conversation_id: number
}

export const Main = () => {
  const [messages, setMessages] = useState<Messages[]>([])
  const [messageField, setMessageField] = useState("")
  const { userId, conversationId, ws } = useWebsocket()

  const copyConversationId = () => {
    if (conversationId) {
      navigator.clipboard.writeText(conversationId.toString())
    }
  }

  if (!userId) return <UserForm userId={userId} />

  if (!conversationId) {
    return <ConversationForm conversationId={conversationId} />
  }

  ws.onmessage = function (event) {
    const currentMessages = [...messages, JSON.parse(event.data)]
    if (currentMessages.length > 20) {
      currentMessages.shift()
    }
    setMessages(currentMessages)
  }

  const sendMessage = () => {
    if (messageField.length > 0) {
      const msg = {
        body: messageField,
        from_user_id: userId,
        conversation_id: conversationId,
      }
      ws.send(JSON.stringify(msg))
      setMessageField("")
    }
  }

  return (
    <div>
      {messages.length > 0 ? (
        <Box mt={4}>
          {messages.map((message) => (
            <p key={message.id}>
              {message.from_user_id}: {message.body}
            </p>
          ))}
        </Box>
      ) : (
        <Box mt={4}>
          <Text colorScheme="gray">No messages exist yet...</Text>
        </Box>
      )}
      <form onSubmit={() => sendMessage()}>
        <FormControl mt={4}>
          <Input
            type="text"
            placeholder="Type a message..."
            value={messageField}
            onChange={(e) => setMessageField(e.target.value)}
          />
        </FormControl>
      </form>
      <Flex mt={2}>
        <Button onClick={() => copyConversationId()}>Copy conversation ID</Button>
        <Spacer />
        <Button colorScheme="purple" onClick={() => sendMessage()}>
          Send Message
        </Button>
      </Flex>
    </div>
  )
}
