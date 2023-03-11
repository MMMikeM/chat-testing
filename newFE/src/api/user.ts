type User = {
  id: number
  name: string
  created_at: string
  updated_at: string
}

export const createUser = async (username: string) => {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name: username }),
  }

  return fetch("http://localhost:3001/api/v1/users", requestOptions)
    .then((response) => response.json())
    .then((data) => data as User)
    .catch((error) => console.log(error))
}
