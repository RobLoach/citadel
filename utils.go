package citadel

import (
	"fmt"
	"strings"

	"github.com/samalba/dockerclient"
)

// ValidateContainer ensures that the required fields are set on the container
func ValidateContainer(c *Container) error {
	switch {
	case c.Cpus == 0:
		return fmt.Errorf("container cannot have cpus equal to 0")
	case c.Memory == 0:
		return fmt.Errorf("container cannot have memory equal to 0")
	case c.Image == "":
		return fmt.Errorf("container must have an image")
	case c.Name == "":
		return fmt.Errorf("container must have a name")
	case c.Type == "":
		return fmt.Errorf("container must have a type")
	}

	return nil
}
func AsCitadelContainer(container *dockerclient.Container, engine *Docker) (*Container, error) {
	info, err := engine.client.InspectContainer(container.Id)
	if err != nil {
		return nil, err
	}
	cType := ""
	labels := []string{}
	env := make(map[string]string)
	for _, e := range info.Config.Env {
		vals := strings.Split(e, "=")
		k, v := vals[0], vals[1]
		switch k {
		case "_citadel_type":
			cType = v
		case "_citadel_labels":
			labels = strings.Split(v, ",")
		default:
			env[k] = v
		}
	}
	return &Container{
		Name:        info.Name,
		Image:       container.Image,
		Cpus:        float64(info.Config.CpuShares),
		Memory:      float64(info.Config.Memory),
		Environment: env,
		Hostname:    info.Config.Hostname,
		Domainname:  info.Config.Domainname,
		Type:        cType,
		Labels:      labels,
	}, nil
}
