package servicediscovery

import (
	"context"
	"strings"

	"log"

	gokit_log "github.com/go-kit/kit/log"
	"github.com/prometheus/common/model"
	prometheus_discovery "github.com/prometheus/prometheus/discovery"
	prometheus_discovery_config "github.com/prometheus/prometheus/discovery/config"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/pkg/relabel"
	"github.com/soundcloud/periskop/config"
)

type ResolvedAddresses struct {
	Addresses []string
}

func EmptyResolvedAddresses() ResolvedAddresses {
	return ResolvedAddresses{
		Addresses: make([]string, 0),
	}
}

type Resolver struct {
	sdConfig       map[string]prometheus_discovery_config.ServiceDiscoveryConfig
	relabelConfigs []*relabel.Config
}

func NewResolver(service config.Service) Resolver {
	sdConfig := map[string]prometheus_discovery_config.ServiceDiscoveryConfig{
		service.Name: service.ServiceDiscovery,
	}

	return Resolver{
		sdConfig:       sdConfig,
		relabelConfigs: service.RelabelConfigs,
	}
}

func (r Resolver) Resolve() <-chan ResolvedAddresses {
	ctx := context.Background()
	out := make(chan ResolvedAddresses)
	manager := prometheus_discovery.NewManager(ctx, gokit_log.NewNopLogger())

	err := manager.ApplyConfig(r.sdConfig)
	if err != nil {
		log.Fatal("Could not apply SD configuration")
	}

	go func() {
		err = manager.Run()
	}()

	if err != nil {
		log.Fatal("Could not initialize SD manager")
	}

	go func() {
		for {
			var addresses []string
			groups := <-manager.SyncCh()

			for _, groupArr := range groups {
				for i := 0; i < len(groupArr); i++ {
					group := groupArr[i]
					for _, target := range group.Targets {
						discoveredLabels := group.Labels.Merge(target)
						var labelMap = make(map[string]string)
						for k, v := range discoveredLabels.Clone() {
							labelMap[string(k)] = string(v)
						}

						processedLabels := relabel.Process(labels.FromMap(labelMap), r.relabelConfigs...)

						var labels = make(model.LabelSet)
						for k, v := range processedLabels.Map() {
							labels[model.LabelName(k)] = model.LabelValue(v)
						}

						// Drop empty targets (drop in relabeling).
						if processedLabels == nil {
							continue
						}

						for k := range labels {
							if strings.HasPrefix(string(k), "__") {
								delete(labels, k)
							}
						}

						addresses = append(addresses, string(target["__address__"]))
						log.Println(addresses)
					}
				}
			}
			out <- ResolvedAddresses{
				Addresses: addresses,
			}
		}
	}()

	return out
}
