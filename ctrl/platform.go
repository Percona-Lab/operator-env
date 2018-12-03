package ctrl

// Platform represents enumeration of container orchestration platform type.
type Platform int

const (
	// Kubernetes represents Kubernetes test environment.
	Kubernetes Platform = iota
	// OpenShift represents OpenShift test environment.
	OpenShift
)

// String return string representation of concrete Platform type
func (p Platform) String() string {
	platforms := [...]string{
		"kubernetes",
		"openshift",
	}
	return platforms[p]
}
