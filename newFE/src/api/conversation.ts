type Conversation = {
  ID: number
  name: string
  created_at: string
  updated_at: string
}

export const createConversation = async () => {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  }

  return fetch("http://localhost:3001/api/v1/conversations", requestOptions)
    .then((response) => response.json())
    .then((data) => data as Conversation)
    .catch((error) => console.log(error))
}
