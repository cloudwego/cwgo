import { repoFetchOption } from "../../../types";
import service from "../../../utils/request";

async function getRepo(fetchOption: repoFetchOption) {
	const repoData = await service.get("/api/repo", {
		params: {
			store_type: fetchOption.repoType === "idl" ? 1 : 2,
			page: fetchOption.currentPage,
			limit: fetchOption.pageSize,
			order: 0,
			order_by: "create_time",
			repository_domain: fetchOption.repository_domain,
			repository_name: fetchOption.repository_name
		}
	});
	return repoData.data;
}

async function createRepo(repoType: number, url: string, token: string) {
	const repoData: { msg: string } = await service.post("/api/repo", {
		repository_type: repoType,
		repository_url: url,
		token: token,
		store_type: 2
	});
	return repoData.msg;
}

async function updateRepo(id: number, branch: string, status: number) {
	let params;
	if (branch) {
		params = {
			id,
			branch
		};
	} else {
		params = {
			id,
			status
		};
	}
	const repoData: { msg: string } = await service.patch(`/api/repo`, params);
	return repoData.msg;
}

async function syncRepo(id: number) {
	const repoData: { msg: string } = await service.post(`/api/repo/sync`, {
		ids: [id]
	});
	return repoData.msg;
}

async function deleteRepo(id: number) {
	const repoData: { msg: string } = await service.delete(`/api/repo`, {
		data: {
			ids: [id]
		}
	});
	return repoData.msg;
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

export { getRepo, createRepo, updateRepo, deleteRepo, syncRepo, createIdl };
