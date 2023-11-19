import { deleteRepo, getRepo, updateRepo } from "./api";
import { useState, useEffect } from "react";
import {
	Button,
	ConfigProvider,
	Modal,
	Space,
	Table,
	Tag,
	Toast,
	Tooltip
} from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/table";
import { UpdateRepo } from "../../types";
import { IconInfoCircle } from "@douyinfe/semi-icons";
import ContextHolder from "./contextHolder";

export default function RepositoryPage() {
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	const [modal, contextHolder] = Modal.useModal();
	const [total, setTotal] = useState(0);
	let destroyFn = () => {};
	const pageSize = 10;

	/**
	 * 更新数据
	 * @param currentPage 当前页码
	 */
	const fetchData = async (currentPage = 1) => {
		setLoading(true);
		setPage(currentPage);
		const curDataSource = await new Promise((res) => {
			getRepo(currentPage, pageSize).then((data) => {
				res(data.repositories);
				setTotal(data.total);
				console.log(data.total);
			});
		});
		setData(curDataSource);
		setLoading(false);
		destroyFn();
	};

	useEffect(() => {
		fetchData();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

	// 列表项
	const columns = [
		{
			title: "仓库类型",
			dataIndex: "store_type",
			width: 100,
			render: (value: number) => {
				return value === 1 ? (
					<Tag color="red" size="large">
						Gitlab
					</Tag>
				) : (
					<Tag color="blue" size="large">
						Github
					</Tag>
				);
			}
		},
		{
			title: "仓库地址",
			dataIndex: "repository_url",
			width: 150,
			render: (value: string) => {
				return (
					<a href={value} target="_blank">
						{value}
					</a>
				);
			}
		},
		{
			title: "仓库最后更新时间",
			dataIndex: "last_update_time",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "仓库最后同步时间",
			dataIndex: "last_sync_time",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "状态",
			dataIndex: "status",
			width: 100,
			render: (value: number) => {
				return value === 1 ? (
					<Tag color="green" size="large">
						激活
					</Tag>
				) : (
					<Tag color="red" size="large">
						未激活
					</Tag>
				);
			}
		},
		{
			title: "快捷命令",
			render: ({ token }: { token: string }) => {
				return (
					<Space>
						<Tooltip content={token}>
							<Button
								onClick={() => {
									navigator.clipboard.writeText(token);
									Toast.success({
										content: "已复制到剪贴板"
									});
								}}
							>
								查看 token
							</Button>
						</Tooltip>
					</Space>
				);
			}
		},
		{
			title: "操作",
			render: ({ token, id, status }: UpdateRepo) => {
				return (
					<Space>
						<Button
							type="warning"
							onClick={() => {
								const toast = Toast.info({
									content: "正在更新仓库",
									duration: 0
								});
								updateRepo(id, token, status)
									.then((res) => {
										Toast.success(res);
										fetchData(currentPage);
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
							强制更新仓库
						</Button>
						{/* <Button
							type="warning"
							onClick={() => {
								const toast = Toast.info({
									content: "正在同步仓库",
									duration: 0
								});
								syncRepo(id)
									.then((res) => {
										Toast.success(res);
										fetchData(currentPage);
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
							同步仓库
						</Button> */}
						<Button
							type="danger"
							onClick={() => {
								const toast = Toast.info({
									content: "正在删除仓库",
									duration: 0
								});
								deleteRepo(id)
									.then((res) => {
										Toast.success(res);
										fetchData(currentPage);
									})
									.catch((err) => {
										Toast.error(err);
									})
									.finally(() => {
										Toast.close(toast);
									});
							}}
						>
							删除仓库
						</Button>
					</Space>
				);
			}
		},
		{
			title: "记录更新时间",
			dataIndex: "update_time",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "创建时间",
			dataIndex: "create_time",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		}
	];

	// 添加仓库弹窗配置
	const config = {
		size: "medium",
		title: "添加仓库",
		content: <ContextHolder update={fetchData} />,
		icon: <IconInfoCircle />,
		footer: null
	} as ModalReactProps;

	return (
		<ConfigProvider>
			<div>
				<div
					style={{
						paddingBottom: "1rem"
					}}
				>
					<Button
						style={{
							width: "100%"
						}}
						onClick={() => {
							const temp = modal.confirm(config);
							destroyFn = temp.destroy;
						}}
					>
						添加仓库
					</Button>
				</div>
				<div className={styles["content"]}>
					<Table
						columns={columns}
						dataSource={dataSource as Data[]}
						pagination={{
							currentPage,
							onPageChange: fetchData,
							pageSize,
							total: total
						}}
						loading={loading}
					/>
				</div>
			</div>
			{contextHolder}
		</ConfigProvider>
	);
}
