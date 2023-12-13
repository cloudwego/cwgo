import { deleteToken, getToken } from "./api";
import { useState, useEffect } from "react";
import {
	Button,
	ConfigProvider,
	Form,
	Modal,
	Select,
	Space,
	Table,
	Tag,
	Toast,
	Tooltip
} from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/table";
import { IconInfoCircle } from "@douyinfe/semi-icons";
import ContextHolder from "./contextHolder";

export default function TokenPage() {
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	const [modal, contextHolder] = Modal.useModal();
	const [total, setTotal] = useState(0);
	const [searchInfo, setSearchInfo] = useState({
		repository_domain: ""
	});
	const [pageSize, setPageSize] = useState(5);
	let destroyFn = () => {};

	/**
	 * 更新数据
	 * @param currentPage 当前页码
	 */
	const fetchData = async (currentPage = 1) => {
		setLoading(true);
		setPage(currentPage);
		const curDataSource = await new Promise((res) => {
			getToken(currentPage, pageSize).then((data) => {
				res(data.tokens);
				setTotal(data.total);
			});
		});
		console.log(curDataSource);
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
			dataIndex: "repository_type",
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
			title: "仓库域名",
			dataIndex: "repository_domain",
			width: 180,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "令牌所有者",
			dataIndex: "owner",
			width: 120,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "令牌值",
			width: 180,
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
								查看令牌
							</Button>
						</Tooltip>
					</Space>
				);
			}
		},
		{
			title: "令牌状态",
			dataIndex: "status",
			width: 100,
			render: (value: number) => {
				return value === 2 ? (
					<Tag color="green" size="large">
						有效
					</Tag>
				) : (
					<Tag color="red" size="large">
						过期
					</Tag>
				);
			}
		},
		{
			title: "令牌过期时间",
			dataIndex: "expiration_time",
			width: 180,
			render: (value: string) => {
				return <div>{value}</div>;
			},
			sorter: (
				a: {
					expiration_time: string;
				},
				b: {
					expiration_time: string;
				}
			) => {
				return (
					new Date(a.expiration_time).getTime() -
					new Date(b.expiration_time).getTime()
				);
			}
		},
		{
			title: "创建时间",
			dataIndex: "create_time",
			width: 180,
			render: (value: string) => {
				return <div>{value}</div>;
			},
			sorter: (
				a: {
					create_time: string;
				},
				b: {
					create_time: string;
				}
			) => {
				return (
					new Date(a.create_time).getTime() - new Date(b.create_time).getTime()
				);
			}
		},
		{
			title: "更新时间",
			dataIndex: "update_time",
			width: 180,
			render: (value: string) => {
				return <div>{value}</div>;
			},
			sorter: (
				a: {
					update_time: string;
				},
				b: {
					update_time: string;
				}
			) => {
				return (
					new Date(a.update_time).getTime() - new Date(b.update_time).getTime()
				);
			}
		},
		{
			title: "操作",
			width: 180,
			render: ({ id }: { id: number }) => {
				return (
					<Space>
						<Button
							type="danger"
							onClick={() => {
								const toast = Toast.info({
									content: "正在删除 Token",
									duration: 0
								});
								deleteToken(id)
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
							删除 Token
						</Button>
					</Space>
				);
			}
		}
	];

	// 添加仓库弹窗配置
	const config = {
		size: "medium",
		title: "添加 Token",
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
						添加 Token
					</Button>
				</div>
				<Form
					layout="horizontal"
					onValueChange={(values) => {
						setSearchInfo(values);
					}}
					style={{
						paddingBottom: "1rem",
						display: "flex",
						justifyContent: "space-between",
						alignItems: "end"
					}}
				>
					<Space align="end">
						<Form.Input
							showClear
							field="repository_domain"
							label="仓库域名"
							style={{ width: 180 }}
						/>
						<Button
							onClick={() => {
								const toast = Toast.info({
									content: "正在搜索仓库",
									duration: 0
								});
								console.log(searchInfo);
								getToken(currentPage, pageSize, searchInfo.repository_domain)
									.then((res) => {
										console.log(res);
										setData(res.tokens);
										setTotal(res.total);
										Toast.success({
											content: "搜索成功"
										});
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
							搜索
						</Button>
					</Space>
					<Select
						style={{
							width: 120,
							textAlign: "right"
						}}
						onChange={(value) => {
							setPageSize(value as number);
						}}
						defaultValue={pageSize}
					>
						<Select.Option value={5}>5 条/页</Select.Option>
						<Select.Option value={10}>10 条/页</Select.Option>
						<Select.Option value={20}>20 条/页</Select.Option>
					</Select>
				</Form>
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
