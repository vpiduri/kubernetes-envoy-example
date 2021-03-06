package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bakins/kubernetes-envoy-example/order"
	"github.com/spf13/cobra"
)

var address = ":8080"
var endpoint = "127.0.0.1:9090"

var rootCmd = &cobra.Command{
	Use:   "order",
	Short: "simple grpc order service",
	Run:   runServer,
}

func main() {
	f := rootCmd.Flags()
	f.StringVarP(&address, "address", "a", address, "listening address")
	f.StringVarP(&endpoint, "endpoint", "e", endpoint, "endpoint for contacting other services")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {

	s, err := order.New(
		order.SetAddress(address),
		order.SetEndpoint(endpoint),
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		s.Stop()
	}()

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
