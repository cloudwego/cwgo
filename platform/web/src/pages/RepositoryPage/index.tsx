import { getIDLsRes } from "./api";
import { useState, useEffect } from "react";
import {
	Button,
	ConfigProvider,
	Input,
	Modal,
	Space,
	Table,
	Tag,
	Toast,
	Tooltip
} from "@douyinfe/semi-ui";
import { IconInfoCircle, IconSearch } from "@douyinfe/semi-icons";
import styles from "./index.module.scss";
import AllInfoBox from "../../components/AllInfoBox";
import en_GB from "@douyinfe/semi-ui/lib/es/locale/source/en_GB";
import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/table";

const pageSize = 10;

const columns = [
	{
		title: "仓库地址",
		dataIndex: "repository_url",
		width: 150,
		render: (value: string) => {
			return <a href="#">{value}</a>;
		}
	},
	{
		title: "仓库最近一次更新时间",
		dataIndex: "last_update_time",
		width: 200,
		render: (value: string) => {
			return <div>{value}</div>;
		}
	},
	{
		title: "IDL 变动检测时间",
		dataIndex: "last_sync_time",
		width: 200,
		render: (value: string) => {
			return <div>{value}</div>;
		}
	},
	{
		title: "状态",
		dataIndex: "status",
		width: 100,
		render: (value: string) => {
			return value ? (
				<Tag color="green" size="large">
					OK
				</Tag>
			) : (
				<Tag color="red" size="large">
					ERROR
				</Tag>
			);
		}
	},
	{
		title: "快捷命令",
		render: (value: string) => {
			console.log("value", value);
			return (
				<Space>
					<Tooltip content={"adsadasdadasdasdadasdad"}>
						<Button>添加依赖</Button>
					</Tooltip>
					<Tooltip content={"adsadasdadasdasdadasdad"}>
						<Button>下载分支</Button>
					</Tooltip>
					<Tooltip content={"adsadasdadasdasdadasdad"}>
						<Button>import 路径</Button>
					</Tooltip>
					<Tooltip content={"adsadasdadasdasdadasdad"}>
						<Button>主结构体 import 路径</Button>
					</Tooltip>
				</Space>
			);
		}
	},
	{
		title: "操作",
		render: (value: string) => {
			console.log("value", value);
			return (
				<Space>
					<Button type="warning">强制更新仓库</Button>
					<Button type="danger">删除仓库</Button>
				</Space>
			);
		}
	}
];

export default function RepositoryPage() {
	const data = getIDLsRes();
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	const [modal, contextHolder] = Modal.useModal();
	const config = {
		size: "medium",
		title: "添加仓库",
		content: (
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
						仓库 URL
					</div>
					<Input
						style={{
							width: "30rem"
						}}
						showClear
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
						TOKEN
					</div>
					<Input
						style={{
							width: "30rem"
						}}
						showClear
					></Input>
				</Space>
			</Space>
		),
		cancelText: "取消",
		okText: "确定",
		icon: <IconInfoCircle />,
		onOk: () => {
			// 返回一个延时的 Promise
			return new Promise((resolve, reject) => {
				setTimeout(
					Math.random() > 0.5
						? () => {
								Toast.success("添加成功！");
								resolve(true);
						  }
						: () => {
								Toast.error("Oops errors!");
								reject(false);
						  },
					1000
				);
			}).catch(() => console.log("Oops errors!"));
		}
	} as ModalReactProps;

	const fetchData = async (currentPage = 1) => {
		setLoading(true);
		setPage(currentPage);
		const curDataSource = await new Promise((res) => {
			setTimeout(() => {
				const data = getIDLsRes();
				const dataSource = data.slice(
					(currentPage - 1) * pageSize,
					currentPage * pageSize
				);
				res(dataSource);
			}, 300);
		});
		setLoading(false);
		setData(curDataSource);
	};

	const handlePageChange = (page: number) => {
		fetchData(page);
	};

	useEffect(() => {
		fetchData();
	}, []);

	return (
		<ConfigProvider locale={en_GB}>
			<div>
				<AllInfoBox type={"repo"} />
				<div
					style={{
						padding: "1rem 0"
					}}
				>
					<Space
						style={{
							display: "flex",
							justifyContent: "space-between"
						}}
					>
						<Space>
							<Input
								style={{
									width: "20rem"
								}}
								prefix={<IconSearch />}
								showClear
							></Input>
							<Button type="primary" htmlType="submit">
								提交
							</Button>
							<div
								style={{
									color: "var(--semi-color-text-2)"
								}}
							>
								（最多显示 10 条数据）
							</div>
						</Space>
						<Button
							style={{
								width: "10rem"
							}}
							onClick={() => {
								console.log("modal", modal);
								modal.confirm(config);
							}}
						>
							添加仓库
						</Button>
					</Space>
				</div>
				<div className={styles["content"]}>
					<Table
						columns={columns}
						dataSource={dataSource as Data[]}
						pagination={{
							currentPage,
							pageSize: 10,
							total: data.length,
							onPageChange: handlePageChange
						}}
						loading={loading}
					/>
				</div>
			</div>
			{contextHolder}
		</ConfigProvider>
	);
}
