package consts

//goland:noinspection GoUnusedConst
const (
	APIRouteSchema   = "/api/%s/%s" // INFO: Schema => /api/{version}/{connection_type}
	ConnTypeAdmin    = "admin"
	ConnTypeRegistry = "registry"
	ConnTypeWorker   = "worker"
	StatusError      = "error"
	StatusOK         = "ok"
	StatusNone       = "none"
)

const (
	InfoPrefix    = "[+]"
	SkipPrefix    = "[-]"
	ErrorPrefix   = "[!]"
	ProcessPrefix = "[*]"
	TaskPrefix    = "[>]"
	HookPrefix    = "[~]"
	WarningPrefix = "[?]"
)

const PACKAGE_DIR_SHARE = "share"
const PACKAGE_FILE_METADATA = ".metadata"
const PACKAGE_FORMAT_VERSION = 1

const STAFILE_SHELL_PREFIX = "shell:"
const STAFILE_GROUP_PREFIX = "group:"
const STAFILE_ALIAS_PREFIX = "alias:"

const BLOB_REQ_TYPE = "blob_req_type"
const BLOB_REQ_TYPE_UPLOAD = "type_upload"

//goland:noinspection GoUnusedConst
const BLOB_REQ_TYPE_DOWNLOAD = "type_download"
const BLOB_FIELD_PACKAGE = "PACKAGE_UPLOAD"
