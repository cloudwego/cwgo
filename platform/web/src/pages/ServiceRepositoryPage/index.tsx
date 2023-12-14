import { deleteRepo, getRepo, updateRepo } from "./api";
import { useState, useEffect } from "react";
import {
	Button,
	ConfigProvider,
	Form,
	// Modal,
	Popconfirm,
	Select,
	Space,
	Table,
	Tag,
	Toast,
	Tooltip
} from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
// import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/table";
// import { IconInfoCircle } from "@douyinfe/semi-icons";
// import ContextHolder from "./contextHolder";

export default function RepositoryPage({ repoType }: { repoType: string }) {
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	// const [modal, contextHolder] = Modal.useModal();
	const [total, setTotal] = useState(0);
	const [statusActive, setStatusActive] = useState(1);
	const [searchInfo, setSearchInfo] = useState({
		service_name: ""
	});
	// let destroyFn = () => {};
	const [pageSize, setPageSize] = useState(5);

	const fetchData = async (currentPage = 1) => {
		setLoading(true);
		setPage(currentPage);

		const fetchOption = {
			currentPage,
			pageSize,
			repoType
		};
		const curDataSource = await new Promise((res) => {
			getRepo(fetchOption).then((data) => {
				res(data.idls);
				setTotal(data.total);
				console.log(data.total);
			});
		});
		setData(curDataSource);

		setLoading(false);
		// destroyFn();
	};

	useEffect(() => {
		fetchData();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [repoType, pageSize]);

	// 列表项
	const columns = [
		{
			title: "仓库类型",
			dataIndex: "service_repository",
			width: 100,
			render: (value: { repository_type: number }) => {
				return value.repository_type === 1 ? (
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
			dataIndex: "service_repository.repository_domain",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "仓库名",
			dataIndex: "service_repository.repository_name",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "仓库所有者",
			dataIndex: "service_repository.repository_owner",
			width: 120,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "分支",
			dataIndex: "service_repository.repository_branch",
			width: 120,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "仓库最后更新时间",
			dataIndex: "service_repository.last_update_time",
			width: 180,
			render: (value: string) => {
				return <div>{value}</div>;
			},
			sorter: (
				a: {
					last_update_time: string;
				},
				b: {
					last_update_time: string;
				}
			) => {
				const aTime = new Date(a.last_update_time).getTime();
				const bTime = new Date(b.last_update_time).getTime();
				return aTime - bTime;
			}
		},
		{
			title: "仓库最后同步时间",
			dataIndex: "service_repository.last_sync_time",
			width: 180,
			render: (value: string) => {
				return <div>{value}</div>;
			},
			sorter: (
				a: {
					last_sync_time: string;
				},
				b: {
					last_sync_time: string;
				}
			) => {
				const aTime = new Date(a.last_sync_time).getTime();
				const bTime = new Date(b.last_sync_time).getTime();
				return aTime - bTime;
			}
		},
		{
			title: "状态",
			dataIndex: "service_repository.status",
			width: 100,
			render: (value: number) => {
				return value === 2 ? (
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
			render: ({
				commit_hash,
				service_repository,
				service_name
			}: {
				commit_hash: string;
				service_repository: {
					repository_domain: string;
					repository_owner: string;
					repository_name: string;
				};
				service_name: string;
			}) => {
				const temp = `go get ${service_repository.repository_domain}/${service_repository.repository_owner}/${service_repository.repository_name}`;
				return (
					<Space>
						<Tooltip
							content={temp}
							style={{
								maxWidth: "100vw"
							}}
						>
							<Button
								onClick={() => {
									navigator.clipboard.writeText(temp);
									Toast.success({
										content: "已复制到剪贴板"
									});
								}}
							>
								复制添加依赖
							</Button>
						</Tooltip>
						<Button
							onClick={() => {
								window.open(
									`https://${service_repository.repository_domain}/${service_repository.repository_owner}/${service_repository.repository_name}/commit/${commit_hash}`
								);
							}}
						>
							跳转 commit
						</Button>
						<Tooltip
							content={`import "${service_repository.repository_domain}/${service_repository.repository_owner}/rpc/${service_name}"`}
							style={{
								maxWidth: "100vw"
							}}
						>
							<Button
								onClick={() => {
									navigator.clipboard.writeText(
										`import "${service_repository.repository_domain}/${service_repository.repository_owner}/rpc/${service_name}"`
									);
									Toast.success({
										content: "已复制到剪贴板"
									});
								}}
							>
								复制添加依赖
							</Button>
						</Tooltip>
					</Space>
				);
			}
		},
		{
			title: "操作",
			render: ({
				service_repository: { id }
			}: {
				service_repository: {
					id: number;
					status: number;
				};
			}) => {
				return (
					<Space>
						<Popconfirm
							title="修改仓库状态"
							content={
								<div>
									<Select
										defaultValue={1}
										style={{ width: 120 }}
										onChange={(value) => {
											setStatusActive(value as number);
										}}
									>
										<Select.Option value={1}>未激活</Select.Option>
										<Select.Option value={2}>激活</Select.Option>
									</Select>
								</div>
							}
							onConfirm={() => {
								const toast = Toast.info({
									content: "正在修改仓库状态",
									duration: 0
								});
								updateRepo(id, "", statusActive)
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
										setStatusActive(1);
										Toast.close(toast);
									});
							}}
							onCancel={() => {
								setStatusActive(1);
							}}
						>
							<Button type="warning">修改仓库状态</Button>
						</Popconfirm>
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
			dataIndex: "service_repository.update_time",
			width: 150,
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
				const aTime = new Date(a.update_time).getTime();
				const bTime = new Date(b.update_time).getTime();
				return aTime - bTime;
			}
		},
		{
			title: "创建时间",
			dataIndex: "service_repository.create_time",
			width: 150,
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
				const aTime = new Date(a.create_time).getTime();
				const bTime = new Date(b.create_time).getTime();
				return aTime - bTime;
			}
		}
	];

	// 添加仓库弹窗配置
	// const config = {
	// 	size: "medium",
	// 	title: "添加服务代码生成仓库",
	// 	content: <ContextHolder update={fetchData} />,
	// 	icon: <IconInfoCircle />,
	// 	footer: null
	// } as ModalReactProps;

	return (
		<ConfigProvider>
			<div>
				{/* <div
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
						添加服务代码生成仓库
					</Button>
				</div> */}
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
							field="service_name"
							label="服务名"
							style={{ width: 180 }}
						/>
						<Button
							htmlType="submit"
							onClick={() => {
								const toast = Toast.info({
									content: "正在搜索仓库",
									duration: 0
								});
								console.log(searchInfo);
								getRepo({
									currentPage,
									pageSize,
									repoType,
									...searchInfo
								})
									.then((res) => {
										console.log(res);
										setData(res.idls);
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
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
						columns={columns}
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
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
			{/* {contextHolder} */}
		</ConfigProvider>
	);
}
