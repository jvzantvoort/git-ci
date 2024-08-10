[![forthebadge](https://forthebadge.com/images/badges/made-with-out-pants.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/it-works-why.svg)](https://forthebadge.com)
[![forthebadge](https://forthebadge.com/images/badges/designed-in-etch-a-sketch.svg)](https://forthebadge.com)

# git-ci

Lazy commit wrapper for branches created from JIRA tickets.

Given a branch from Jira for example
"feature/ABC-1234-fix-teh-interwebs" the prefix for a commit should
be "ABC-1234" for it to propperly be mentioned in the tickets. This
command wraps the default message with that ticket id.

```
# git rev-parse --abbrev-ref HEAD
feature/ABC-1234-fix-teh-interwebs

# git ci -m "reset router"

# git log --oneline -n 1
15ad934 ABC-1234 reset router
```
## Install



```
ln -s git-ci git-bs
ln -s git-ci git-ci
ln -s git-ci git-feat
ln -s git-ci git-fix
ln -s git-ci git-chore
ln -s git-ci git-docs
ln -s git-ci git-style
ln -s git-ci git-refactor
ln -s git-ci git-perf
ln -s git-ci git-test
ln -s git-ci git-new
ln -s git-ci git-build
ln -s git-ci git-revert
ln -s git-ci git-config
ln -s git-ci git-hotfix
ln -s git-ci git-release
```

*Warning*: git bs add random messages
