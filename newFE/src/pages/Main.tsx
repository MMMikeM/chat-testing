import { FormEvent, useState } from "react"
import { Button, Flex, FormControl, Input, Spacer } from "@chakra-ui/react"
import useWebsocket from "@/hooks/useWebsocket"
import { useParams } from "react-router"
import MessageList from "./MessageList"

export type Message = {
  body: string
  conversation_id: string
  created_at: string
  from_user_id: string
  uuid: string
}

const Main = () => {
  return (
    <>
      <h1>Main</h1>
      <Form />
      <MessageList />
    </>
  )
}

export default Main

const Form = () => {
  const { userId, conversationId } = useParams()
  const [messageField, setMessageField] = useState("")
  const { sendMessage } = useWebsocket()

  const handleSubmit = (e: FormEvent<HTMLFormElement>, message: string) => {
    e.preventDefault()
    if (message.length > 0) {
      const msg = {
        body: message,
        from_user_id: userId,
        conversation_id: conversationId,
      }
      sendMessage(JSON.stringify(msg))
      setMessageField((m) => (parseInt(m, 10) + 1).toString())
    }
  }

  const copyConversationId = () => {
    if (conversationId) {
      navigator.clipboard.writeText(conversationId.toString())
    }
  }

  return (
    <Flex mt={2} direction={"column"}>
      <div>
        <Spacer />
        <Button onClick={() => copyConversationId()} marginBottom="4">
          Copy conversation ID
        </Button>
      </div>

      <form onSubmit={(e) => handleSubmit(e, messageField)}>
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
    </Flex>
  )
}
