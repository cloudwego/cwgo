import { Layout } from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
import logo from "../../assets/favicon.ico";

export default function Footer() {
	const { Footer } = Layout;
	return (
		<Footer className={styles["footer"]}>
			<span className={styles["copyright"]}>
				<img src={logo} className={styles["logo"]} />
				<span>Copyright © 2023 CloudWeGo. All Rights Reserved. </span>
			</span>
			<span>
				<a href="https://github.com/cloudwego/cwgo/issues" target="__blank">
					反馈建议
				</a>
			</span>
		</Footer>
	);
}
