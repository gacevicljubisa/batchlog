package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethersphere/bee/v2/pkg/config"
	"github.com/ethersphere/bee/v2/pkg/util/abiutil"
	"github.com/gacevicljubisa/batchlog/pkg/ethclient"
	"github.com/gacevicljubisa/batchlog/pkg/filestore"
	"github.com/gacevicljubisa/batchlog/pkg/full"
	"github.com/spf13/cobra"
)

func (c *command) initFullCmd() (err error) {
	var (
		startBlock      uint64
		endBlock        uint64
		rpcEndpoint     string
		maxRequest      int
		blockRangeLimit uint32
	)

	cmd := &cobra.Command{
		Use:   "full",
		Short: "Fetch Swarm Postage Stamp contract event logs within a block range.",
		Long: `Fetches event logs for the Swarm Postage Stamp contract from a specified Ethereum RPC endpoint
within a given block range (--start to --end). It handles large ranges by querying in chunks (--block-range-limit)
and respects RPC rate limits (--max-request).

The retrieved logs are saved to 'export.ndjson' in NDJSON format.
The process can be interrupted at any time (Ctrl+C), and it will attempt to save already retrieved logs before exiting.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx := cmd.Context()

			ec, err := ethclient.NewClient(ctx, rpcEndpoint, ethclient.WithRateLimit(maxRequest))
			if err != nil {
				return fmt.Errorf("failed to connect to the Ethereum client: %w", err)
			}
			defer ec.Close()

			chainID, err := ec.ChainID(ctx)
			if err != nil {
				return fmt.Errorf("failed to get chainID: %w", err)
			}

			chainCfg, found := config.GetByChainID(chainID.Int64())
			if !found {
				return fmt.Errorf("chain config not found for chain ID %d", chainID.Int64())
			}

			postageStampContractABI := abiutil.MustParseABI(chainCfg.PostageStampABI)

			client := full.NewClient(ec, postageStampContractABI, blockRangeLimit)

			if startBlock == 0 {
				startBlock = chainCfg.PostageStampStartBlock
			}

			log.Printf("Retrieving logs from block %d to %d...\n", startBlock, endBlock)

			logChan, errorChan := client.GetLogs(ctx, &full.Request{
				Address:    chainCfg.PostageStampAddress,
				StartBlock: startBlock,
				EndBlock:   endBlock,
			})

			var wg sync.WaitGroup
			wg.Add(1)

			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			go func() {
				defer wg.Done()
				if err := filestore.SaveLogsAsync(ctx, logChan, "export.ndjson"); err != nil {
					if errors.Is(err, context.Canceled) {
						log.Fatalf("not all logs have been saved: %v", err)
					}
					log.Fatalf("failed to save logs: %v", err)
				}
			}()

			for {
				select {
				case err, ok := <-errorChan:
					if !ok {
						errorChan = nil
					} else {
						return fmt.Errorf("error retrieving logs: %w", err)
					}
				case <-ticker.C:
					log.Println("processing...")
				case <-ctx.Done():
					log.Println("shutting down...")
					wg.Wait()
					return ctx.Err()
				}

				if errorChan == nil {
					break
				}
			}

			wg.Wait()
			log.Println("all logs have been saved.")
			return nil
		},
	}

	cmd.Flags().Uint64VarP(&startBlock, "start", "", 31306381, "Start block number")
	cmd.Flags().Uint64VarP(&endBlock, "end", "", 39810670, "End block number")
	cmd.Flags().StringVarP(&rpcEndpoint, "endpoint", "e", "https://wandering-evocative-gas.xdai.quiknode.pro/0f2525676e3ba76259ab3b72243f7f60334b0000/", "ETH RPC endpoint")
	cmd.Flags().IntVarP(&maxRequest, "max-request", "m", 15, "Maximum number of requests per second")
	cmd.Flags().Uint32VarP(&blockRangeLimit, "block-range-limit", "b", 5, "Block range limit")

	c.root.AddCommand(cmd)

	return nil
}
