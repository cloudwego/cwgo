import { Layout } from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
import logo from "../../assets/favicon.ico";

export default function Footer() {
	const { Footer } = Layout;
	return (
		<Footer className={styles["footer"]}>
			<span className={styles["copyright"]}>
				<img src={logo} className={styles["logo"]} />
				<span>Copyright © 2023. All Rights Reserved. </span>
			</span>
			<span>
				<span>反馈建议</span>
			</span>
		</Footer>
	);
}
