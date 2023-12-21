# Golang Monorepo Project Template

[![CI/CD on main push](https://github.com/lorenzoranucci/golang-monorepo-project-template/actions/workflows/on-main-push-ci-cd.yml/badge.svg)](https://github.com/lorenzoranucci/golang-monorepo-project-template/actions/workflows/on-main-push-ci-cd.yml)

This repository serves as a template for creating a monorepo for Golang projects. It is designed to streamline the
development process, making it easier for teams to manage web applications and backend systems efficiently.

## Features

- **Monorepo Structure:** Organize your Golang project as a monorepo, allowing you to manage multiple related services
  and libraries in a single repository.

- **Go Workspace:** Leverage Go workspace to structure your project and manage dependencies effectively.

- **GitHub Actions:** Automate your workflows with GitHub Actions. CI/CD pipelines are set up to ensure a smooth
  development and deployment process.

- **Single Version Policy:** Adopt a single version policy for internal libraries to maintain consistency and simplify
  dependency management.

## Getting Started

To use this template for your own project, follow these steps:

1. Clone the newly created repository to your local machine.

2. Customize the project structure, adding your Golang applications and libraries.

3. Update the GitHub Actions workflows to fit the specific needs of your project.

4. Start building your Golang project with a robust and scalable monorepo structure!

### How to configure the `.netrc` file

The .netrc file is needed if you want to make this project private.
Go and in turn git will use the .netrc file to authenticate with GitHub to download the private packages in the `lib`.

1. Create a GitHub personal access token [here](https://github.com/settings/tokens).
2. Create a `.netrc` file in this project's root folder:
    - `cp .netrc.dist .netrc`  Note that the .netrc file is gitignored.
3. Replace `your-github-username` with your GitHub username.
4. Replace `your-github-personal-access-token` with the token you created in step 1.
