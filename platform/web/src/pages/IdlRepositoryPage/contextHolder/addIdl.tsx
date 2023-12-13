import { Button, Input, Space, Toast } from "@douyinfe/semi-ui";
import { useState } from "react";
import { createIdl } from "../api";

export default function ContextHolder({
	update,
	id: idProp
}: {
	update: () => void;
	id: number;
}) {
	// 表单状态
	const [idlPath, setIdlPath] = useState("");
	const [serviceName, setServiceName] = useState("");
	const [serviceRepoName, setServiceRepoName] = useState("");

	return (
		<Space vertical>
			<Space
				style={{
					display: "flex",
					justifyContent: "space-between",
					width: "100%"
				}}
			>
				<div
					style={{
						width: "5rem"
					}}
				>
					主 idl 路径
				</div>
				<Input
					style={{
						width: "30rem"
					}}
					showClear
					onChange={(value) => {
						setIdlPath(value);
					}}
				></Input>
			</Space>
			<Space
				style={{
					display: "flex",
					justifyContent: "space-between",
					width: "100%"
				}}
			>
				<div
					style={{
						width: "5rem"
					}}
				>
					服务名
				</div>
				<Input
					style={{
						width: "30rem"
					}}
					showClear
					onChange={(value) => {
						setServiceName(value);
					}}
				></Input>
			</Space>
			<Space
				style={{
					display: "flex",
					justifyContent: "space-between",
					width: "100%"
				}}
			>
				<div
					style={{
						width: "5rem"
					}}
				>
					服务仓库名
				</div>
				<Input
					style={{
						width: "30rem"
					}}
					showClear
					onChange={(value) => {
						setServiceRepoName(value);
					}}
				></Input>
			</Space>
			<Button
				style={{
					width: "100%",
					margin: "2rem 0"
				}}
				type="primary"
				onClick={() => {
					if (!idlPath) {
						Toast.error({
							content: "主 idl 路径不能为空"
						});
						return;
					}
					if (!serviceName) {
						Toast.error({
							content: "服务名不能为空"
						});
						return;
					}
					if (!serviceRepoName) {
						Toast.error({
							content: "服务仓库名不能为空"
						});
						return;
					}
					const toast = Toast.info({
						content: "正在更新仓库",
						duration: 0
					});
					createIdl(idProp, idlPath, serviceName, serviceRepoName)
						.then((res) => {
							Toast.success({
								content: res
							});
							update();
						})
						.catch((err) => {
							Toast.error({
								content: err.response.data.msg
							});
						})
						.finally(() => {
							Toast.close(toast);
						});
				}}
			>
				添加
			</Button>
		</Space>
	);
}
