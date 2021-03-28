# Graph Build
## Problem definition
Many software projects today contain multiple languages, modules and deployment targets.
This can cause a few problems for teams, such as:

* Naive CI builds that build serially, instead of building the most efficient way in parallel.
* Hodge-podge home-built builds, where glue code are built into the build-tools of one language, to call build tools in other languages.
* Needing to "rebuild the world" on each build, because there is no concept of discovery of which module actually changed (for instance, rebuilding a frontend, when only backend code has changed).

A potential solution to this is adopting a tool like [Bazel](https://bazel.build), which is an excellent solution for monorepos and multi-language builds. But unfortunately, this leads to another set of problems, like a high-learning curve, having to throw away existing build-tools and potentially sacrifice IDE integration with these tools.

## The Solution/who this is for
_Graph Build_ is a _meta-build tool_, it is intended to work with your existing language specific build tools, such as `yarn`, `npm`, `sbt`, `gradle` or any others. It is intended as an overlay over these tools for efficient CI builds and local workflows.

What the tool does, is allow you to configure your multi-language, multi-module build as a directed graph of `targets`. Each `target` can depend on one or more other `targets`.
These targets are then bundled into an `execution plan`, which is the set of targets you want to build for a given purpose. `GraphBuild` will look at the `targets` in the `execution plan`, analyse their dependencies and then execute as much as possible in parallel for speed and efficiency, while honouring the dependency order of the targets that it has derived.

You can have multiple `execution plans`, for instance one for `CI`, one for `main releases`, one for `local development` etc. The intent is to allow local development to be as much a first class citizen as CI pipelines.

In the future, we intend to add target caching, so as to avoid rebuilding targets unecessarily, if they are unchanged.

## Example


```
targets:
- name: Frontend
  max_retries: 2
  path: front-end/
  run:
    |-
    yarn test
    yarn build
- name: Backend
  run:
    |-
    sbt test assembly
- name: PackageFE
  path: front-end/
  depends_on:
    - Frontend
  run: 
    |-
    docker build . -t my_image:latest
    docker push my_image:latest
- name: PackageBE
  depends_on:
    - Frontend
  run: 
    |-
    docker build . -t my_backend:latest
    docker push my_backend:latest
- name: Deploy
  depends_on:
    - Frontend
  run: 
    |-
    terraform apply
execution_plans:
  - name: CI
    targets:
    - Frontend
    - Backend
    - PackageFE
    - PackageBE
    - Deploy
  - name: Local
    targets:
    - Frontend
    - Backend
```

## TODO

* Isolated execution
* Caching of outputs and avoid re-running unchanged targets