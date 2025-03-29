// Package autoload provides functionality related to PHP Composer autoloading
package autoload

// Autoload defines how the package should be autoloaded
type Autoload struct {
	PSR0        interface{} `json:"psr-0,omitempty"`
	PSR4        interface{} `json:"psr-4,omitempty"`
	Classmap    []string    `json:"classmap,omitempty"`
	Files       []string    `json:"files,omitempty"`
	ExcludeFrom []string    `json:"exclude-from-classmap,omitempty"`
}

// GetPSR4Map returns the PSR-4 map as a string to path mapping
func GetPSR4Map(a *Autoload) (map[string]string, bool) {
	psr4Map, ok := a.PSR4.(map[string]interface{})
	if !ok {
		return nil, false
	}

	result := make(map[string]string)
	for ns, path := range psr4Map {
		if pathStr, ok := path.(string); ok {
			result[ns] = pathStr
		}
	}
	return result, true
}

// SetPSR4 sets a PSR-4 namespace mapping
func SetPSR4(a *Autoload, namespace, path string) {
	psr4Map, ok := a.PSR4.(map[string]interface{})
	if !ok {
		psr4Map = make(map[string]interface{})
		a.PSR4 = psr4Map
	}
	psr4Map[namespace] = path
}

// RemovePSR4 removes a PSR-4 namespace mapping
func RemovePSR4(a *Autoload, namespace string) bool {
	psr4Map, ok := a.PSR4.(map[string]interface{})
	if !ok {
		return false
	}

	if _, exists := psr4Map[namespace]; !exists {
		return false
	}

	delete(psr4Map, namespace)
	return true
}
