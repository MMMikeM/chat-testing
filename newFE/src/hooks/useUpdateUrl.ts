import { useNavigate, useParams } from "react-router-dom"

export type UpdateURLArgs = {
  newUserId?: string
  newConversationId?: string
}

const useUpdateUrl = () => {
  const navigate = useNavigate()
  const params = useParams()

  const { userId, conversationId } = params

  const updateUrl = ({ newUserId, newConversationId }: UpdateURLArgs) => {
    const url = `/${newUserId ?? userId}/${newConversationId ?? conversationId ?? ""}`

    navigate(url)
  }
  return updateUrl
}

export default useUpdateUrl
