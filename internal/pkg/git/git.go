package git

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"go.uber.org/zap"
	"time"
)

// Git 封装 go-git
// go-git 是 git 的子集，很多功能并不支持，查看不同：https://github.com/go-git/go-git/blob/master/COMPATIBILITY.md
// 比如 pull 只支持 fast-forward，不支持 stash，所以在封装的时候回做一些取舍：
// merge 与 解决冲突是十分复杂的操作，hollow 无法实现它们，所以在同步的时候采取以下策略：
//   - pull 时 如果传递 force=true，如果遇到 non-fast-forward，则会将远端文件全部下载下来，cp 到本地，相同文件保留最新的一个。尽量将降低影响。
//   - push：为了避免 push 的冲突，每次 push 都是 force 的，为了避免远端文件丢失，每次 push 之前都会 pull 一次。
//
// non-fast-forward: 当本地有提交，pull 都会报错 non-fast-forward。
type Git struct {
	log  *zap.SugaredLogger
	dir  billy.Filesystem
	r    *git.Repository
	auth *http.BasicAuth
}

type logWrite struct {
	log *zap.SugaredLogger
}

func (l *logWrite) Write(p []byte) (n int, err error) {
	l.log.Infof("%s", p)
	return len(p), nil
}

// NewGit return Git
// https://docs.github.com/cn/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
func NewGit(personalAccessTokens string, dir billy.Filesystem, log *zap.SugaredLogger) (g *Git, err error) {
	var auth *http.BasicAuth
	if personalAccessTokens != "" {
		auth = &http.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			// https://docs.github.com/cn/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
			Password: personalAccessTokens,
		}
	}
	g = &Git{
		log:  log,
		dir:  dir,
		r:    nil,
		auth: auth,
	}
	g.r, err = g.initRepo(dir)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Git) initRepo(dir billy.Filesystem) (*git.Repository, error) {
	dot, _ := dir.Chroot(".git")
	s := filesystem.NewStorage(dot, cache.NewObjectLRUDefault())

	r, err := git.Init(s, dir)
	if err != nil {
		if err == git.ErrRepositoryAlreadyExists {
		} else {
			return nil, fmt.Errorf("PlainInit error: %w", err)
		}
	} else {
		g.log.Infof("git init %v", dir.Root())
	}

	if r == nil {
		g.log.Infof("git open %v", dir.Root())
		r, err = git.Open(s, dir)
		if err != nil {
			return nil, fmt.Errorf("PlainOpen error: %w", err)
		}
	}

	return r, nil
}

func getFileLastCommitAt(r *git.Repository, filename string) (t time.Time) {
	l, _ := r.Log(&git.LogOptions{
		From:       plumbing.Hash{},
		Order:      0,
		FileName:   &filename,
		PathFilter: nil,
		All:        false,
		Since:      nil,
		Until:      nil,
	})
	l.ForEach(func(commit *object.Commit) error {
		t = commit.Author.When
		return fmt.Errorf("skip")
	})

	return
}

func (g *Git) Pull(remote string, branch string, force bool) error {
	remoteName := "origin-temp"
	err := g.r.DeleteRemote(remoteName)
	if err != nil {
		if err == git.ErrRemoteNotFound {
			err = nil
		} else {
			return err
		}
	}

	_, err = g.r.CreateRemote(&config.RemoteConfig{
		Name:  remoteName,
		URLs:  []string{remote},
		Fetch: nil,
	})
	if err != nil {
		return fmt.Errorf("CreateRemote error: %w", err)
	}
	defer func() {
		err = g.r.DeleteRemote(remoteName)
		if err != nil {
			if err == git.ErrRemoteNotFound {
				err = nil
			} else {
				g.log.Errorf("DeleteReomte error: %v", err)
			}
		}
	}()

	wt, err := g.r.Worktree()
	if err != nil {
		return fmt.Errorf("worktree error: %w", err)
	}

	referenceName := plumbing.NewBranchReferenceName(branch)
	err = wt.Pull(&git.PullOptions{
		RemoteName:        remoteName,
		ReferenceName:     referenceName,
		SingleBranch:      true,
		Depth:             0,
		Auth:              g.auth,
		RecurseSubmodules: 0,
		Progress: &logWrite{
			log: g.log,
		},
		Force:           force,
		InsecureSkipTLS: false,
		CABundle:        nil,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			g.log.Infof("AlreadyUpToDate")
			err = nil
		} else {
			return fmt.Errorf("pull error: %v", err)
		}
	}

	// reset to Head
	head, err := g.r.Head()
	if err != nil {
		return fmt.Errorf("get Head error: %v", err)
	}
	err = wt.Reset(&git.ResetOptions{
		Commit: head.Hash(),
		Mode:   git.HardReset,
	})
	if err != nil {
		return fmt.Errorf("reset to head '%v' error: %v", head.Hash(), err)
	}
	g.log.Infof("reset success")

	return nil
}
