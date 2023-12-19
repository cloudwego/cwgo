import { Layout, Nav } from "@douyinfe/semi-ui";
import { SwitchButton } from "../SwitchButton";
import styles from "./index.module.scss";
import logo from "../../assets/favicon.ico";

export default function Header() {
	const { Header } = Layout;
	const Img = () => <img src={logo} alt="logo" className={styles["logo"]} />;
	return (
		<Header className={styles["header"]}>
			<Nav mode="horizontal" className={styles["nav"]}>
				<Nav.Header>
					{/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
					{/* @ts-ignore */}
					<Img />
				</Nav.Header>
				<span className={styles["logo-title"]}>一站式 RPC 调用平台</span>
				<Nav.Footer>
					{/* eslint-disable-next-line @typescript-eslint/ban-ts-comment */}
					{/* @ts-ignore */}
					<SwitchButton />
				</Nav.Footer>
			</Nav>
		</Header>
	);
}
