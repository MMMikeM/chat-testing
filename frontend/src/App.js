import { useState, useEffect } from 'react'
import { createConversation } from './api/conversation';
import { createUser } from './api/user';

const params = new Proxy(new URLSearchParams(window.location.search), {
  get: (searchParams, prop) => searchParams.get(prop),
});

const ws = new WebSocket(`ws://localhost:3001/ws?conversation_id=${params.conversationId}`);
ws.onopen = function() {
  console.log('Connected')
}

function App() {
  const [userId, setUserId] = useState(null)
  const [conversationId, setConversationId] = useState(null)
  const [username, setUsername] = useState("")
  const [message, setMessage] = useState("")
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
      setUserId(user.uuid)
      window.location.href = `http://localhost:3000?userId=${user.uuid}`;
    }
  }

  useEffect(() => {

    const conversationId = params.conversationId
    const userId = params.userId

    setUserId(userId)
    setConversationId(conversationId)
  }, [])

  if (userId && conversationId) {
    ws.onmessage = function (event) {
      let currentMessages = [...messages, JSON.parse(event.data)]
      if (currentMessages.length > 20) {
        currentMessages.shift()
      }
      setMessages(currentMessages)
      console.log(messages.length)
    };

    const sendMessage = () => {
      let msg = {body: message, from_user_id: userId, conversation_id: conversationId}
      ws.send(JSON.stringify(msg))
      setMessage("")
    }

    return (
      <div>
        <header></header>
        <main>
          <label>Message</label>
          <input value={ message } onChange={ (event) => setMessage(event.target.value) } />
          <button onClick={ () => sendMessage() }>Send Message</button>
          <p>Messages</p>
          <div>
            {messages.map((message) => <p key={message.uuid}>{message.from_user_id}: {message.body}</p>)}
          </div>
        </main>
        <footer></footer>
      </div>
    )
  }

  return (
    <div>
      <header></header>
      <main>
        <label>Username</label>
        <input value={ username } onChange={ (event) => setUsername(event.target.value) } />
        <button onClick={ () => user() }>Create user</button>
        <button onClick={ () => convo(userId) }>Create conversation</button>
      </main>
      <footer></footer>
    </div>
  );
}

export default App;
