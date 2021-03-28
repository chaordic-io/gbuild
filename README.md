# Graph Build
Simple meta-build tool that builds multi-language projects.

This is not intended to replace existing build tools like `yarn`, `npm`, `sbt`, `gradle`, but rather be an overlay, that allows you to compose actions from multiple build-tools in a safe and parallelised fashion for optimal build-speed. It will also allow you to run arbitrary scripts.

A build consists of a directed graph of `tasks`, where tasks can define dependencies on other `tasks`.
The tool calculates the optimal execution plan based on the dependency graph, and executes as many of the tasks in parallel as possible.
