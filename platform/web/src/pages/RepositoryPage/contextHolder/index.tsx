import { Button, Input, Select, Space, Toast } from "@douyinfe/semi-ui";
import { useState } from "react";
import { createRepo } from "../api";

export default function ContextHolder({ update }: { update: () => void }) {
	// 表单状态
	const [url, setUrl] = useState("");
	const [token, setToken] = useState("");
	const [storeType, setStoreType] = useState(1);
	const [repoType, setRepoType] = useState(1);

	return (
		<Space vertical>
			<Space
				style={{
					display: "flex",
					justifyContent: "space-between",
					width: "100%"
				}}
			>
				<div>仓库 URL</div>
				<Input
					style={{
						width: "30rem"
					}}
					showClear
					onChange={(value) => {
						setUrl(value);
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
					token
				</div>
				<Input
					style={{
						width: "30rem"
					}}
					showClear
					onChange={(value) => {
						setToken(value);
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
					存储类型
				</div>
				<div
					style={{
						width: "30rem"
					}}
				>
					<Select
						defaultValue={1}
						style={{ width: 120 }}
						onChange={(value) => {
							setStoreType(value as number);
						}}
					>
						<Select.Option value={1}>IDL 文件</Select.Option>
						<Select.Option value={2}>服务代码</Select.Option>
					</Select>
				</div>
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
					仓库类型
				</div>
				<div
					style={{
						width: "30rem"
					}}
				>
					<Select
						defaultValue={1}
						style={{ width: 120 }}
						onChange={(value) => {
							setRepoType(value as number);
						}}
					>
						<Select.Option value={1}>Gitlab</Select.Option>
						<Select.Option value={2}>Github</Select.Option>
					</Select>
				</div>
			</Space>
			<Button
				style={{
					width: "100%",
					margin: "2rem 0"
				}}
				type="primary"
				onClick={() => {
					if (!url || !token) {
						Toast.error({
							content: "请填写完整信息"
						});
						return;
					}
					const toast = Toast.info({
						content: "正在更新仓库",
						duration: 0
					});
					createRepo(repoType, url, token, storeType)
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
