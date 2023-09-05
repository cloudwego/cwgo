import ReactDOM from "react-dom/client";
import {
	createBrowserRouter,
	RouterProvider,
	Navigate
} from "react-router-dom";
import App from "./App.tsx";
import "./index.css";
import IdlPage from "./pages/IdlPage/index.tsx";
import RepositoryPage from "./pages/RepositoryPage/index.tsx";
import TemplatePage from "./pages/TemplatePage/index.tsx";

const router = createBrowserRouter([
	{
		path: "/",
		element: <App />,
		children: [
			{
				path: "/idl",
				element: <IdlPage />
			},
			{
				path: "/repository",
				element: <RepositoryPage />
			},
			{
				path: "/template",
				element: <TemplatePage />
			},
			{
				path: "*",
				element: <Navigate to={"/idl"} />
			}
		]
	}
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
	<RouterProvider router={router} />
);
