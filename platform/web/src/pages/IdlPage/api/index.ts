import { IdlRes } from "../../../types";
import idlsResData from "./mock.json";

function getIDLsRes() {
	return idlsResData as IdlRes[];
}

export { getIDLsRes };
