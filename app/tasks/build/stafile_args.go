package build

type StaployFile struct {
	Config  *Configs  `hcl:"configure,block"`
	Alias   *Alias    `hcl:"alias,block"`
	Targets []*Target `hcl:"target,block"`
	Manages []*Manage `hcl:"manage,block"`
	Builds  []*Build  `hcl:"build,block"`
}

type Configs struct {
	Address   string `hcl:"address"`
	Port      string `hcl:"port"`
	UseIdOnly bool   `hcl:"enforce_uuid,optional"`
}

type Alias struct {
	Worker []*WorkerAlias `hcl:"worker,block"`
	App    []*AppAlias    `hcl:"app,block"`
}

type WorkerAlias struct {
	Alias     string   `hcl:"alias,label"`
	WorkerIds []string `hcl:"workers"`
	Where     *Where   `hcl:"where,block"`
}

type AppAlias struct {
	Alias   string `hcl:"alias,label"`
	AppName string `hcl:"name"`
	Version string `hcl:"version,optional"`
}

type Target struct {
	TargetName string   `hcl:"name,label"`
	WorkerIds  []string `hcl:"workers"`
	Where      *Where   `hcl:"where,block"`

	Deploy  *DeployTask  `hcl:"deploy,block"`
	Push    *PushTask    `hcl:"push_only,block"`
	Set     *SetTask     `hcl:"set_active,block"`
	Unset   *UnsetTask   `hcl:"unset,block"`
	Remove  *RemoveTask  `hcl:"remove,block"`
	Disconn *DisconnTask `hcl:"disconnect,block"`
}

type Where struct {
	Name       string   `hcl:"name,optional"`
	Arch       []string `hcl:"arch,optional"`
	Cpu        string   `hcl:"cpu,optional"`
	Memory     string   `hcl:"memory,optional"`
	WorkingDir string   `hcl:"workdir,optional"`

	ShellEnabled       *bool    `hcl:"shell_enable,optional"`
	SkipIntegrityCheck *bool    `hcl:"integrity_skip,optional"`
	CpuBigEndian       *bool    `hcl:"cpu_big_endian,optional"`
	CpuCapabilities    []string `hcl:"cpu_capability,optional"`
	CpuFlags           []string `hcl:"cpu_flags,optional"`
}

type DeployTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version,optional"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type PushTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version,optional"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type SetTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version,optional"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type UnsetTask struct {
	AppName    string `hcl:"name,label"`
	PreDeploy  string `hcl:"pre_deploy,optional"`
	PostDeploy string `hcl:"post_deploy,optional"`
}

type RemoveTask struct {
	AppName    string `hcl:"name,label"`
	Version    string `hcl:"version,optional"`
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
	PackageFile string `hcl:"path,optional"`
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
	AppName      string    `hcl:"name,label"`
	OutputDir    string    `hcl:"output_dir"`
	Version      string    `hcl:"version,optional"`
	Executables  []string  `hcl:"executable"`
	Environments *[]string `hcl:"envs"`

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
	Target_mips     *BuildTarget `hcl:"mips,block"`
	Target_mips64   *BuildTarget `hcl:"mips64,block"`
	Target_ppc64    *BuildTarget `hcl:"ppc64,block"`
	Target_ppc64le  *BuildTarget `hcl:"ppc64le,block"`
	Target_s390x    *BuildTarget `hcl:"s390x,block"`
	Target_loong64  *BuildTarget `hcl:"loong64,block"`
}

type BuildTarget struct {
	Path      string `hcl:"path"`
	PreBuild  string `hcl:"pre_build,optional"`
	PostBuild string `hcl:"post_build,optional"`
}
