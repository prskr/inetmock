package main

import (
	"context"
	"fmt"
	"net/http"

	"code.gitea.io/sdk/gitea"
)

type CommitState string

func (s CommitState) String() string {
	return string(s)
}

func setCommitStatus(ctx context.Context, notification gitea.CreateStatusOption) error {
	if GiteaClient == nil {
		return nil
	}

	GiteaClient.SetContext(ctx)

	_, resp, err := GiteaClient.CreateStatus("inetmock", "inetmock", GitCommit, notification)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Close
	}()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed to update commit status - %d - %s", resp.StatusCode, resp.Status)
	}

	return nil
}
