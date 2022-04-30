# 🚌 Bus, monorepo/workspace tool (tutut)

## Why use a monorepo tool ?

> A monorepo is a version-controlled code repository that holds many projects. While these projects may be related, they are often logically independent and run by different teams.

[- what is a monorepo](https://semaphoreci.com/blog/what-is-monorepo)

[Everything you need to know about monorepos, and the tools to build them.](https://monorepo.tools/)

Monorepo tool is a software create for simplifying the management of monorepo, like [RushJS](https://github.com/microsoft/rushstack) or [Lerna](https://github.com/lerna/lerna).

### So, why use a monorepo ?

A monorepo is useful for any type of project, it allows you to group several packages and subprojects in a single repo. For example, for an application, you have a part for the website, another for the phone application and another for the desktop application. It's the same application, but for different devices. So instead of creating multiple repo's we'll put everything in one repo.

## Why bus?

Most of the monorepo tools work with a limited number of languages *(such as [Pants](https://github.com/pantsbuild/pants) with Python, Go, Java, Scala and Shell or [Bazel](https://github.com/bazelbuild/bazel/) with Java, C++ and Go or [RushJS](https://github.com/microsoft/rushstack) which only works with JavaScript)*. Bus want to simplify all and create a simple tool to simplify execution and configuration of projects. Bus ease without opinion the scripting for building, testing and executing your project.

## Quickstarting

### Install 

For the moment, the only way to install bus it's via git on unix with

```bash
git clone --depth 1 https://github.com/Erternys/bus.git

make all
make install
```

### Example

Example of config for a minimal project

```yml
# ./.bus.yaml
name: project name
version: 1.0.0
description: project description
repository: project repo
js_manager: npm / yarn / pnpm
packages:
    - path: path/to/sub-package
      name: sub-package
      extend: default/nodejs
```

```yml
# path/to/sub-package/.bus.yaml
name: sub-package
description: sub package description
version: 1.0.0
license: ISC
scripts: 
  script-name: command
  # or
  script-name: |
    command1
    command2
```

For generating these files, Bus has the command `bus init` and `bus init path/to/sub-package` and with `bus run script` you can run all scripts in subprojects.
