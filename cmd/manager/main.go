package main

import (
	"fmt"
	"os"

	"github.com/mongodb/mongodb-kubernetes-operator/pkg/apis"
	"github.com/mongodb/mongodb-kubernetes-operator/pkg/controller"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// Change below variables to serve metrics on different host or port.
var (
	metricsHost               = "0.0.0.0"
	metricsPort         int32 = 8383
	operatorMetricsPort int32 = 8686
)

func configureLogger() (*zap.Logger, error) {
	// TODO: configure non development logger
	logger, err := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	return logger, err
}

func hasRequiredVariables(logger *zap.Logger, envVariables ...string) bool {
	allPresent := true
	for _, envVariable := range envVariables {
		if _, envSpecified := os.LookupEnv(envVariable); !envSpecified {
			logger.Error(fmt.Sprintf("required environment variable %s not found", envVariable))
			allPresent = false
		}
	}
	return allPresent
}

func main() {
	log, err := configureLogger()
	if err != nil {
		os.Exit(1)
	}

	if !hasRequiredVariables(log, "AGENT_IMAGE") {
		os.Exit(1)
	}

	// get watch namespace from environment variable
	namespace, nsSpecified := os.LookupEnv("WATCH_NAMESPACE")
	if !nsSpecified {
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("Watching namespace: %s", namespace))

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace: namespace,
	})

	if err != nil {
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		os.Exit(1)
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		os.Exit(1)
	}
}
