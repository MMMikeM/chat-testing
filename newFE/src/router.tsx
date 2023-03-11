import { createBrowserRouter } from "react-router-dom"
import UserForm from "./pages/UserForm"
import ConversationForm from "./pages/ConversationForm"
import Main from "./pages/Main"

export const router = createBrowserRouter([
  {
    path: "/",
    element: <UserForm />,
  },
  {
    path: "/:userId",
    element: <ConversationForm />,
  },
  {
    path: "/:userId/:conversationId",
    element: <Main />,
  },
])
