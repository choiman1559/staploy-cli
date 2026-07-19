package build

type StaployFile struct {
	Config  *Configs  `hcl:"configure,block"`
	Targets []*Target `hcl:"target,block"`
	Manages []*Manage `hcl:"manage,block"`
	Builds  []*Build  `hcl:"build,block"`
}

type Configs struct {
	Address   string `hcl:"address"`
	Port      int    `hcl:"port"`
	UseIdOnly bool   `hcl:"enforce_uuid,optional"`
}

type Target struct {
	TargetName string   `hcl:"name,label"`
	WorkerIds  []string `hcl:"workers"`

	Deploy  *DeployTask  `hcl:"deploy,block"`
	Push    *PushTask    `hcl:"push_only,block"`
	Set     *SetTask     `hcl:"set_active,block"`
	Remove  *RemoveTask  `hcl:"remove,block"`
	Disconn *DisconnTask `hcl:"disconnect,block"`
}

type DeployTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type PushTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type SetTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type RemoveTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version"`
	AutoRemove bool   `hcl:"autoremove,optional"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type DisconnTask struct {
	PreDeploy string `hcl:"pre_deploy,optional"`
}

type Manage struct {
	AppName string      `hcl:"name,label"`
	Create  *CreateTask `hcl:"create,block"`
	Upload  *UploadTask `hcl:"upload,block"`
	Delete  *DeleteTask `hcl:"delete,block"`
	Pull    *PullTask   `hcl:"pull,block"`
}

type CreateTask struct {
	AppDescription string `hcl:"description,optional"`
}

type UploadTask struct {
	PackageFile string `hcl:"path"`
}

type DeleteTask struct {
	Versions []string `hcl:"versions,optional"`
}

type PullTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version,optional"`
	Repository string `hcl:"repository,optional"`
}

type Build struct {
	AppName     string   `hcl:"name,label"`
	OutputDir   string   `hcl:"output_dir"`
	Version     string   `hcl:"version"`
	Executables []string `hcl:"executable"`

	LibVersion string `hcl:"lib_version,optional"`
	PreBuild   string `hcl:"pre_build,optional"`
	PostBuild  string `hcl:"post_build,optional"`

	Target_share    *BuildTarget `hcl:"share,block"`
	Target_i386     *BuildTarget `hcl:"i386,block"`
	Target_x86_64   *BuildTarget `hcl:"x86_64,block"`
	Target_arm      *BuildTarget `hcl:"arm,block"`
	Target_aarch64  *BuildTarget `hcl:"aarch64,block"`
	Target_riscv32  *BuildTarget `hcl:"riscv32,block"`
	Target_riscv64  *BuildTarget `hcl:"riscv64,block"`
	Target_mipsel   *BuildTarget `hcl:"mipsel,block"`
	Target_mips64el *BuildTarget `hcl:"mips64el,block"`
}

type BuildTarget struct {
	Path      string `hcl:"path"`
	PreBuild  string `hcl:"pre_build,optional"`
	PostBuild string `hcl:"post_build,optional"`
}
