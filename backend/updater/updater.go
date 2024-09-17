// Package updater implements a listener that listens the database for build
// updates and creates/updates check runs on GitHub.
package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/bee-ci/bee-ci-system/internal/auth"
	"github.com/google/go-github/v64/github"

	"github.com/bee-ci/bee-ci-system/data"

	"github.com/lib/pq"
)

const channelName = "builds_channel"

type Updater struct {
	logger        *slog.Logger
	httpClient    *http.Client
	dbListener    *pq.Listener
	channelName   string
	repoRepo      data.RepoRepo
	userRepo      data.UserRepo
	buildRepo     data.BuildRepo
	githubService *GithubService
}

func New(
	dbListener *pq.Listener,
	repoRepo data.RepoRepo,
	userRepo data.UserRepo,
	buildRepo data.BuildRepo,
	githubService *GithubService,
) *Updater {
	return &Updater{
		logger:        slog.Default(), // TODO: add some "subsystem name" to this logger
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		channelName:   channelName,
		dbListener:    dbListener,
		repoRepo:      repoRepo,
		userRepo:      userRepo,
		buildRepo:     buildRepo,
		githubService: githubService,
	}
}

// Start starts the updater. It will listen for updates from the database and
// create check runs on GitHub when the updates happen.
//
// To shutdown the updater, cancel the context.
func (u Updater) Start(ctx context.Context) error {
	err := u.dbListener.Listen(channelName)
	if err != nil {
		return fmt.Errorf("listen on channel %s: %w", channelName, err)
	}

	u.logger.Info("updater started, listens to db changes", slog.String("channel", channelName))

	for {
		select {
		case <-ctx.Done():
			u.logger.Debug("context cancelled, db listener will be closed")
			err = u.dbListener.Close()
			if err != nil {
				u.logger.Error("failed to close db listener", slog.Any("error", err))
				return err
			}
			return nil
		case msg := <-u.dbListener.Notify:
			updatedBuild := data.Build{}
			err := json.Unmarshal([]byte(msg.Extra), &updatedBuild)
			if err != nil {
				u.logger.Error("db listener got notification but it failed to unmarshal build", slog.Any("error", err))
				break
			}

			u.logger.Debug("db listener got notification",
				slog.Any("channel", msg.Channel),
				slog.Any("build", updatedBuild),
			)

			if updatedBuild.Status == "queued" && updatedBuild.CheckRunID == nil {
				// The build is new and hasn't been sent to GitHub yet. Create a new check run.
				checkRunID, err := u.createCheckRun(ctx, updatedBuild)
				if err != nil {
					u.logger.Error("failed to create check run", slog.Any("error", err))
					break
				}

				err = u.buildRepo.SetCheckRunID(ctx, updatedBuild.ID, checkRunID)
				if err != nil {
					u.logger.Error("failed to update check run ID in the database", slog.Any("error", err))
					break
				}
			} else if updatedBuild.Status == "queued" && updatedBuild.CheckRunID != nil {
				// Nothing to be done. Te check run on Gitub has been already created.
			} else {
				if updatedBuild.CheckRunID == nil {
					// This should never happen, but let's be extra safe.
					u.logger.Error("check run ID is nil", slog.Any("build", updatedBuild))
					break
				}
				// The build isn't new and has been sent to GitHub before. Update the check run.
				err = u.updateCheckRun(ctx, *updatedBuild.CheckRunID, updatedBuild)
				if err != nil {
					u.logger.Error("failed to update check run", slog.Any("error", err))
				}
			}
		}
	}
}

func (u Updater) createCheckRun(ctx context.Context, build data.Build) (checkRunID int64, err error) {
	repo, err := u.repoRepo.Get(ctx, build.RepoID)
	if err != nil {
		return 0, fmt.Errorf("ger repo: %w", err)
	}

	user, err := u.userRepo.Get(ctx, repo.UserID)
	if err != nil {
		return 0, fmt.Errorf("get user: %w", err)
	}

	installationAccessToken, err := u.githubService.GetInstallationAccessToken(ctx, build.InstallationID)
	if err != nil {
		return 0, fmt.Errorf("get installation access token: %w", err)
	}

	ghClient := github.NewClient(&http.Client{
		Transport: &auth.BearerTransport{Token: installationAccessToken},
	})

	createCheckRunOptions := github.CreateCheckRunOptions{
		// TODO: Get name from the BeeCI config file?
		Name:        build.CommitMsg + ", started at: " + fmt.Sprint(time.Now().Format(time.RFC822Z)),
		HeadSHA:     build.CommitSHA,
		DetailsURL:  github.String("https://bee-ci.vercel.app/dashboad/"), // TODO: Use actual URL of the backend
		ExternalID:  github.String(strconv.FormatInt(build.ID, 10)),
		Status:      github.String(build.Status),
		Conclusion:  nil,
		StartedAt:   &github.Timestamp{Time: build.CreatedAt},
		CompletedAt: nil,
		Output: &github.CheckRunOutput{
			Title:            github.String("This is check run title"),
			Summary:          github.String("This is check run summary"),
			Text:             github.String("This is check run text"),
			AnnotationsCount: nil,
			AnnotationsURL:   nil,
			Annotations:      nil,
			Images:           nil,
		},
		Actions: nil,
	}

	checkRun, _, err := ghClient.Checks.CreateCheckRun(ctx, user.Username, repo.Name, createCheckRunOptions)
	if err != nil {
		return 0, fmt.Errorf("create check run for repo %s/%s: %w", user.Username, repo.Name, err)
	}

	u.logger.Info("check run created",
		slog.String("html_url", *checkRun.HTMLURL),
		slog.Any("build", build),
	)

	return *checkRun.ID, nil
}

func (u Updater) updateCheckRun(ctx context.Context, checkRunID int64, build data.Build) error {
	repo, err := u.repoRepo.Get(ctx, build.RepoID)
	if err != nil {
		return fmt.Errorf("ger repo: %w", err)
	}

	user, err := u.userRepo.Get(ctx, repo.UserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	installationAccessToken, err := u.githubService.GetInstallationAccessToken(ctx, build.InstallationID)
	if err != nil {
		return fmt.Errorf("get installation access token: %w", err)
	}

	ghClient := github.NewClient(&http.Client{
		Transport: &auth.BearerTransport{Token: installationAccessToken},
	})

	// TODO: Do I need to set these options again, or if I set them to null they will be removed?
	checkRunUpdateOptions := github.UpdateCheckRunOptions{
		Name:        build.CommitMsg + ", started at: " + fmt.Sprint(time.Now().Format(time.RFC822Z)),
		DetailsURL:  github.String("https://bee-ci.vercel.app/dashboad/"), // TODO: Use actual URL of the backend
		ExternalID:  nil,
		Status:      nil,
		Conclusion:  nil,
		CompletedAt: nil,
		Output:      nil,
		Actions:     nil,
	}
	if build.Conclusion != nil {
		checkRunUpdateOptions.Conclusion = build.Conclusion
		checkRunUpdateOptions.CompletedAt = &github.Timestamp{Time: build.UpdatedAt}
	}

	checkRun, _, err := ghClient.Checks.UpdateCheckRun(ctx, user.Username, repo.Name, checkRunID, checkRunUpdateOptions)
	if err != nil {
		return fmt.Errorf("update check run for repo %s/%s: %w", user.Username, repo.Name, err)
	}

	u.logger.Info("check run updated",
		slog.String("html_url", *checkRun.HTMLURL),
		slog.Any("build", build),
	)

	return nil
}
