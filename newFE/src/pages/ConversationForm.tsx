import { useState } from "react"
import { Box, Button, FormControl, Input } from "@chakra-ui/react"
import { createConversation } from "@/api/conversation"
import useUpdateUrl from "@/hooks/useUpdateUrl"

const ConversationForm = () => {
  const [input, setInput] = useState("")
  const updateUrl = useUpdateUrl()

  const newConversation = async () => {
    const conversation = await createConversation()
    console.log("conversation:", conversation)

    if (conversation?.ID !== undefined) {
      updateUrl({ newConversationId: conversation.ID.toString() })
    }
  }

  const joinConversation = () => {
    updateUrl({ newConversationId: input })
  }

  return (
    <Box mt={4} p={4}>
      <h1>Conversation Form</h1>
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
