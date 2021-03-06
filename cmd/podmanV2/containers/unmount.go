package containers

import (
	"fmt"

	"github.com/containers/libpod/cmd/podmanV2/parse"
	"github.com/containers/libpod/cmd/podmanV2/registry"
	"github.com/containers/libpod/cmd/podmanV2/utils"
	"github.com/containers/libpod/pkg/domain/entities"
	"github.com/spf13/cobra"
)

var (
	description = `Container storage increments a mount counter each time a container is mounted.

  When a container is unmounted, the mount counter is decremented. The container's root filesystem is physically unmounted only when the mount counter reaches zero indicating no other processes are using the mount.

  An unmount can be forced with the --force flag.
`
	umountCommand = &cobra.Command{
		Use:     "umount [flags] CONTAINER [CONTAINER...]",
		Aliases: []string{"unmount"},
		Short:   "Unmounts working container's root filesystem",
		Long:    description,
		RunE:    unmount,
		PreRunE: preRunE,
		Args: func(cmd *cobra.Command, args []string) error {
			return parse.CheckAllLatestAndCIDFile(cmd, args, false, false)
		},
		Example: `podman umount ctrID
  podman umount ctrID1 ctrID2 ctrID3
  podman umount --all`,
	}
)

var (
	unmountOpts entities.ContainerUnmountOptions
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode},
		Command: umountCommand,
	})
	flags := umountCommand.Flags()
	flags.BoolVarP(&unmountOpts.All, "all", "a", false, "Umount all of the currently mounted containers")
	flags.BoolVarP(&unmountOpts.Force, "force", "f", false, "Force the complete umount all of the currently mounted containers")
	flags.BoolVarP(&unmountOpts.Latest, "latest", "l", false, "Act on the latest container podman is aware of")
}

func unmount(cmd *cobra.Command, args []string) error {
	var errs utils.OutputErrors
	reports, err := registry.ContainerEngine().ContainerUnmount(registry.GetContext(), args, unmountOpts)
	if err != nil {
		return err
	}
	for _, r := range reports {
		if r.Err == nil {
			fmt.Println(r.Id)
		} else {
			errs = append(errs, r.Err)
		}
	}
	return errs.PrintErrors()
}
