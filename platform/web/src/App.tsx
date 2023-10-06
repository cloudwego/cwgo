import { Layout, Nav, Breadcrumb, Skeleton } from "@douyinfe/semi-ui";
import {
	IconCode,
	IconBytedanceLogo,
	IconHome,
	IconHistogram,
	IconLive
} from "@douyinfe/semi-icons";
import { Outlet, useNavigate } from "react-router-dom";
import styles from "./App.module.scss";
import { SwitchButton } from "./components/SwitchButton";
import { useEffect, useState } from "react";

export default function App() {
	const { Header, Footer, Sider, Content } = Layout;
	const [selectedKey, setSelectedKey] = useState("idl");
	const navigate = useNavigate();
	const themeMode = localStorage.getItem("theme-mode");
	const navList = [
		{
			itemKey: "idl",
			text: "IDL 信息查询",
			icon: <IconHome size="large" />
		},
		{
			itemKey: "repository",
			text: "仓库信息管理",
			icon: <IconHistogram size="large" />
		},
		{
			itemKey: "template",
			text: "模版管理",
			icon: <IconLive size="large" />
		}
	];

	if (themeMode) {
		document.body.setAttribute("theme-mode", themeMode);
	}
	// 开启监听 storage 事件
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

	useEffect(() => {
		navigate(selectedKey);
	}, [navigate, selectedKey]);

	return (
		<Layout className={styles["layout"]}>
			<Header className={styles["header"]}>
				<Nav mode="horizontal" className={styles["nav"]}>
					<Nav.Header>
						<IconCode className={styles["logo"]} />
					</Nav.Header>
					<span className={styles["logo-title"]}>一站式 RPC 调用平台</span>
					<Nav.Footer>
						<SwitchButton />
					</Nav.Footer>
				</Nav>
			</Header>
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
			<Footer className={styles["footer"]}>
				<span className={styles["copyright"]}>
					<IconBytedanceLogo size="large" style={{ marginRight: "8px" }} />
					<span>Copyright © 2023. All Rights Reserved. </span>
				</span>
				<span>
					<span>反馈建议</span>
				</span>
			</Footer>
		</Layout>
	);
}
