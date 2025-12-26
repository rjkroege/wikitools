You are in the root of a golang project. There are multiple places
where the go code walks over the a tree of files. There is
now a channel based concurrency-safe wrapper around the
state that is updated corresponding to each file. Analyze the
code and append a Markdown task list to this file of a series
of steps where each such step will become a separate CL
such that the complete set of steps will invoke the per-file
actions in separate go routines. Do not yet write code.
Document each proposed step with
a bulleted sub-list of relative path references into the source
files where you expect changes to occur.

Pause at this point to permit me to review and possibly update your
descriptions.

Once I indicate to you to continue, complete the
previously generated markdown list of
tasks. Do each task not yet marked as done one at a time. Once `go test
./...` passes for each task, git commit the change with a description
of the change to the current branch. Stop in the task list if you
cannot get the tests to pass. Continue through the tasks until
are completed. Update this file (prompt.md) to record the completion
of each task.