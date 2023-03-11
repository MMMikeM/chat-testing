import { FormEvent, useState } from "react"
import { Box, Button, Flex, FormControl, Input, Spacer, Text } from "@chakra-ui/react"
import useWebsocket from "../hooks/useWebsocket"

type Messages = {
  id: number
  body: string
  from_user_id: number
  conversation_id: number
}

const Main = () => {
  const [messages, setMessages] = useState<Messages[]>([])
  const [messageField, setMessageField] = useState("")
  const { userId, conversationId, ws } = useWebsocket()

  const copyConversationId = () => {
    if (conversationId) {
      navigator.clipboard.writeText(conversationId.toString())
    }
  }

  ws.onmessage = (event) => {
    const currentMessages = [...messages, JSON.parse(event.data)]
    if (currentMessages.length > 20) {
      currentMessages.shift()
    }
    setMessages(currentMessages)
  }

  const sendMessage = (e: FormEvent<HTMLFormElement> | undefined) => {
    e?.preventDefault()
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
    <>
      <h1>Main</h1>
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
      <Flex mt={2}>
        <Spacer />
        <Button onClick={() => copyConversationId()} marginBottom="4">
          Copy conversation ID
        </Button>
      </Flex>
      <form onSubmit={sendMessage}>
        <FormControl mt={4}>
          <Input
            type="text"
            placeholder="Type a message..."
            value={messageField}
            onChange={(e) => setMessageField(e.target.value)}
          />
        </FormControl>
        <Flex mt={2}>
          <Spacer />
          <Button colorScheme="purple" type="submit">
            Send Message
          </Button>
        </Flex>
      </form>
    </>
  )
}

export default Main
