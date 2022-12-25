# gobase
A base golang setup

Set up in a repo:
```
  git remote add go-base git@github.com:leinadao/go-base.git
  git fetch go-base
  git rebase go-base/main # Resolve any conflicts.
  git push # possibly need --force
```

Update latest base changes in a repo:
```
  git fetch go-base
  git rebase go-base/main # Resolve any conflicts.
  git push # possibly need --force
```
