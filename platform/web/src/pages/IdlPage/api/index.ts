import service from "../../../utils/request";

async function getIdl(currentPage = 1, pageSize = 10) {
	const idlData = await service.get("/api/idl", {
		params: {
			page: currentPage,
			limit: pageSize,
			order: 0,
			order_by: "create_time"
		}
	});
	return idlData.data.idls;
}

async function createIdl(
	id: number,
	idlPath: string,
	serviceName: string,
	serviceRepoName: string
) {
	const repoData: { msg: string } = await service.post("/api/idl", {
		repository_id: id,
		main_idl_path: idlPath,
		service_name: serviceName,
		service_repository_name: serviceRepoName
	});
	return repoData.msg;
}

async function updateIdl(id: number, service_name: string) {
	console.log(id, service_name);
	const repoData: { msg: string } = await service.post(`/api/idl/sync`, {
		ids: [id]
		// service_name: service_name
	});
	return repoData.msg;
}

async function deleteIdl(id: number) {
	const repoData: { msg: string } = await service.delete(`/api/idl`, {
		data: {
			ids: [id]
		}
	});
	return repoData.msg;
}

export { getIdl, createIdl, updateIdl, deleteIdl };
