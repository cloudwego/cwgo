import { deleteIdl, getIdl, updateIdl } from "./api";
import { useState, useEffect } from "react";
import {
	Button,
	ConfigProvider,
	Modal,
	Space,
	Table,
	Toast,
	// Toast,
	Tooltip
} from "@douyinfe/semi-ui";
import { Dropdown } from "@douyinfe/semi-ui";
import styles from "./index.module.scss";
import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/table";
import { IconInfoCircle } from "@douyinfe/semi-icons";
import ContextHolder from "./contextHolder";

interface Idls {
	commit_hash: string;
	idl_path: string;
}

export default function RepositoryPage() {
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	const [modal, contextHolder] = Modal.useModal();
	const [total, setTotal] = useState(0);
	const pageSize = 10;
	let destroyFn = () => {};

	function InnerIdls({ data }: { data: Idls[] }) {
		return (
			<div>
				{data.map((item, index) => {
					return (
						<div>
							<Dropdown.Title>idl {index + 1}</Dropdown.Title>
							<Dropdown.Item
								type="primary"
								style={{
									marginRight: "0.5rem",
									overflow: "hidden",
									textOverflow: "ellipsis",
									whiteSpace: "nowrap"
								}}
								onClick={() => {
									navigator.clipboard.writeText(item.commit_hash);
									Toast.success({
										content: "已复制到剪贴板"
									});
								}}
							>
								{item.commit_hash}
							</Dropdown.Item>
							<Dropdown.Item
								type="primary"
								onClick={() => {
									navigator.clipboard.writeText(item.idl_path);
									Toast.success({
										content: "已复制到剪贴板"
									});
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
		destroyFn();
	};

	useEffect(() => {
		fetchData();
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, []);

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
			title: "快捷命令",
			render: ({ commit_hash }: { commit_hash: string }) => {
				return (
					<Space>
						<Tooltip content={commit_hash}>
							<Button>查看 commit_hash</Button>
						</Tooltip>
					</Space>
				);
			}
		},
		{
			title: "idl 最后同步时间",
			dataIndex: "last_sync_time",
			width: 150,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "import idls",
			dataIndex: "import_idls",
			width: 150,
			render: (value: []) => {
				console.log(value);
				return value.length ? (
					<Dropdown
						trigger={"hover"}
						showTick
						position={"bottomLeft"}
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
						render={<Dropdown.Menu children={<InnerIdls data={value} />} />}
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
			render: ({ id, service_name }: { id: number; service_name: string }) => {
				return (
					<Space>
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
							同步 idl
						</Button>
						<Button
							type="danger"
							onClick={() => {
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
							删除 idl
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
		title: "添加 idl",
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
						添加 idl
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
