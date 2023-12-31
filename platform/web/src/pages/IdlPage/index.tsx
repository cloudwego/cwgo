import { deleteIdl, getIdl, updateIdl } from "./api";
import { useState, useEffect } from "react";
import {
	Button,
	ConfigProvider,
	Form,
	Popconfirm,
	Select,
	Space,
	// Modal,
	Table,
	Tag,
	Toast
	// Toast,
} from "@douyinfe/semi-ui";
import { Dropdown } from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
// import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/table";
import { updateRepo } from "../ServiceRepositoryPage/api";
// import { IconInfoCircle } from "@douyinfe/semi-icons";
// import ContextHolder from "./contextHolder";

interface Idls {
	commit_hash: string;
	idl_path: string;
}

export default function RepositoryPage() {
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	// const [modal, contextHolder] = Modal.useModal();
	const [total, setTotal] = useState(0);
	const [searchInfo, setSearchInfo] = useState({
		service_name: ""
	});
	const [statusActive, setStatusActive] = useState(1);
	// const pageSize = 10;
	// 从本地存储中获取 pageSize
	const [pageSize, setPageSize] = useState(
		Number(localStorage.getItem("pageSize")) || 10
	);
	// let destroyFn = () => {};

	function InnerIdls({
		data,
		repo
	}: {
		data: Idls[];
		repo: {
			repository_domain: string;
			repository_owner: string;
			repository_name: string;
			repository_branch: string;
		};
	}) {
		return (
			<div>
				{data.map((item) => {
					return (
						<div>
							<Dropdown.Item
								type="primary"
								style={{
									maxWidth: "100vw"
								}}
								onClick={() => {
									let url: string = "";
									switch (repo.repository_domain) {
									case "github.com":
										url = `https://${repo.repository_domain}/${repo.repository_owner}/${repo.repository_name}/blob/${repo.repository_branch}/${item.idl_path}`;
										break;
									case "gitlab.com":
										url = `https://${repo.repository_domain}/${repo.repository_owner}/${repo.repository_name}/-/blob/${repo.repository_branch}/${item.idl_path}`;
										break;
									default:
										url = `https://${repo.repository_domain}/${repo.repository_owner}/${repo.repository_name}/-/blob/${repo.repository_branch}/${item.idl_path}`;
									}

									// 跳转到对应的 idl
									window.open(url);
								}}
							>
								{item.idl_path}
							</Dropdown.Item>
						</div>
					);
				})}
			</div>
		);
	}

	/**
	 * 更新数据
	 * @param currentPage 当前页码
	 */
	const fetchData = async (currentPage = 1) => {
		setLoading(true);
		setPage(currentPage);
		const curDataSource = await new Promise((res) => {
			getIdl(currentPage, pageSize).then((data) => {
				res(data);
				setTotal(data.total);
			});
		});
		console.log(curDataSource);
		setData(curDataSource);
		setLoading(false);
		// destroyFn();
	};

	useEffect(() => {
		fetchData();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [pageSize]);

	// 列表项
	const columns = [
		{
			title: "服务名",
			dataIndex: "service_name",
			width: 100,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "主 idl 路径",
			dataIndex: "main_idl_path",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "idl 最后同步时间",
			dataIndex: "last_sync_time",
			width: 180,
			render: (value: string) => {
				return <div>{value}</div>;
			},
			sorter: (
				a: { last_sync_time: string },
				b: { last_sync_time: string }
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
						同步
					</Tag>
				) : (
					<Tag color="red" size="large">
						未激活
					</Tag>
				);
			}
		},
		{
			title: "import idls",
			width: 150,
			render: ({
				import_idls,
				idl_repository
			}: {
				import_idls: [];
				idl_repository: {
					repository_domain: string;
					repository_owner: string;
					repository_name: string;
					repository_branch: string;
				};
			}) => {
				return import_idls.length ? (
					<Dropdown
						trigger={"hover"}
						showTick
						position={"bottomLeft"}
						render={
							<Dropdown.Menu
								// eslint-disable-next-line @typescript-eslint/ban-ts-comment
								// @ts-ignore
								children={
									<InnerIdls data={import_idls} repo={idl_repository} />
								}
							/>
						}
					>
						<Button>查看 import idls</Button>
					</Dropdown>
				) : (
					<div>无 import idls</div>
				);
			}
		},
		{
			title: "操作",
			render: ({
				id,
				service_repository,
				service_name
			}: {
				id: number;
				service_repository: {
					id: number;
				};
				service_name: string;
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
										<Select.Option value={2}>同步</Select.Option>
									</Select>
								</div>
							}
							onConfirm={() => {
								const toast = Toast.info({
									content: "正在修改仓库状态",
									duration: 0
								});
								updateRepo(service_repository.id, "", statusActive)
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
							<Button type="warning">更改同步状态</Button>
						</Popconfirm>
						<Button
							type="warning"
							onClick={() => {
								const toast = Toast.info({
									content: "正在同步 idl",
									duration: 0
								});
								updateIdl(id, service_name)
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
							强制同步 idl
						</Button>
						<Popconfirm
							title="确定删除"
							onConfirm={() => {
								const toast = Toast.info({
									content: "正在删除 idl",
									duration: 0
								});
								deleteIdl(id)
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
							<Button type="danger">删除 idl</Button>
						</Popconfirm>
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
			},
			sorter: (a: { update_time: string }, b: { update_time: string }) => {
				const aTime = new Date(a.update_time).getTime();
				const bTime = new Date(b.update_time).getTime();
				return aTime - bTime;
			}
		},
		{
			title: "创建时间",
			dataIndex: "create_time",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			},
			sorter: (a: { create_time: string }, b: { create_time: string }) => {
				const aTime = new Date(a.create_time).getTime();
				const bTime = new Date(b.create_time).getTime();
				return aTime - bTime;
			}
		}
	];

	// // 添加仓库弹窗配置
	// const config = {
	// 	size: "medium",
	// 	title: "添加 idl",
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
						添加 idl
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
								getIdl(currentPage, pageSize, searchInfo.service_name)
									.then((res) => {
										console.log(res);
										setData(res);
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
							localStorage.setItem("pageSize", String(value));
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
