import { useState } from 'react'
import {
  ChakraProvider,
  Container,
  Box,
  Button,
  Flex,
  Spacer,
  ButtonGroup,
  FormControl,
  FormLabel,
  FormErrorMessage,
  FormHelperText,
  Input,
  Text
} from '@chakra-ui/react'

import { createConversation } from '../api/conversation';
import { createUser } from '../api/user';

export const Main = ({ws, userId, conversationId, setConversationId}) => {
  const [message, setMessage] = useState("")
  const [username, setUsername] = useState("")
  const [messages, setMessages] = useState([])

  const convo = async (userId) => {
    const conversation = await createConversation()
    if (conversation.uuid !== undefined) {
      window.location.href = `http://localhost:3000?userId=${userId}&conversationId=${conversation.uuid}`;
    }
  }

  const user = async () => {
    const user = await createUser(username)
    if (user.uuid !== undefined) {
      window.location.href = `http://localhost:3000?userId=${user.uuid}`;
    }
  }

  const joinConversation = () => {
    if (conversationId !== undefined) {
      window.location.href = `http://localhost:3000?userId=${userId}&conversationId=${conversationId}`;
    }
  }

  const copyConversationId = () => {
    navigator.clipboard.writeText(conversationId)
  }

  if (userId !== null && conversationId !== null) {
    if (ws) {
      ws.onmessage = function (event) {
        let currentMessages = [...messages, JSON.parse(event.data)]
        if (currentMessages.length > 20) {
          currentMessages.shift()
        }
        setMessages(currentMessages)
      };
    }

    const sendMessage = () => {
      if (message.length > 0) {
        let msg = {body: message, from_user_id: userId, conversation_id: conversationId}
        if (ws) {
          ws.send(JSON.stringify(msg))
        }
        setMessage("")
      }
    }

    const handleSubmit = (event) => {
      event.preventDefault();
      sendMessage()
    }

    return (
      <div>
        {messages.length > 0 ? 
        <Box mt={4}>
          {messages.map((message) => <p key={message.uuid}>{message.from_user_id}: {message.body}</p>)}
        </Box>
        :
        <Box mt={4}>
          <Text colorScheme="gray">No messages exist yet...</Text>
        </Box>
        }
        <form onSubmit={handleSubmit}>
          <FormControl mt={4}>
            <Input
              type="text"
              placeholder="Type a message..."
              value={message}
              onInput={(e) => setMessage(e.target.value)}
            />
          </FormControl>
        </form>
        <Flex mt={2}>
          <Button onClick={ () => copyConversationId() }>Copy conversation ID</Button>
          <Spacer />
          <Button colorScheme='purple' onClick={ () => sendMessage() }>Send Message</Button>
        </Flex>
      </div>
    )
  }

  if (userId !== null && conversationId === null) {
    return (
      <Box mt={4} p={4}>
        <FormControl>
          <Input
            placeholder="Paste conversation ID here..."
            value={ conversationId }
            onChange={ (event) => setConversationId(event.target.value) }
            key="conversationId"
            onKeyPress={(event) => event.key === 'Enter' ? joinConversation() : ""}
            autoFocus
          />
        </FormControl>
        <Button mt={4} colorScheme='purple' onClick={ () => joinConversation() }>Join conversation</Button>
        <Button mt={4} ml={4} onClick={ () => convo(userId) }>Create conversation</Button>
      </Box>
    )
  }

  return (
    <Box mt={4} p={4}>
      <FormControl>
        <Input
          placeholder="Enter username..."
          value={ username }
          onInput={ (event) => setUsername(event.target.value) }
          onKeyPress={(event) => event.key === 'Enter' ? user() : ""}
          key="username"
          autoFocus
        />
      </FormControl>
      <Button colorScheme='purple' mt={4} onClick={ () => user() }>Create user</Button>
    </Box>
  )
}
