import { Layout, Nav, Breadcrumb, Skeleton } from "@douyinfe/semi-ui";
import { IconCode, IconGitlabLogo, IconFolderOpen } from "@douyinfe/semi-icons";
import { Outlet, useNavigate } from "react-router-dom";
import styles from "./App.module.scss";
import { useEffect, useState } from "react";
import Header from "./components/Header";
import Footer from "./components/Footer";

export default function App() {
	const { Sider, Content } = Layout;
	const [selectedKey, setSelectedKey] = useState("idl");
	const navList = [
		{
			itemKey: "idl",
			text: "IDL 信息查询",
			icon: <IconCode size="large" />
		},
		{
			itemKey: "repository",
			text: "仓库信息管理",
			icon: <IconGitlabLogo size="large" />
		},
		{
			itemKey: "template",
			text: "模版管理",
			icon: <IconFolderOpen size="large" />
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
