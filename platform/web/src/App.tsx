import { Layout, Nav, Breadcrumb, Skeleton } from "@douyinfe/semi-ui";
import { IconCode, IconGitlabLogo } from "@douyinfe/semi-icons";
import { Outlet, useNavigate } from "react-router-dom";
import styles from "./App.module.scss";
import { useEffect, useState } from "react";
import Header from "./components/Header";
import Footer from "./components/Footer";

export default function App() {
	const { Sider, Content } = Layout;
	// 获取当前的路由参数
	const navigatePath = window.location.pathname.split("/")[1];
	const [selectedKey, setSelectedKey] = useState(navigatePath || "tokeninfo");
	const navList = [
		{
			itemKey: "tokeninfo",
			text: "令牌管理",
			icon: <IconCode size="large" />
		},
		{
			itemKey: "idlrepo",
			text: "IDL 仓库信息管理",
			icon: <IconGitlabLogo size="large" />
		},
		{
			itemKey: "idlinfo",
			text: "IDL 信息同步",
			icon: <IconCode size="large" />
		},
		{
			itemKey: "servicerepo",
			text: "服务代码产物仓库查询",
			icon: <IconGitlabLogo size="large" />
		}
	];

	// 主题模式
	const themeMode = localStorage.getItem("theme-mode");
	if (themeMode) {
		document.body.setAttribute("theme-mode", themeMode);
	}
	window.addEventListener("storage", (event) => {
		if (event.key === "theme-mode") {
			const body = document.body;
			if (event.newValue) {
				body.setAttribute("theme-mode", event.newValue);
			} else {
				body.removeAttribute("theme-mode");
			}
		}
	});

	// 路由跳转
	const navigate = useNavigate();
	useEffect(() => {
		navigate(selectedKey);
	}, [navigate, selectedKey]);

	return (
		<Layout className={styles["layout"]}>
			<Header />
			<Layout className={styles["inner-layout"]}>
				<Sider className={styles["sider"]}>
					<Nav
						className={styles["nav"]}
						selectedKeys={[selectedKey]}
						items={navList}
						onSelect={(item) => {
							setSelectedKey(String(item.itemKey));
						}}
						footer={{
							collapseButton: true
						}}
					/>
				</Sider>
				<Content
					style={{
						padding: "24px",
						backgroundColor: "var(--semi-color-bg-0)"
					}}
					className={styles["content-wrapper"]}
				>
					<Breadcrumb
						style={{
							marginBottom: "24px"
						}}
						routes={[
							navList.find((item) => item.itemKey === selectedKey)?.text || ""
						]}
					/>
					<Skeleton
						placeholder={<Skeleton.Paragraph rows={2} />}
						loading={false}
					>
						<Outlet />
					</Skeleton>
				</Content>
			</Layout>
			<Footer />
		</Layout>
	);
}
