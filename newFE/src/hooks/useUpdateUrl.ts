import { useNavigate } from "react-router-dom"

export type UpdateURLArgs = {
  newUserId?: number | string
  newConversationId?: number | string
}

const useUpdateUrl = ({ userId, conversationId }: { userId?: number | null; conversationId?: number | null }) => {
  const navigate = useNavigate()

  const updateUrl = ({ newUserId, newConversationId }: UpdateURLArgs) => {
    const newParams = new URLSearchParams(window.location.search)
    if (userId !== null && newUserId) {
      newParams.set("userId", newUserId.toString())
    }
    if (conversationId !== null && newConversationId) {
      newParams.set("conversationId", newConversationId.toString())
    }
    navigate(`/?${newParams.toString()}`)
  }
  return updateUrl
}

export default useUpdateUrl
