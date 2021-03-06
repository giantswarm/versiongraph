package versionbundle

import (
	"sort"
	"time"

	"github.com/giantswarm/microerror"
)

const releaseTimestampFormat = "2006-01-02T15:04:05.000000Z"

type ReleaseConfig struct {
	Active  bool
	Apps    []App
	Bundles []Bundle
	Date    time.Time
	Version string
}

type Release struct {
	apps       []App
	bundles    []Bundle
	components []Component
	timestamp  time.Time
	version    string
	active     bool
}

func NewRelease(config ReleaseConfig) (Release, error) {
	if len(config.Bundles) == 0 {
		return Release{}, microerror.Maskf(invalidConfigError, "%T.Bundles must not be empty", config)
	}

	r := Release{
		active:     config.Active,
		apps:       config.Apps,
		bundles:    config.Bundles,
		components: aggregateReleaseComponents(config.Bundles),
		timestamp:  config.Date,
		version:    config.Version,
	}

	return r, nil
}

func (r Release) Active() bool {
	return r.active
}

func (r Release) Apps() []App {
	return CopyApps(r.apps)
}

func (r Release) Bundles() []Bundle {
	return CopyBundles(r.bundles)
}

func (r Release) Components() []Component {
	return CopyComponents(r.components)
}

func (r Release) Timestamp() string {
	if r.timestamp.IsZero() {
		// This maintains existing behavior.
		return ""
	}

	return r.timestamp.Format(releaseTimestampFormat)
}

func (r Release) Version() string {
	return r.version
}

func aggregateReleaseComponents(bundles []Bundle) []Component {
	var components []Component

	for _, b := range bundles {
		bundleAsComponent := Component{
			Name:    b.Name,
			Version: b.Version,
		}
		components = append(components, bundleAsComponent)
		components = append(components, b.Components...)
	}

	sort.Sort(SortComponentsByName(components))

	return components
}

func GetNewestRelease(releases []Release) (Release, error) {
	if len(releases) == 0 {
		return Release{}, microerror.Maskf(executionFailedError, "releases must not be empty")
	}

	s := SortReleasesByVersion(releases)
	sort.Sort(s)

	return s[len(s)-1], nil
}
