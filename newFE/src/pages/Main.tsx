import { FormEvent, useState } from "react"
import { Box, Button, Flex, FormControl, Input, Spacer, Text } from "@chakra-ui/react"
import useWebsocket from "@/hooks/useWebsocket"

type Messages = {
  body: string
  conversation_id: string
  created_at: string
  from_user_id: string
  uuid: string
}

const Main = () => {
  const [messages, setMessages] = useState<Messages[]>([])
  const [messageField, setMessageField] = useState("")
  const { userId, conversationId, ws } = useWebsocket()

  if (ws === undefined) return <>Loading</>

  const copyConversationId = () => {
    if (conversationId) {
      navigator.clipboard.writeText(conversationId.toString())
    }
  }

  ws.onmessage = (event) => {
    setMessages((m) => [...m, JSON.parse(event.data)])
  }

  const sendMessage = (e: FormEvent<HTMLFormElement>, message: string) => {
    e.preventDefault()
    if (message.length > 0) {
      const msg = {
        body: message,
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

      <Flex mt={2}>
        <Spacer />
        <Button onClick={() => copyConversationId()} marginBottom="4">
          Copy conversation ID
        </Button>
      </Flex>
      <form onSubmit={(e) => sendMessage(e, messageField)}>
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
      {messages.length > 0 ? (
        <Box mt={4}>
          {messages.map((message) => (
            <p key={message.uuid}>
              {message.from_user_id}: {message.body}
            </p>
          ))}
        </Box>
      ) : (
        <Box mt={4}>
          <Text colorScheme="gray">No messages exist yet...</Text>
        </Box>
      )}
    </>
  )
}

export default Main
