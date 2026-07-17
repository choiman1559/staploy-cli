package build

type StaployFile struct {
	Config  Configs   `hcl:"configure,block"`
	Targets []*Target `hcl:"target,block"`
	Manages []*Manage `hcl:"manage,block"`
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
	Version    string `hcl:"version"`
	Repository string `hcl:"repository"`
}
