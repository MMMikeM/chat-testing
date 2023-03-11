import React from "react"
import ReactDOM from "react-dom/client"
import { ChakraProvider, Container } from "@chakra-ui/react"
import { RouterProvider } from "react-router-dom"
import { router } from "./routes/router"

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <React.StrictMode>
    <ChakraProvider>
      <Container maxW={"6xl"}>
        <RouterProvider router={router} />
      </Container>
    </ChakraProvider>
  </React.StrictMode>
)
