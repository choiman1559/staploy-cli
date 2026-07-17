package cmds

// ### TaskAppsTypes
// create 		(TYPE_APP_REGISTER)
// delete 		(TYPE_APP_DELETE)
// apps 		(TYPE_APP_LISTS)
// upload 		(TYPE_APP_BLOB)

type CreateCmd struct {
	AppName        string `arg:"-n,--app-name,required" help:"name of the application (warning: cannot be changed once created)"`
	AppDescription string `arg:"-d,--description" help:"Specifies the description of the app"`
}

type DeleteCmd struct {
	AppName     string   `arg:"-n,--app-name,required" help:"name of the application to delete"`
	VersionName []string `arg:"-e,--version-name" help:"names of version to delete"`
}

type AppsCmd struct {
	AppName string `arg:"-n,--app-name" help:"name of the application to list"`
}

type UploadCmd struct {
	PackageFile string `arg:"-f,--file,required" help:"package file to upload"`
}

// ### TaskNodeTypes
// list 		(TYPE_NODE_CONNECTED)
// fetch 		(TYPE_NODE_REQ_WORKER_INFO, TYPE_NODE_REQ_APP_INFO)
// bash 		(TYPE_NODE_EXECUTE_SHELL)
// disconn 		(TYPE_NODE_DISCONN_WORKER)

type ListCmd struct {
	WorkerId []string `arg:"-w,--worker-id" help:"list of worker ids to search"`
	Refresh  bool     `arg:"-r,--refresh" help:"request to refresh data from worker"`
	Detail   bool     `arg:"-d" help:"show all detail of workers"`
}

type FetchCmd struct {
	WorkerId    string   `arg:"-w,--worker-id,required" help:"worker id to use for command"`
	AppName     string   `arg:"-n,--app-name" help:"name of app to fetch from worker"`
	VersionName []string `arg:"-e,--version-name" help:"names of version to fetch from worker"`
	Detail      bool     `arg:"-d" help:"show all detail of versions"`
}

type BashCmd struct {
	Command  string `arg:"positional,required" help:"shell command to execute"`
	WorkerId string `arg:"-w,--worker-id,required" help:"worker id to use for command"`
	NoOutput bool   `arg:"-o,--no-output" help:"don't waiting for output and exit"`
}

type DisconnCmd struct {
	WorkerId []string `arg:"positional,required" help:"list of worker ids to disconnect"`
}

/// ### TaskDeployTypes
// push 		(TYPE_DEPLOY_PUSH_VERSION)
// remove 		(TYPE_DEPLOY_DEL_VERSION)
// set 			(TYPE_DEPLOY_SET_VERSION)

type PushCmd struct {
	WorkerId []string `arg:"-w,--worker-id,required" help:"worker id to install new package"`
	AppName  string   `arg:"-n,--app-name, required" help:"name of app to push"`
	Version  string   `arg:"-e,--version,required" help:"version of app to push"`
}

type RemoveCmd struct {
	WorkerId   []string `arg:"-w,--worker-id,required" help:"worker id to remove package"`
	AppName    string   `arg:"-n,--app-name" help:"name of app to remove"`
	Version    string   `arg:"-e,--version" help:"version of app to remove, remove all versions if not specified"`
	AutoRemove bool     `arg:"--autoremove" help:"automatically remove unused versions"`
}

type SetCmd struct {
	WorkerId []string `arg:"-w,--worker-id,required" help:"worker id to set package"`
	AppName  string   `arg:"-n,--app-name, required" help:"name of app to set"`
	Version  string   `arg:"-e,--version" help:"version of app to set, disable package if not specified"`
}

/// ### TaskGroupTypes
// group
// -> create 	(TYPE_GROUP_CREATE)
// -> delete	(TYPE_GROUP_DELETE)
// -> add		(TYPE_GROUP_ADD_WORKER)
// -> remove	(TYPE_GROUP_REMOVE_WORKER)
// -> list		(TYPE_QUERY_GROUP_LIST)

type GroupCmd struct {
	GroupCreate *GroupCreateCmd `arg:"subcommand:create" help:"create group"`
	GroupDelete *GroupDeleteCmd `arg:"subcommand:delete" help:"delete group"`
	GroupList   *GroupListCmd   `arg:"subcommand:list" help:"list groups"`
	GroupAdd    *GroupAddCmd    `arg:"subcommand:add" help:"add workers to group"`
	GroupRemove *GroupRemoveCmd `arg:"subcommand:remove" help:"remove workers from group"`
}

type GroupCreateCmd struct {
	GroupName string `arg:"-g,--group-name,required" help:"name of group to create"`
}

type GroupDeleteCmd struct {
	GroupName string `arg:"-g,--group-name,required" help:"name of group to delete"`
}

type GroupListCmd struct {
	GroupName string `arg:"-g,--group-name" help:"name of group to list"`
}

