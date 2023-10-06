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
	Toast
} from "@douyinfe/semi-ui";
import { IconInfoCircle, IconSearch } from "@douyinfe/semi-icons";
import { IdlRes } from "../../types";
import AllInfoBox from "../../components/AllInfoBox";
import styles from "./index.module.scss";
import en_GB from "@douyinfe/semi-ui/lib/es/locale/source/en_GB";
import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";

const pageSize = 10;

const columns = [
	{
		title: "PSM",
		dataIndex: "service_name",
		width: "10rem",
		render: (value: string) => {
			return <div>{value}</div>;
		}
	},
	{
		title: "IDL 仓库",
		dataIndex: "url",
		width: "15rem",
		render: (value: string) => {
			return <a href="#">{value}</a>;
		}
	},
	{
		title: "主 IDL 文件路径",
		dataIndex: "main_idl_path",
		width: "15rem",
		render: (value: string) => {
			return <a href="#">{value}</a>;
		}
	},
	{
		title: "最近一次同步时间",
		dataIndex: "last_sync_time",
		width: "20rem",
		render: (value: string) => {
			return <div>{value}</div>;
		}
	},
	{
		title: "状态",
		dataIndex: "status",
		width: "10rem",
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
		title: "操作",
		render: (value: string) => {
			console.log("value", value);
			return (
				<Space>
					<Button type="warning">强制同步 IDL 信息</Button>
					<Button>分支生成</Button>
					<Button type="danger">删除 IDL</Button>
				</Space>
			);
		}
	}
];

export default function IdlPage() {
	const data = getIDLsRes();
	const [dataSource, setData] = useState<IdlRes[]>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	const [modal, contextHolder] = Modal.useModal();
	const config = {
		size: "medium",
		title: "添加 IDL",
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
							width: "8rem"
						}}
					>
						IDL 仓库 ID
					</div>
					<Input
						style={{
							width: "25rem"
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
							width: "8rem"
						}}
					>
						主 IDL 文件路径
					</div>
					<Input
						style={{
							width: "25rem"
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
							width: "8rem"
						}}
					>
						服务名
					</div>
					<Input
						style={{
							width: "25rem"
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
		const curDataSource: IdlRes[] = await new Promise((res) => {
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
				<AllInfoBox type={"idl"} />
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
							添加 IDL
						</Button>
					</Space>
				</div>
				<div className={styles["content"]}>
					<Table
						columns={columns}
						dataSource={dataSource}
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
