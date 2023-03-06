import { useState, useEffect } from 'react'
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

import { Main  } from './components/main';

const params = new Proxy(new URLSearchParams(window.location.search), {
  get: (searchParams, prop) => searchParams.get(prop),
});

const ws = new WebSocket(`ws://localhost:3001/ws?conversation_id=${params.conversationId}`);
ws.onopen = function() {
  console.log('Connected')
}

function App() {
  const [userId, setUserId] = useState("")
  const [conversationId, setConversationId] = useState("")

  const MainComponentLayout = ({ children }) => {
    return (
      <Container maxW='6xl'>
        {children}
      </Container>
    )
  }

  useEffect(() => {
    setUserId(params.userId)
    setConversationId(params.conversationId)
  }, [])


  return (
    <ChakraProvider>
      <header></header>
      <main>
        <MainComponentLayout>
          <Main
            ws={ws}
            userId={userId}
            conversationId={conversationId}
          />
        </MainComponentLayout>
      </main>
      <footer></footer>
    </ChakraProvider>
  );
}

export default App;
