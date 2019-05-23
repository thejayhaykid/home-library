package actions

import (
	"github.com/gobuffalo/envy"
)

// ENV determines if the environment is dev, test, or prod
var ENV = envy.Get("GO_ENV", "development")
