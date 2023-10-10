import { Layout, Nav } from "@douyinfe/semi-ui";
import { SwitchButton } from "../SwitchButton";
import styles from "./index.module.scss";
import logo from "../../../public/favicon.ico";

export default function Header() {
	const { Header } = Layout;
	return (
		<Header className={styles["header"]}>
			<Nav mode="horizontal" className={styles["nav"]}>
				<Nav.Header>
					<img src={logo} alt="logo" className={styles["logo"]} />
				</Nav.Header>
				<span className={styles["logo-title"]}>一站式 RPC 调用平台</span>
				<Nav.Footer>
					<SwitchButton />
				</Nav.Footer>
			</Nav>
		</Header>
	);
}
