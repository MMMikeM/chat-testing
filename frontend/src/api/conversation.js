export const createConversation = async () => {
  const requestOptions = {method: "POST", headers: { "Content-Type": "application/json"}}
  return fetch('http://localhost:3001/api/v1/conversations', requestOptions)
    .then((response) => {
      return response.json()
    })
    .then((data) => {
      return data
    })
    .catch((error) => {
      console.log(error)
    })
}
