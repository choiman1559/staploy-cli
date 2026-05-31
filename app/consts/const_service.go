package consts

const (
	APIRouteSchema = "/api/%s/%s" // INFO: Schema => /api/{version}/{connection_type}
	ConnTypeAdmin  = "admin"
	ConnTypeWorker = "worker"
	StatusError    = "error"
	StatusOK       = "ok"
)

const (
	InfoPrefix    = "[+] "
	SkipPrefix    = "[-] "
	ErrorPrefix   = "[!] "
	ProcessPrefix = "[*] "
	WarningPrefix = "[?] "
)

const PACKAGE_DIR_SHARE = "share"
const PACKAGE_FILE_METADATA = ".metadata"
const PACKAGE_FORMAT_VERSION = 1

const BLOB_REQ_TYPE = "blob_req_type"
const BLOB_REQ_TYPE_UPLOAD = "type_upload"
const BLOB_REQ_TYPE_DOWNLOAD = "type_download"
const BLOB_FIELD_PACKAGE = "PACKAGE_UPLOAD"
