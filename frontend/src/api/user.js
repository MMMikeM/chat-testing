export const createUser = async (username) => {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json"},
    body: JSON.stringify({ name: username })
  }

  return fetch('http://localhost:3001/api/v1/users', requestOptions)
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
