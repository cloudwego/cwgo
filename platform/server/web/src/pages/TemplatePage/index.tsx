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
import styles from "./index.module.scss";
import en_GB from "@douyinfe/semi-ui/lib/es/locale/source/en_GB";
import InnerTable from "./components/InnerTable";
import { ModalReactProps } from "@douyinfe/semi-ui/lib/es/modal";
import { Data } from "@douyinfe/semi-ui/lib/es/cascader";
// import AllInfoBox from "../../components/AllInfoBox/index";

const pageSize = 10;

export default function TemplatePage() {
	const data = getIDLsRes();
	const [dataSource, setData] = useState<unknown>([]);
	const [loading, setLoading] = useState(false);
	const [currentPage, setPage] = useState(1);
	const [modal, contextHolder] = Modal.useModal();
	const config = {
		size: "medium",
		title: "添加模版",
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
						模版名称
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
						模版类型
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
	const itemConfig = {
		width: "80%",
		title: "模版元素列表",
		content: (
			<div>
				<InnerTable />
			</div>
		),
		cancelText: "取消",
		okText: "确定",
		icon: null
	};

	const columns = [
		{
			title: "模版名称",
			dataIndex: "name",
			width: 200,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "类型",
			dataIndex: "type",
			width: 200,
			render: (value: number) => {
				return value === 1 ? (
					<Tag color="purple" size="large">
						hertz
					</Tag>
				) : (
					<Tag color="blue" size="large">
						kitex
					</Tag>
				);
			}
		},
		{
			title: "创建时间",
			dataIndex: "create_time",
			width: 300,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "更新时间",
			dataIndex: "update_time",
			width: 300,
			render: (value: string) => {
				return <div>{value}</div>;
			}
		},
		{
			title: "操作",
			render: (value: string) => {
				console.log("value", value);
				return (
					<Space>
						<Button
							onClick={() => {
								console.log("modal", modal);
								modal.confirm(itemConfig);
							}}
						>
							查看 / 更新模版
						</Button>
					</Space>
				);
			}
		}
	];

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
				<div
					style={{
						paddingBottom: "1rem"
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
							添加模版
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
