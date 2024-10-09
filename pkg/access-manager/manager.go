package access_manager

import (
	"context"
	"fmt"
	"github.com/t34-dev/go-utils/pkg/etcd"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	cliv3 "go.etcd.io/etcd/client/v3"
)

const (
	etcdTimeout  = 5 * time.Second
	configPrefix = "/casbin/"
	modelKey     = configPrefix + "model"
	policyKey    = configPrefix + "policy"
)

type AccessManager interface {
	CheckAccess(role, resourceName, action string) (bool, error)
	UpdateConfigsFromEtcd(context.Context) error
	WatchConfig(context.Context, func(error, string, []byte)) error
	StopWatchConfig() error
	UpdateConfigsFromFiles(model, policy string) error
	UpdateEtcdStore(context.Context) error
}

type accessManagerData struct {
	modelStr  string
	policyStr string
}
type accessManager struct {
	enforcer   *casbin.Enforcer
	etcdClient etcd.Client
	mu         sync.RWMutex
	watchCh    cliv3.WatchChan
	stopCh     chan struct{}
	data       *accessManagerData
}

func NewAccessManager(etcdClient etcd.Client) AccessManager {
	return &accessManager{
		mu:         sync.RWMutex{},
		etcdClient: etcdClient,
		stopCh:     make(chan struct{}),
	}
}

func (a *accessManager) CheckAccess(role, resourceName, action string) (bool, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.enforcer == nil {
		return false, fmt.Errorf("enforcer not initialized")
	}

	return a.enforcer.Enforce(role, resourceName, action)
}

func (a *accessManager) UpdateConfigsFromEtcd(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	modelStr, err := a.getEtcdValue(ctx, modelKey)
	if err != nil {
		return fmt.Errorf("failed to get `model` from etcdClient: %w", err)
	}

	policyStr, err := a.getEtcdValue(ctx, policyKey)
	if err != nil {
		return fmt.Errorf("failed to get `policy` from etcdClient: %w", err)
	}

	enforcer, err := createEnforcer(modelStr, policyStr)
	if err != nil {
		return fmt.Errorf("failed to create enforcer: %w", err)
	}
	a.enforcer = enforcer
	a.data = &accessManagerData{modelStr: modelStr, policyStr: policyStr}

	return nil
}

func (a *accessManager) WatchConfig(ctx context.Context, callback func(error, string, []byte)) error {
	// Use a common prefix to track both keys
	a.watchCh = a.etcdClient.EtcdClient().Watch(context.Background(), configPrefix, cliv3.WithPrefix())

	go func() {
		for {
			select {
			case watchResp := <-a.watchCh:
				for _, event := range watchResp.Events {
					key := string(event.Kv.Key)
					newValue := watchResp.Events[0].Kv.Value
					if a.data == nil {
						callback(fmt.Errorf("no data to initialize key %s", key), key, newValue)
						continue
					}
					if (strings.HasPrefix(key, modelKey) && a.data.modelStr != string(newValue)) ||
						(strings.HasPrefix(key, policyKey) && a.data.policyStr != string(newValue)) {
						err := a.UpdateConfigsFromEtcd(ctx)
						callback(err, key, newValue)
					}
				}
			case <-a.stopCh:
				return
			}
		}
	}()

	return nil
}

func (a *accessManager) StopWatchConfig() error {
	if a.watchCh != nil {
		close(a.stopCh)
	}
	return nil
}

func (a *accessManager) UpdateConfigsFromFiles(modelFile, policyFile string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Load model from file
	modelBytes, err := os.ReadFile(modelFile)
	if err != nil {
		return fmt.Errorf("failed to read model file: %w", err)
	}
	// Load policy from file
	policyBytes, err := os.ReadFile(policyFile)
	if err != nil {
		return fmt.Errorf("failed to read policy file: %w", err)
	}

	enforcer, err := createEnforcer(string(modelBytes), string(policyBytes))
	if err != nil {
		return fmt.Errorf("failed to create enforcer: %w", err)
	}
	a.enforcer = enforcer
	a.data = &accessManagerData{modelStr: string(modelBytes), policyStr: string(policyBytes)}

	return nil
}

func (a *accessManager) UpdateEtcdStore(ctx context.Context) error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.enforcer == nil || a.data == nil {
		return fmt.Errorf("model or enforcer not initialized")
	}

	// model
	if modelVal, err := a.getEtcdValue(ctx, modelKey); err != nil {
		return fmt.Errorf("failed to get `model` from etcdClient: %w", err)
	} else if modelVal != a.data.modelStr {
		if err := a.setETCDValue(ctx, modelKey, a.data.modelStr); err != nil {
			return fmt.Errorf("failed to update `model` in etcdClient: %w", err)
		}
	}
	// policy
	if policyVal, err := a.getEtcdValue(ctx, policyKey); err != nil {
		return fmt.Errorf("failed to get `policy` from etcdClient: %w", err)
	} else if policyVal != a.data.policyStr {
		if err := a.setETCDValue(ctx, policyKey, a.data.policyStr); err != nil {
			return fmt.Errorf("failed to update `policy` in etcdClient: %w", err)
		}
	}

	return nil
}

func (a *accessManager) getEtcdValue(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	resp, err := a.etcdClient.EtcdClient().Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return string(resp.Kvs[0].Value), nil
}

func (a *accessManager) setETCDValue(ctx context.Context, key, value string) error {
	ctx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	_, err := a.etcdClient.EtcdClient().Put(ctx, key, value)
	return err
}
