package v1

const (
	Service               ResourceNames = "Service"
	Deployment            ResourceNames = "Deployment"
	Secret                ResourceNames = "Secret"
	ConfigMap             ResourceNames = "ConfigMap"
	StatefulSet           ResourceNames = "StatefulSet"
	Job                   ResourceNames = "Job"
	CronJob               ResourceNames = "CronJob"
	PersistantVolume      ResourceNames = "PersistantVolume"
	PersistantVolumeClaim ResourceNames = "PersistantVolumeClaim"
	ServiceAccount        ResourceNames = "ServiceAccount"
	Role                  ResourceNames = "Role"
	RoleBinding           ResourceNames = "RoleBinding"
	ClusterRole           ResourceNames = "ClusterRole"
	ClusterRoleBinding    ResourceNames = "ClusterRoleBinding"
	NetworkPolicy         ResourceNames = "NetworkPolicy"
	LimitRange            ResourceNames = "LimitRange"
	ResourceQuota         ResourceNames = "ResourceQuota"
	Namespaces            ResourceNames = "Namespaces"
)

const (
	AWS   CloudName = "AWS"
	GCP   CloudName = "GCP"
	Azure CloudName = "Azure"
)

const (
	Low      SwipePolicyName = "low"
	Moderate SwipePolicyName = "moderate"
	High     SwipePolicyName = "high"
)

const (
	Serve   OperationName = "SERVE"
	CleanUp OperationName = "CLEANUP"
)
