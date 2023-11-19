import service from "../../../utils/request";

async function getRepo(currentPage = 1, pageSize = 10) {
	const repoData = await service.get("/api/repo", {
		params: {
			page: currentPage,
			limit: pageSize,
			order: 0,
			order_by: "create_time"
		}
	});
	return repoData.data;
}

async function createRepo(
	repoType: number,
	url: string,
	token: string,
	storeType: number
) {
	const repoData: { msg: string } = await service.post("/api/repo", {
		repository_type: repoType,
		repository_url: url,
		token: token,
		store_type: storeType
	});
	return repoData.msg;
}

async function updateRepo(id: number, token: string, status: number) {
	const repoData: { msg: string } = await service.patch(`/api/repo`, {
		id: id,
		token: token,
		status: status
	});
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

export { getRepo, createRepo, updateRepo, deleteRepo, syncRepo };
