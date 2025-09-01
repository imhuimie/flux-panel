package socket

import (
	"errors"
	"github.com/go-gost/x/config"
	parser "github.com/go-gost/x/config/parsing/limiter"
	"github.com/go-gost/x/registry"
	"strings"
)

func createLimiter(req createLimiterRequest) error {
	name := strings.TrimSpace(req.Data.Name)
	if name == "" {
		return errors.New("limiter name is required")
	}
	req.Data.Name = name

	if registry.TrafficLimiterRegistry().IsRegistered(name) {
		return errors.New("limiter " + name + " already exists")
	}

	v := parser.ParseTrafficLimiter(&req.Data)

	if err := registry.TrafficLimiterRegistry().Register(name, v); err != nil {
		return errors.New("limiter " + name + " already exists")
	}

	config.OnUpdate(func(c *config.Config) error {
		c.Limiters = append(c.Limiters, &req.Data)
		return nil
	})

	return nil
}

func updateLimiter(req updateLimiterRequest) error {

	name := strings.TrimSpace(req.Limiter)

	if !registry.TrafficLimiterRegistry().IsRegistered(name) {
		return errors.New("limiter " + name + " not found")
	}

	req.Data.Name = name

	v := parser.ParseTrafficLimiter(&req.Data)

	registry.TrafficLimiterRegistry().Unregister(name)

	if err := registry.TrafficLimiterRegistry().Register(name, v); err != nil {
		return errors.New("limiter " + name + " already exists")
	}

	config.OnUpdate(func(c *config.Config) error {
		for i := range c.Limiters {
			if c.Limiters[i].Name == name {
				c.Limiters[i] = &req.Data
				break
			}
		}
		return nil
	})

	return nil
}

func deleteLimiter(req deleteLimiterRequest) error {

	name := strings.TrimSpace(req.Limiter)

	if !registry.TrafficLimiterRegistry().IsRegistered(name) {
		return errors.New("limiter " + name + " not found")
	}
	registry.TrafficLimiterRegistry().Unregister(name)

	config.OnUpdate(func(c *config.Config) error {
		limiteres := c.Limiters
		c.Limiters = nil
		for _, s := range limiteres {
			if s.Name == name {
				continue
			}
			c.Limiters = append(c.Limiters, s)
		}
		return nil
	})

	return nil
}

type createLimiterRequest struct {
	Data config.LimiterConfig `json:"data"`
}

type updateLimiterRequest struct {
	Limiter string               `json:"limiter"`
	Data    config.LimiterConfig `json:"data"`
}

type deleteLimiterRequest struct {
	Limiter string `json:"limiter"`
}
