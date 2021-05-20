package utils

import (
	"context"
	"io"

	"github.com/go-git/go-git/v5"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/isnlan/coral/pkg/errors"
)

func GitClone(ctx context.Context, username, password, remotePath, localPath string, writer io.Writer) error {
	if writer == nil {
		return errors.New("writer is null")
	}

	auth := &gitHttp.BasicAuth{Username: username, Password: password}
	_, err := git.PlainCloneContext(ctx, localPath, false, &git.CloneOptions{
		URL:      remotePath,
		Auth:     auth,
		Progress: writer,
	})

	return err
}
