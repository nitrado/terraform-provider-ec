package armada

import (
	"errors"

	"gitlab.com/nitrado/b2b/ec/armada/pkg/apiclient/clientset"
)

func resolveClientSet(m any) (clientset.Interface, error) {
	clientSet, ok := m.(clientset.Interface)
	if !ok {
		return nil, errors.New("invalid clientset")
	}
	return clientSet, nil
}
