# Contributing to p8s

First off, thank you for considering contributing to `p8s`! It's people like you that make open source such a great community. We welcome contributions of all kinds, from bug reports to feature requests, documentation improvements, and code submissions.

## Where to Start

*   **Discussions:** For general questions or to discuss ideas, please start a conversation on our [GitHub Discussions](https://github.com/p8s-dev/p8s/discussions) page.
*   **Bug Reports & Feature Requests:** To report a bug or request a new feature, please open an issue on our [GitHub Issues](https://github.com/p8s-dev/p8s/issues) page. Please use the provided templates to ensure you provide all the necessary information.

## Setting Up Your Development Environment

The `pace` CLI is written in Go. To contribute to the code, you will need to have the Go toolchain installed on your machine.

1.  **Fork the Repository:** Start by forking the main `p8s` repository to your own GitHub account.
2.  **Clone Your Fork:** Clone your forked repository to your local machine.
    ```bash
    git clone git@github.com:YOUR_USERNAME/p8s.git
    cd p8s
    ```
3.  **Build the CLI:** You can build the `pace` binary from the source code.
    ```bash
    cd pace-cli
    go build -o pace .
    ```
    You can now run the local build with `./pace`.

## Running Tests

`p8s` has a robust suite of tests to ensure the reliability of the code generation engine. Before submitting a pull request, please ensure that all tests pass.

The tests are located in the `pace-cli/` directory.

```bash
cd pace-cli
go test ./...
```

## Submitting a Pull Request

1.  Create a new branch for your feature or bug fix.
    ```bash
    git checkout -b my-awesome-feature
    ```
2.  Make your changes and commit them with a clear and descriptive commit message.
3.  Push your branch to your fork on GitHub.
    ```bash
    git push origin my-awesome-feature
    ```
4.  Open a pull request from your fork to the `main` branch of the `p8s-dev/p8s` repository.
5.  Please fill out the pull request template with as much detail as possible.

Thank you again for your contribution!
