package conda

import (
	"time"

	"github.com/robocorp/rcc/common"
)

func Cleanup(daylimit int, dryrun, all bool) error {
	deadline := time.Now().Add(-24 * time.Duration(daylimit) * time.Hour)
	for _, template := range TemplateList() {
		whenLive, err := LastUsed(LiveFrom(template))
		if err != nil {
			return err
		}
		if !all && whenLive.After(deadline) {
			continue
		}
		whenBase, err := LastUsed(TemplateFrom(template))
		if err != nil {
			return err
		}
		if !all && whenBase.After(deadline) {
			continue
		}
		if dryrun {
			common.Log("Would be removing %v.", template)
			continue
		}
		RemoveEnvironment(template)
		common.Debug("Removed environment %v.", template)
	}
	return nil
}
