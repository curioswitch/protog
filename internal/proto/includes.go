package proto

var includeSpecs = []includeSpec{
	{
		prefix:  "google/api/",
		repo:    "github.com/googleapis/googleapis",
		ref:     "master",
		repoDir: "google",
		dir:     "google",
	},
	{
		prefix:  "google/rpc/",
		repo:    "github.com/googleapis/googleapis",
		ref:     "master",
		repoDir: "google",
		dir:     "google",
	},
	{
		prefix:  "gogoproto/",
		repo:    "github.com/gogo/protobuf",
		ref:     "master",
		repoDir: "gogoproto",
		dir:     "gogoproto",
	},
	{
		prefix: "k8s.io/api/",
		repo:   "github.com/kubernetes/api",
		ref:    "master",
		dir:    "k8s.io/api",
	},
	{
		prefix: "k8s.io/apimachinery/",
		repo:   "github.com/kubernetes/apimachinery",
		ref:    "master",
		dir:    "k8s.io/apimachinery",
	},
	{
		prefix:  "validate/",
		repo:    "github.com/envoyproxy/protoc-gen-validate",
		ref:     "main",
		repoDir: "validate",
		dir:     "validate",
	},
}