type GroupAddCmd struct {
	GroupName string   `arg:"-g,--group-name,required" help:"name of group to add workers"`
	WorkerId  []string `arg:"-w,--worker-id,required" help:"id or name of worker for add to group"`
}

type GroupRemoveCmd struct {
	GroupName string   `arg:"-g,--group-name,required" help:"name of group to remove workers"`
	WorkerId  []string `arg:"-w,--worker-id,required" help:"id or name of worker to remove from group"`
}

/// ### TaskBuilds
// build	(package building)
// file 	(eval *.Staployfile)

type BuildCmd struct {
	I386     string `arg:"--i386" help:"specifies directory of i386 executable"`
	X86_64   string `arg:"--x86_64" help:"specifies directory of x86_64 executable"`
	Arm      string `arg:"--arm" help:"specifies directory of ARM executable"`
	Aarch64  string `arg:"--aarch64" help:"specifies directory of aarch64 executable"`
	Riscv32  string `arg:"--riscv32" help:"specifies directory of riscv32 executable"`
	Riscv64  string `arg:"--riscv64" help:"specifies directory of riscv64 executable"`
	Mipsel   string `arg:"--mipsel" help:"specifies directory of mipsel executable"`
	Mips64el string `arg:"--mips64el" help:"specifies directory of mips64el executable"`
	Share    string `arg:"--share" help:"specifies directory of shared executable and resources"`

	Executable  []string `arg:"-e,separate" help:"specifies names of executable and resources that will be linked on target worker"`
	AppName     string   `arg:"-m,--app-name,required" help:"specifies package name"`
	VersionName string   `arg:"-n,--version-name,required" help:"specifies package version name"`
	LibVersion  string   `arg:"--lib-version" help:"specifies package lib version"`
	OutputDir   string   `arg:"-o,--output-dir" help:"specifies directory of output"`
}

type StaFileCmd struct {
	ConfigFile string `arg:"-f,--file,required" help:"Staployfile to configure"`
	ParseOnly  bool   `arg:"--parse-only" help:"Only parse and analyze staployfile for debugging"`
}

type UserCmd struct {
	UserLogin  *UserLoginCmd       `arg:"subcommand:login" help:"login to user"`
	UserCreate *UserCreateCmd      `arg:"subcommand:create" help:"create a user"`
	UserRemove *UserRemoveCmd      `arg:"subcommand:remove" help:"remove a user"`
	UserPerm   *UserPermissionsCmd `arg:"subcommand:perm" help:"perm permission to user"`
	UserList   *UserListCmd        `arg:"subcommand:list" help:"list users"`
	UserAudit  *UserAuditCmd       `arg:"subcommand:audit" help:"audit the user"`
}

type UserLoginCmd struct {
	UserName       string `arg:"-n,--user-name,required" help:"user name"`
	PrintTokenOnly bool   `arg:"--print-token" help:"prints token only. useful when pipelining in script."`
}

type UserCreateCmd struct {
	UserName string `arg:"-n,--user-name,required" help:"user name"`
	Refresh  bool   `arg:"-r,--refresh" help:"refresh user token, not creating new user."`
}

type UserRemoveCmd struct {
	UserName string `arg:"-n,--user-name,required" help:"user name"`
}

type UserListCmd struct {
	UserName string `arg:"-n,--user-name" help:"user name to find, blank if get all users"`
}

type UserAuditCmd struct {
}

type CommandPermEnum string

const (
	Create  CommandPermEnum = "create"
	Delete  CommandPermEnum = "delete"
	Upload  CommandPermEnum = "upload"
	Bash    CommandPermEnum = "bash"
	Disconn CommandPermEnum = "disconn"
	Push    CommandPermEnum = "push"
	Remove  CommandPermEnum = "remove"
	Set     CommandPermEnum = "set"
	Group   CommandPermEnum = "group"
	Query   CommandPermEnum = "query"
	User    CommandPermEnum = "user"
	Admin   CommandPermEnum = "admin"
)

type UserPermissionsCmd struct {
	UserName string            `arg:"-n,--user-name,required" help:"user name"`
	PermCmd  []CommandPermEnum `arg:"-c,--perm-cmds,required" help:"perm commands to enable (commands: create, delete, upload, bash, disconn, push, remove, set, group, user)"`
}

