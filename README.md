# scrivgit

This is a crude wrapper around `git diff` for use with [Scrivener](https://www.literatureandlatte.com/scrivener/overview).

If you're using git to track changes to a Scrivener document you
can use `git log` to get a list of changes. Each change is identified
by a long hex number.

You can then use `scrivgit <that long hex number>` to show the changes
between that checkin and the current state. It will show all the
scrivener pages that have changed, and for each page show a diff
of the rtf format of the document. You have to be inside the scrivener
document directory, i.e. inside the Whatever.scriv directory.

It's a very crude wrapper around `git diff` and you can pass it
versions to compare in the same way. `scrivgit` will show changes
since the last checkin and so on.
