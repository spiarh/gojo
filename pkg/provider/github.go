package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/blang/semver/v4"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"

	"github.com/lcavajani/gojo/pkg/core"
	"github.com/lcavajani/gojo/pkg/util"
)

const (
	gitHubTokenEnvVar = "GITHUB_TOKEN"
)

type GitHub struct {
	client *github.Client
	log    zerolog.Logger

	owner      string            `yaml:"owner"`
	repository string            `yaml:"repository"`
	object     core.GitHubObject `yaml:"repository"`
}

type Versions struct {
	stable   []string
	unstable []string
}

func newGitHubClient() *github.Client {
	token := os.Getenv(gitHubTokenEnvVar)
	if token == "" {
		return github.NewClient(nil)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	return github.NewClient(oauth2.NewClient(ctx, ts))
}

func NewGitHub(owner, repo string, object core.GitHubObject) *GitHub {
	return &GitHub{
		client:     newGitHubClient(),
		log:        log.With().Str("provider", string(ProviderGitHub)).Logger(),
		owner:      owner,
		repository: repo,
		object:     object,
	}
}

func (g *GitHub) GetLatest() (string, error) {
	g.log.Info().Msg("get latest version")
	v, err := g.GetAll()
	if err != nil {
		return "", err
	}

	// TODO: add unstable & semver
	if len(v.stable) == 0 {
		return "", fmt.Errorf("no stable version found")
	}

	return v.stable[0], nil
}

func (g *GitHub) GetAll() (*Versions, error) {
	var v Versions

	switch g.object {
	case core.GitHubObjectRelease:
		releases, err := g.getRepoReleases()
		if err != nil {
			return nil, err
		}

		for _, r := range releases {
			if *r.Prerelease {
				v.unstable = append(v.unstable, *r.TagName)
				continue
			}
			v.stable = append(v.stable, *r.TagName)
		}
	case core.GitHubObjectTag:
		tags, err := g.getRepoTags()
		if err != nil {
			return nil, err
		}

		for _, t := range tags {
			tagName := util.SanitizeVersion(*t.Name)
			ver, err := semver.Parse(tagName)
			if err != nil {
				g.log.Warn().
					Str("tag", tagName).
					Err(err).
					Msg("parsing tag name failed")
				v.unstable = append(v.unstable, *t.Name)
				continue
			}
			if len(ver.Pre) > 0 {
				v.unstable = append(v.unstable, *t.Name)
				continue
			}
			v.stable = append(v.stable, *t.Name)
		}
	default:
		return nil, fmt.Errorf("github object type not recognized: %s", string(g.object))
	}

	g.log.Info().Int("len", len(v.stable)).
		Str("version", strings.Join(v.stable, ",")).
		Msg("stable versions")

	g.log.Info().Int("len", len(v.unstable)).
		Str("version", strings.Join(v.unstable, ",")).
		Msg("unstable versions")

	return &v, nil
}

func (g *GitHub) getRepoReleases() ([]*github.RepositoryRelease, error) {
	opt := &github.ListOptions{PerPage: 30}
	releases, _, err := g.client.Repositories.ListReleases(
		context.Background(),
		g.owner, g.repository,
		opt,
	)
	if err != nil {
		return nil, err
	}

	return releases, nil
}

func (g *GitHub) getRepoTags() ([]*github.RepositoryTag, error) {
	opt := &github.ListOptions{PerPage: 30}
	releases, _, err := g.client.Repositories.ListTags(
		context.Background(),
		g.owner, g.repository,
		opt,
	)
	if err != nil {
		return nil, err
	}
	return releases, nil
}
