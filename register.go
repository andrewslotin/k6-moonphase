package moonphase

import (
	"os"

	"go.k6.io/k6/js/modules"
)

const stormglassURL = "https://api.stormglass.io/v2"

func init() {
	modules.Register("k6/x/moonphase", New(stormglassURL, os.Getenv("STORMGLASS_API_KEY")))
}
