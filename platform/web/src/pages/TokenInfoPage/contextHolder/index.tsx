import { Button, Input, Select, Space, Toast } from "@douyinfe/semi-ui";
import { useState } from "react";
import { createIdl } from "../api";

export default function ContextHolder({ update }: { update: () => void }) {
	// 表单状态
	const [repoType, setRepoType] = useState(1);
	const [repoDomain, setRepoDomain] = useState("");
	const [tokenValue, setTokenValue] = useState("");

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
					仓库域名（可选）
				</div>
				<Input
					style={{
						width: "30rem"
					}}
					showClear
					placeholder={"Github 默认为 github.com，Gitlab 默认为 gitlab.com"}
					onChange={(value) => {
						setRepoDomain(value);
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
					Token
				</div>
				<Input
					style={{
						width: "30rem"
					}}
					showClear
					onChange={(value) => {
						setTokenValue(value);
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
					if (!repoType) {
						Toast.error({
							content: "仓库类型不能为空"
						});
						return;
					}
					if (!tokenValue) {
						Toast.error({
							content: "Token 不能为空"
						});
						return;
					}
					const toast = Toast.info({
						content: "正在添加 Token",
						duration: 0
					});
					createIdl(repoType, repoDomain, tokenValue)
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
