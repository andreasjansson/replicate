package cli

import (
	"github.com/spf13/cobra"

	"replicate.ai/cli/pkg/console"
	"replicate.ai/cli/pkg/project"
)

func newRmCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm <experiment or checkpoint ID> [experiment or checkpoint ID...]",
		Short: "Remove experiments or checkpoint",
		Long: `Remove experiments or checkpoints.

To remove experiments or checkpoints, pass any number of IDs (or prefixes).
`,
		Run:        handleErrors(removeExperimentOrCheckpoint),
		Args:       cobra.MinimumNArgs(1),
		Aliases:    []string{"delete"},
		SuggestFor: []string{"remove"},
		Example: `Delete an experiment and its checkpoints
(where a1b2c3d4 is an experiment ID):
replicate rm a1b2c3d4

Delete all experiments where the metric "val_accuracy" is less
than 0.2 at the best checkpoints:
replicate rm $(replicate ls --filter "val_accuracy < 0.2")
`,
	}

	addStorageURLFlag(cmd)

	return cmd
}

func removeExperimentOrCheckpoint(cmd *cobra.Command, prefixes []string) error {
	storageURL, projectDir, err := getStorageURLFromFlagOrConfig(cmd)
	if err != nil {
		return err
	}
	store, err := getStorage(storageURL, projectDir)
	if err != nil {
		return err
	}
	proj := project.NewProject(store)
	if err != nil {
		return err
	}

	for _, prefix := range prefixes {
		comOrExp, err := proj.CheckpointOrExperimentFromPrefix(prefix)
		if err != nil {
			return err
		}
		if comOrExp.Checkpoint != nil {
			console.Info("Removing checkpoint %s...", comOrExp.Checkpoint.ShortID())
			if err := proj.DeleteCheckpoint(comOrExp.Checkpoint); err != nil {
				return err
			}
		} else {
			console.Info("Removing experiment %s and its checkpoints...", comOrExp.Experiment.ShortID())
			experiment := comOrExp.Experiment
			for _, checkpoint := range experiment.Checkpoints {
				if err := proj.DeleteCheckpoint(checkpoint); err != nil {
					return err
				}
			}
			if err := proj.DeleteExperiment(experiment); err != nil {
				return err
			}
		}
	}

	return nil
}
