import { useState } from "react"
import { Box, Button, FormControl, Input } from "@chakra-ui/react"
import { createUser } from "@/api/user"
import useUpdateUrl from "@/hooks/useUpdateUrl"

const UserForm = () => {
  const [userName, setUserName] = useState("")
  const updateUrl = useUpdateUrl()

  const newUserId = async () => {
    const user = await createUser(userName)
    if (user?.ID !== undefined) {
      updateUrl({ newUserId: user.ID.toString() })
    }
  }

  return (
    <Box mt={4} p={4}>
      <h1>User Form</h1>
      <FormControl>
        <Input
          placeholder="Enter username..."
          value={userName}
          onChange={(event) => setUserName(event.target.value)}
          onKeyPress={(event) => (event.key === "Enter" ? newUserId() : null)}
          key="username"
          autoFocus
        />
      </FormControl>
      <Button colorScheme="purple" mt={4} onClick={() => newUserId()}>
        Create user
      </Button>
    </Box>
  )
}

export default UserForm
