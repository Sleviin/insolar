/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/pulsar"
	jww "github.com/spf13/jwalterweatherman"
)

type componentManager struct {
	components     core.Components
	interfaceNames []string
}

func (cm *componentManager) register(interfaceName string, component core.Component) {
	cm.interfaceNames = append(cm.interfaceNames, interfaceName)
	cm.components[interfaceName] = component
}

// linkAll - link dependency for all components
func (cm *componentManager) linkAll() {
	for _, name := range cm.interfaceNames {
		err := cm.components[name].Start(cm.components)
		if err != nil {
			log.Errorln("failed to start component ", name, " : ", err.Error())
		}
	}
}

// stopAll - reverse order stop all components
func (cm *componentManager) stopAll() {
	for i := len(cm.interfaceNames) - 1; i >= 0; i-- {
		name := cm.interfaceNames[i]
		log.Infoln("Stop component: ", name)
		err := cm.components[name].Stop()
		if err != nil {
			log.Errorln("failed to stop component ", name, " : ", err.Error())
		}
	}
}

func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	err := cfgHolder.Load()
	if err != nil {
		log.Warnln("Failed to load configuration from file: ", err.Error())
	}

	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("Failed to load configuration from env:", err.Error())
	}

	initLogger(cfgHolder.Configuration.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	nw, err := servicenetwork.NewServiceNetwork(cfgHolder.Configuration.Host, cfgHolder.Configuration.Node)
	if err != nil {
		log.Fatalln("Failed to start Network: ", err.Error())
	}

	l, err := ledger.NewLedger(cfgHolder.Configuration.Ledger)
	if err != nil {
		log.Fatalln("Failed to start Ledger: ", err.Error())
	}

	lr, err := logicrunner.NewLogicRunner(cfgHolder.Configuration.LogicRunner)
	if err != nil {
		log.Fatalln("Failed to start LogicRunner: ", err.Error())
	}

	eb, err := eventbus.New(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("Failed to start LogicRunner: ", err.Error())
	}

	b, err := bootstrap.NewBootstrapper(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("Failed to start Bootstrapper: ", err.Error())
	}

	ar, err := api.NewRunner(&cfgHolder.Configuration.APIRunner)
	if err != nil {
		log.Fatalln("Failed to start ApiRunner: ", err.Error())
	}

	met, err := metrics.NewMetrics(cfgHolder.Configuration.Metrics)
	if err != nil {
		log.Fatalln("Failed to start Metrics: ", err.Error())
	}

	cm := componentManager{components: make(core.Components), interfaceNames: make([]string, 0)}
	cm.register("core.Network", nw)
	cm.register("core.Ledger", l)
	cm.register("core.LogicRunner", lr)
	cm.register("core.EventBus", eb)
	cm.register("core.Bootstrapper", b)
	cm.register("core.ApiRunner", ar)
	cm.register("core.Metrics", met)
	cm.linkAll()
	lr.OnPulse(*pulsar.NewPulse(0, &pulsar.StandardEntropyGenerator{}))

	defer func() {
		cm.stopAll()
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debugln("caught sig: ", sig)

		cm.stopAll()
		os.Exit(0)
	}()

	fmt.Println("Running interactive mode:")
	repl(nw)
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(cfg.Level)
	if err != nil {
		log.Errorln(err.Error())
	}
}
