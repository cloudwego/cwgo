import { Button } from "@douyinfe/semi-ui";
import { IconMoon, IconSun } from "@douyinfe/semi-icons";
import { useState } from "react";

export function SwitchButton() {
	const [themeMode, setThemeMode] = useState(
		localStorage.getItem("theme-mode") === "dark"
	);

	const switchMode = () => {
		const body = document.body;
		if (body.hasAttribute("theme-mode")) {
			body.removeAttribute("theme-mode");
			// 将主题模式存储到 localStorage 中
			localStorage.removeItem("theme-mode");
			setThemeMode(false);
		} else {
			body.setAttribute("theme-mode", "dark");
			// 将主题模式存储到 localStorage 中
			localStorage.setItem("theme-mode", "dark");
			setThemeMode(true);
		}
	};

	return (
		<Button
			icon={themeMode ? <IconMoon /> : <IconSun />}
			onClick={switchMode}
		></Button>
	);
}
