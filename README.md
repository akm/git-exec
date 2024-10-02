# git-exec

## install

```
go install github.com/akm/git-exec@latest
```

You might have to run `asdf reshim` .

## Usage

```
git-exec (your command with argument)
```

git-exec executes given command and `git add .` and `git commit` automatically.

### Option

You can change message prefix by environment variable `GIT_EXEC_COMMIT_PREFIX` .
