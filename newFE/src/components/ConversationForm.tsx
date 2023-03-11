import { useState } from "react"
import { Box, Button, FormControl, Input } from "@chakra-ui/react"
import { createConversation } from "../api/conversation"
import useUpdateUrl from "../hooks/useUpdateUrl"
import useWebSocket from "../hooks/useWebsocket"

const ConversationForm = () => {
  const [input, setInput] = useState("")
  const { conversationId } = useWebSocket()
  const updateUrl = useUpdateUrl({ conversationId })

  const newConversation = async () => {
    const conversation = await createConversation()
    if (conversation?.id) {
      updateUrl({ newConversationId: conversation.id })
    }
  }

  const joinConversation = () => {
    updateUrl({ newConversationId: input })
  }

  return (
    <Box mt={4} p={4}>
      <FormControl>
        <Input
          placeholder="Paste conversation ID here..."
          value={input}
          type="number"
          onChange={(event) => setInput(event.target.value)}
          key="conversationId"
          onKeyPress={(event) => (event.key === "Enter" ? updateUrl({ newConversationId: input }) : null)}
          autoFocus
        />
      </FormControl>
      <Button mt={4} colorScheme="purple" onClick={() => joinConversation()}>
        Join conversation
      </Button>
      <Button mt={4} ml={4} onClick={() => newConversation()}>
        Create conversation
      </Button>
    </Box>
  )
}

export default ConversationForm