type RegistryCmd struct {
	PushLocal   *RegistryPushLocalCmd   `arg:"subcommand:push-local" help:"push package to local registry"`
	RemoveLocal *RegistryRemoveLocalCmd `arg:"subcommand:remove-local" help:"remove package from local registry"`
	ListLocal   *RegistryListLocalCmd   `arg:"subcommand:list-local" help:"list of packages of local registry"`

	Pull            *RegistryPullCmd            `arg:"subcommand:pull" help:"pull package from remote registry"`
	ListRepo        *RegistryListRepoCmd        `arg:"subcommand:list-repo" help:"list of linked remote registry"`
	AddRepo         *RegistryAddRepoCmd         `arg:"subcommand:add-repo" help:"add a repository into server"`
	RemoveRepo      *RegistryRemoveRepoCmd      `arg:"subcommand:remove-repo" help:"remove a repository from server"`
	ManageRepoToken *RegistryManageRepoTokenCmd `arg:"subcommand:token" help:"associate login access token for remote repository"`
	List            *RegistryListPackageCmd     `arg:"subcommand:list" help:"list of available packages from remote registry"`
	UpdateCache     *RegistryUpdateCacheCmd     `arg:"subcommand:update" help:"update cache of packages from remote registry"`
}

type RegistryPushLocalCmd struct {
	PackageFile string `arg:"-f,--file,required" help:"package file to upload"`
}

type RegistryRemoveLocalCmd struct {
	AppName string `arg:"-n,--app-name, required" help:"name of app to remove"`
	Version string `arg:"-e,--version" help:"version of app to remove"`
}

type RegistryListLocalCmd struct {
	AppName string `arg:"-n,--app-name" help:"name of app to list"`
}

type RegistryPullCmd struct {
	AppName    string `arg:"-n,--app-name, required" help:"name of app to pull"`
	Version    string `arg:"-e,--version" help:"version of app to pull"`
	Repository string `arg:"--repository" help:"specify repository name to pull"`
}

type RegistryListRepoCmd struct {
}

type RegistryAddRepoCmd struct {
	RepoUrl []string `arg:"positional, required" help:"repository url to add"`
}

type RegistryRemoveRepoCmd struct {
	RepoUrl []string `arg:"positional, required" help:"repository url to remove"`
}

type RegistryManageRepoTokenCmd struct {
	RepoUrl string `arg:"-r,--repo-url, required" help:"repository url to add"`
	Token   string `arg:"-t,--token" help:"token to add to repo"`
}

type RegistryListPackageCmd struct {
	AppName string `arg:"-n,--app-name" help:"name of app to list"`
}

type RegistryUpdateCacheCmd struct {
	RepoUrl []string `arg:"positional" help:"repository url to update cache"`
}

var Args struct {
	Address   string `arg:"-a,env:STAPLOY_HOST_ADDR" help:"overrides server address. can be preset with STAPLOY_HOST_ADDR environment variable"`
	Port      int    `arg:"-p,env:STAPLOY_HOST_PORT" help:"overrides server port. can be preset with STAPLOY_HOST_PORT environment variable"`
	UseIdOnly bool   `arg:"--enforce-uuid,env:STAPLOY_USE_WORKER_ID" help:"enforce to use worker uuid only instead of name and groups"`
	NoColor   bool   `arg:"--no-color" help:"disables color output"`
	Verbose   bool   `arg:"-v,--verbose" help:"verbose output"`

	UserJwtToken   string `arg:"-u,--user-token,env:STAPLOY_USER_TOKEN" help:"user jwt token (Tip: get token using 'user login' command)"`
	DisableTls     bool   `arg:"--disable-tls,env:STAPLOY_TLS_DISABLE" help:"disable TLS when connecting"`
	SkipValidation bool   `arg:"--skip-validation,env:STAPLOY_TLS_NOCHECK" help:"skip validation of tls certificates"`

	Create *CreateCmd `arg:"subcommand:create" help:"create newly or update information of application"`
	Delete *DeleteCmd `arg:"subcommand:delete" help:"delete a application"`
	Apps   *AppsCmd   `arg:"subcommand:apps" help:"list all application or versions available on remote server"`
	Upload *UploadCmd `arg:"subcommand:upload" help:"upload a distributable package"`

	List    *ListCmd    `arg:"subcommand:list" help:"lists of connected workers"`
	Fetch   *FetchCmd   `arg:"subcommand:fetch" help:"fetch and update to db from a remote worker"`
	Bash    *BashCmd    `arg:"subcommand:bash" help:"execute shell command"`
	Disconn *DisconnCmd `arg:"subcommand:disconn" help:"disconnect workers"`

	Push   *PushCmd   `arg:"subcommand:push" help:"push an package to a remote worker"`
	Remove *RemoveCmd `arg:"subcommand:remove" help:"remove a package from a remote worker"`
	Set    *SetCmd    `arg:"subcommand:set" help:"set a version of package to executable path from a remote worker"`

	Group    *GroupCmd    `arg:"subcommand:group" help:"manage groups of workers"`
	Build    *BuildCmd    `arg:"subcommand:build" help:"build a distributable package"`
	File     *StaFileCmd  `arg:"subcommand:file" help:"run HCL-like staployfile for deploy & management workers as code"`
	User     *UserCmd     `arg:"subcommand:user" help:"login, create, manage a user"`
	Registry *RegistryCmd `arg:"subcommand:registry" help:"push or pull from a remote registry"`
}
