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

var Args struct {
	Address string `arg:"-a,env:STAPLOY_HOST_ADDR" help:"overrides server address. can be preset with STAPLOY_HOST_ADDR environment variable"`
	Port    int    `arg:"-p,env:STAPLOY_HOST_PORT" help:"overrides server port. can be preset with STAPLOY_HOST_PORT environment variable"`
	UseName bool   `arg:"--use-name,env:STAPLOY_USE_WORKER_NAME" help:"allows using name instead of worker uuid"`
	Verbose bool   `arg:"-v,--verbose" help:"verbose output"`

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
	Build  *BuildCmd  `arg:"subcommand:build" help:"build a distributable package"`
}
