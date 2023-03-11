import {
  createBrowserRouter,
} from "react-router-dom";
import Home from "./home";

export const router = createBrowserRouter([
    {
      path: "/",
      element: <Home />,
      children: [
        {
          path: "team",
          element: <>Team</>,
        },
      ],
    },
  ]);