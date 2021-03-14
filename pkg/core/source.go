package core

func NewAlpineSource(name, pkg, repo, versionId string) Source {
	return Source{
		Name: name,
		Provider: Provider{
			Alpine: &AlpineSource{
				Package:    pkg,
				Repository: repo,
				VersionId:  versionId,
			},
		},
	}
}

func NewGitHubSource(name, repo, owner string) Source {
	return Source{
		Name: name,
		Provider: Provider{
			GitHub: &GitHubSource{
				Owner:      owner,
				Repository: repo,
				Object:     GitHubObjectRelease,
			},
		},
	}

}
