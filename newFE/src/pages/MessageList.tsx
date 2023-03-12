import { memo, useMemo, useRef } from "react"
import { Box, Text } from "@chakra-ui/react"
import { Message } from "./Main"
import useWebsocket from "@/hooks/useWebsocket"

const handleMessageHistory = (messageHistory: Message[], latestMessage?: Message) => {
  if (!latestMessage) return []
  if (latestMessage !== messageHistory[messageHistory.length - 1]) {
    if (messageHistory.length > 10) messageHistory.shift()
    return messageHistory.concat(latestMessage)
  }
  return messageHistory
}

const MessageList = () => {
  const { latestMessage } = useWebsocket()
  const messageHistory = useRef<Message[]>([])
  messageHistory.current = useMemo(() => handleMessageHistory(messageHistory.current, latestMessage), [latestMessage])

  if (!messageHistory.current.length)
    return (
      <Box mt={4}>
        <Text colorScheme="gray">No messages exist yet...</Text>
      </Box>
    )
  return (
    <Box mt={4}>
      {messageHistory.current.map((message) => (
        <Message key={message.uuid} message={message} />
      ))}
    </Box>
  )
}
type MessageProps = {
  message: {
    uuid: string
    from_user_id: string
    body: string
  }
}
const Message = memo(({ message }: MessageProps) => (
  <p key={message.uuid}>
    {message.from_user_id}: {message.body}
  </p>
))

export default memo(MessageList, (p, n) => {
  console.log("p:", p)
  console.log("n:", n)
  return true
})
