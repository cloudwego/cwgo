import { deleteRepo, getRepo, updateRepo } from "./api";
import { useState, useEffect } from "react";
import {
	Button,
	ConfigProvider,
	Form,
	Input,
	Modal,
	Popconfirm,
	Select,
	Space,
	Table,
	Tag,
	Toast
} from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/table";
import { UpdateRepo } from "../../types";
import { IconInfoCircle } from "@douyinfe/semi-icons";
import ContextHolder from "./contextHolder";
import AddIdlContextHolder from "./contextHolder/addIdl";

export default function RepositoryPage({ repoType }: { repoType: string }) {
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	const [modal, contextHolder] = Modal.useModal();
	const [total, setTotal] = useState(0);
	const [serviceRepoName, setServiceRepoName] = useState("");
	const [statusActive, setStatusActive] = useState(1);
	const [searchInfo, setSearchInfo] = useState({
		repository_domain: "",
		repository_name: ""
	});
	let destroyFn = () => {};
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
	}, [repoType, pageSize]);

	// 添加仓库弹窗配置
	const getIdlConfig = (id: number) => {
		const addIdlConfig = {
			size: "medium",
			title: "添加 idl",
			content: <AddIdlContextHolder update={fetchData} id={id} />,
			icon: <IconInfoCircle />,
			footer: null
		} as ModalReactProps;
		return addIdlConfig;
	};

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
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "仓库名",
			dataIndex: "repository_name",
			width: 130,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "仓库所有者",
			dataIndex: "repository_owner",
			width: 120,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "分支",
			dataIndex: "repository_branch",
			width: 120,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "仓库最后更新时间",
			dataIndex: "last_update_time",
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
				return (
					new Date(b.last_update_time).getTime() -
					new Date(a.last_update_time).getTime()
				);
			}
		},
		{
			title: "仓库最后同步时间",
			dataIndex: "last_sync_time",
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
				return (
					new Date(b.last_sync_time).getTime() -
					new Date(a.last_sync_time).getTime()
				);
			}
		},
		{
			title: "状态",
			dataIndex: "status",
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
			title: "操作",
			render: ({ id }: UpdateRepo) => {
				return (
					<Space>
						<Button
							type="warning"
							onClick={() => {
								const addIdlConfig = getIdlConfig(id);
								const temp = modal.confirm(addIdlConfig);
								destroyFn = temp.destroy;
							}}
						>
							添加 IDL 信息
						</Button>
						<Popconfirm
							title="修改仓库分支"
							content={
								<div>
									<Input
										onChange={(value) => {
											setServiceRepoName(value);
										}}
									></Input>
								</div>
							}
							onConfirm={() => {
								if (!serviceRepoName) {
									Toast.error({
										content: "仓库分支不能为空"
									});
									return;
								}
								const toast = Toast.info({
									content: "正在修改仓库分支",
									duration: 0
								});
								serviceRepoName &&
									updateRepo(id, serviceRepoName, 0)
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
											setServiceRepoName("");
										});
							}}
							onCancel={() => {
								setServiceRepoName("");
							}}
						>
							<Button type="warning">修改仓库分支</Button>
						</Popconfirm>
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
			dataIndex: "update_time",
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
				return (
					new Date(a.update_time).getTime() - new Date(b.update_time).getTime()
				);
			}
		},
		{
			title: "创建时间",
			dataIndex: "create_time",
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
				return (
					new Date(a.create_time).getTime() - new Date(b.create_time).getTime()
				);
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
						<Form.Input
							showClear
							field="repository_name"
							label="仓库名"
							style={{ width: 180 }}
						/>
						<Button
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
										setData(res.repositories);
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
							marginLeft: "1rem"
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
