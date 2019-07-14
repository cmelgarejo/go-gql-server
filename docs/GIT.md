# Git Conventions

Git largely follows the `git flow` pattern, with
the main difference being merges happen to master
instead of dev and only fast forward merges are allowed.

Branch naming

- New features are
  - feature/\*

After you've completed changes for a new feature and would like review, create a
new "merge request" via github. If you'll only review, be used to append `WIP:`
to the beginning of the merge request title.

## Promoting code

### Use case for SPA site project (feature to develop)

When your `feature/fix/bug/whatever` branch is ready to be promoted, you should
rebase your current branch onto your target branch. After that, create a
Merge Request to `develop`.

Example using `feature/{project}/some-feature` as current branch, and `develop`
as target branch:

```git
git checkout feature/{project}/some-feature
git rebase origin/develop
git push origin feature/{project}/some-feature -f
```

Now you are able to do a Merge Request without issue.

### General use case (feature/develop to qa, qa to master)

Rebasing qa or master is not always a good choice (sometimes you get repeated
commits). So this is the proposed flow:

When trying to promote code from a feature or develop branch to `qa`, you can
merge target branch into your branch and then you can create a Merge Request:

Example using `develop` as current branch, and `qa` as target branch:

```git
git checkout develop/{project}`
git merge origin/qa --no-ff
git push origin develop/{project}
```

Now you're able to do a Merge Request without issue.

Example using `qa` as current branch and `master` as target branch:

```git
git checkout qa
git merge origin/master --no-ff
git pull
git push origin qa
```

Now you're able to do a Merge Request without issue.

## Conventions

- [Udacity Git Styleguide](https://udacity.github.io/git-styleguide/)

### Branch Naming

To keep it readable use the following branch convention when creating a new
branch:

```git
feature/{TASK-DESCRIPTION}
```

### Commit Messages

Idenfitify each commit message with a **TYPE** and short description.

### Extract from Udacity Git Styleguide

```text
type: subject

body

footer
```

The title consists of the type of the message and subject.

**The Type** _(Mandatory)_
The type is contained within the title and can be one of these types:

- feat: a new feature
- fix: a bug fix
- docs: changes to documentation
- style: formatting, missing semi colons, etc; no code change
- refactor: refactoring production code
- test: adding tests, refactoring test; no production code change
- chore: updating build tasks, package manager configs, etc; no production
  code change

### **The Subject** _(Mandatory)_

Subjects should be no greater than 50 characters, should begin with a capital
letter and do not end with a period.

Use an imperative tone to describe what a commit does, rather than what it did.
For example, use **change**; not changed or changes.

### **The Body** _(Optional)_

Not all commits are complex enough to warrant a body, therefore it is optional
and only used when a commit requires a bit of explanation and context. Use the
body to explain the what and why of a commit, not the how.

When writing a body, the blank line between the title and the body is required
and you should limit the length of each line to no more than 72 characters.

### **The Footer** _(Optional)_

The footer is optional and is used to reference issue tracker IDs.

### **Example**

```text
fix: Fix broken unit test

Fixed {some-feature} broken unit test: was failing because.....

Resolves: #XXX-YYY (Issue/JIRA/YourKanbanTool ticket)
```
