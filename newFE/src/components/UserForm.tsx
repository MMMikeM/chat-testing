import { useState } from "react"
import { Box, Button, FormControl, Input } from "@chakra-ui/react"
import { createUser } from "../api/user"
import useUpdateUrl from "../hooks/useUpdateUrl"

const UserForm = ({ userId }: { userId: number | null }) => {
  const [userName, setUserName] = useState("")
  const updateUrl = useUpdateUrl({ userId })

  const newUserId = async () => {
    const user = await createUser(userName)
    if (user?.id) {
      updateUrl({ newUserId: user.id })
    }
  }

  return (
    <Box mt={4} p={4}>
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
