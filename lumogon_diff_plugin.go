package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/pkg/archive"
	"github.com/puppetlabs/lumogon/capabilities/payloadfilter"
	"github.com/puppetlabs/lumogon/dockeradapter"
	"github.com/puppetlabs/lumogon/logging"
	"github.com/puppetlabs/lumogon/plugin"
	"github.com/puppetlabs/lumogon/types"
)

func main() {
	fmt.Println("Lumogon Diff Plugin")
}

// LP TODO
type LP struct{}

// Metadata TODO
func (p LP) Metadata() *plugin.Metadata {
	return &plugin.Metadata{
		Schema:      "http://puppet.com/lumogon/capability/diff/draft-01/schema#1",
		ID:          "diff",
		Name:        "Changed Files",
		Description: `The diff capability returns files changed from the initial image as a map["changed file"]"change type"`,
		Type:        plugin.DockerAPI,
		Version:     "0.0.1",
		GitSHA:      "yoIHeardYouLikeGitSHAs",
		SupportedOS: map[string]int{"all": 1},
	}
}

// Harvest TODO
func (p LP) Harvest(client dockeradapter.Harvester, ID string, target types.TargetContainer) (map[string]interface{}, error) {
	logging.Debug("[Plugin Diff] Harvesting diff from %s [%s]", target.Name, target.ID)

	ctx := context.Background()

	changedFiles, err := getChangedFiles(ctx, client, ID, target)
	if err != nil {
		logging.Debug("[Plugin Diff] Error getting changed files: %v", err)
		return nil, err
	}

	filtered, err := payloadfilter.Filter(changedFiles)
	if err != nil {
		logging.Debug("[Plugin Diff] Error filtering changedFiles output: %v", changedFiles)
		return nil, err
	}
	return filtered, nil
}

// Impl TODO
var Impl plugin.Plugin = LP{}

// Capability: types.Capability{
// 	Schema:      "http://puppet.com/lumogon/capability/diff/draft-01/schema#1",
// 	Title:       "Changed Files",
// 	Name:        "diff",
// 	Description: diffDescription,
// 	Type:        "dockerapi",
// 	Payload:     nil,
// 	SupportedOS: map[string]int{"all": 1},

func getChangedFiles(ctx context.Context, client dockeradapter.Diff, id string, target types.TargetContainer) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	diffs, err := client.ContainerDiff(ctx, target.ID)
	if err != nil {
		errorMsg := fmt.Sprintf("[Plugin Diff] Error getting diff from targetContainer: %s, error: %s", target.Name, err)
		logging.Debug(errorMsg)
		return nil, err
	}

	for _, diff := range diffs {
		logging.Debug("[Plugin Diff]   Path: %s, Kind %d", diff.Path, diff.Kind)
		var kind string
		switch diff.Kind {
		case archive.ChangeModify:
			kind = "Modified"
		case archive.ChangeAdd:
			kind = "Added"
		case archive.ChangeDelete:
			kind = "Deleted"
		}
		result[diff.Path] = kind
	}
	return result, nil
}
