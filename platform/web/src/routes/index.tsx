import { createBrowserRouter, Navigate } from "react-router-dom";
import App from "../App";
import IdlPage from "../pages/IdlPage";
import RepositoryPage from "../pages/RepositoryPage";
import TemplatePage from "../pages/TemplatePage";

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

export default router;
