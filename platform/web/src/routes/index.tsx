import { createBrowserRouter, Navigate } from "react-router-dom";
import App from "../App";
import IdlPage from "../pages/IdlPage";
import ServiceRepositoryPage from "../pages/ServiceRepositoryPage";
import IdlRepositoryPage from "../pages/IdlRepositoryPage";
import TokenInfoPage from "../pages/TokenInfoPage";
// import TemplatePage from "../pages/TemplatePage";

const router = createBrowserRouter([
	{
		path: "/",
		element: <App />,
		children: [
			{
				path: "/idlrepo",
				element: <IdlRepositoryPage repoType="idl" />
			},
			{
				path: "/idlinfo",
				element: <IdlPage />
			},
			{
				path: "/servicerepo",
				element: <ServiceRepositoryPage repoType="service" />
			},
			{
				path: "/tokeninfo",
				element: <TokenInfoPage />
			},
			// {
			// 	path: "/template",
			// 	element: <TemplatePage />
			// },
			{
				path: "*",
				element: <Navigate to={"/idlrepo"} />
			}
		]
	}
]);

export default router;
