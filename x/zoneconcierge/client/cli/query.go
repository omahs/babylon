package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/babylonchain/babylon/x/zoneconcierge/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group zoneconcierge queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdChainsInfo())
	cmd.AddCommand(CmdFinalizedChainsInfo())
	cmd.AddCommand(CmdEpochChainsInfoInfo())
	return cmd
}

func CmdChainsInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chains-info <chain-ids>",
		Short: "retrieve the latest info for a given list of chains",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			req := types.QueryChainsInfoRequest{ChainIds: args}
			resp, err := queryClient.ChainsInfo(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdFinalizedChainsInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finalized-chains-info <chain-ids>",
		Short: "retrieve the finalized info for a given list of chains",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)
			req := types.QueryFinalizedChainsInfoRequest{ChainIds: args}
			resp, err := queryClient.FinalizedChainsInfo(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdEpochChainsInfoInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epoch-chains-info <epoch-num> <chain-ids>",
		Short: "retrieve the latest info for a list of chains in a given epoch",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryClient(clientCtx)

			var epoch uint64
			_, err := fmt.Sscanf(args[0], "%d", &epoch)
			if err != nil {
				return err
			}
			req := types.QueryEpochChainsInfoRequest{EpochNum: epoch, ChainIds: args[1:]}
			resp, err := queryClient.EpochChainsInfo(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
