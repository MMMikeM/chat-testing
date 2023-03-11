import { createBrowserRouter } from "react-router-dom"
import UserForm from "."
import ConversationForm from "./userid"
import Main from "./userid/conversationid"

export const router = createBrowserRouter([
  {
    path: "/",
    element: <UserForm />,
    children: [
      {
        path: "/:userId",
        element: <ConversationForm />,
        children: [
          {
            path: "/:userId/:conversationId",
            element: <Main />,
          },
        ],
      },
    ],
  },
])
