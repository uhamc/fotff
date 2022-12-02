package vcs

type ManifestStep struct {
	LatestIssueURL string
	Manifest       Manifest
}

type Manifest struct {
}

func ManifestStepsExpand(from, to Manifest) []ManifestStep {
	return nil
}
