import service from "../../../utils/request";

async function getToken(
	currentPage = 1,
	pageSize = 10,
	repository_domain = ""
) {
	const tokenData = await service.get("/api/token", {
		params: {
			page: currentPage,
			limit: pageSize,
			order: 0,
			order_by: "create_time",
			repository_domain
		}
	});
	return tokenData.data;
}

async function createIdl(
	repository_type: number,
	repository_domain: string,
	token: string
) {
	console.log(repository_type, repository_domain, token);
	const repoData: { msg: string } = await service.post("/api/token", {
		repository_type,
		repository_domain,
		token
	});
	return repoData.msg;
}

async function deleteToken(id: number) {
	const repoData: { msg: string } = await service.delete(`/api/token`, {
		data: {
			ids: [id]
		}
	});
	return repoData.msg;
}

export { getToken, createIdl, deleteToken };
