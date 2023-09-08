package cmd

import (
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/ngodn/codereport-cli/extensions"
	"github.com/ngodn/codereport-cli/extensions/options"
	"github.com/ngodn/codereport-cli/pkg/locator"
	"go.riyazali.net/sqlite"

	// bring in sqlite 🙌
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/ngodn/codereport-cli/pkg/sqlite"
)

func registerExt() {
	multiLocOpt := &locator.MultiLocatorOptions{
		CloneDir:        cloneDir,
		InsecureSkipTLS: gitSSLNoVerify != "",
	}
	if githubToken != "" {
		multiLocOpt.HTTPAuth = &http.BasicAuth{Username: githubToken}
	}

	var skipMailmapCtx string
	if skipMailmap {
		skipMailmapCtx = "true"
	}

	sqlite.Register(
		extensions.RegisterFn(
			options.WithExtraFunctions(),
			options.WithRepoLocator(locator.CachedLocator(locator.LoggingLocator(
				&logger,
				locator.MultiLocator(multiLocOpt),
			))),
			options.WithContextValue("defaultRepoPath", repo),
			options.WithContextValue("skipMailmap", skipMailmapCtx),
			options.WithGitHub(),
			options.WithContextValue("githubToken", githubToken),
			options.WithContextValue("githubPerPage", os.Getenv("GITHUB_PER_PAGE")),
			options.WithContextValue("githubRateLimit", os.Getenv("GITHUB_RATE_LIMIT")),
			options.WithSourcegraph(),
			options.WithContextValue("sourcegraphToken", sourcegraphToken),
			options.WithNPM(),
			options.WithLogger(&logger),
		),
	)
}
