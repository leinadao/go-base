[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

# Go Base
A base golang setup

## Set up in another repo
```
  git remote add go-base git@github.com:leinadao/go-base.git
  git fetch go-base
  git rebase go-base/main # Resolve any conflicts.
  git push # possibly need --force
```
Update latest base changes in a repo:
```
  git fetch go-base
  git merge go-base/main --no-commit
  # Resolve any conflicts and stage files.
  # Check diff for new standards etc to match.
  git merge --continue
  git push
```
