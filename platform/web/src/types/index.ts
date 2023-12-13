export interface repoFetchOption {
	currentPage: number;
	pageSize: number;
	repoType: string;
	repository_domain?: string;
	repository_name?: string;
}

export interface IdlRes {
	id: number;
	repository_id: number;
	main_idl_path: string;
	content: string;
	service_name: string;
	last_sync_time: string;
	create_time: string;
	update_time: string;
	url: string;
}

export interface AddRepo {
	repository_type: number;
	repository_url: string;
	token: string;
	store_type: number;
}

export interface UpdateRepo {
	id: number;
	repository_branch: string;
	status: number;
}

export interface AddIdl {
	repository_id: number;
	main_idl_path: string;
	service_name: string;
	service_repository_name: string;
}
