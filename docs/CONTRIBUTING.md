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

## Commit Message Standard

To maintain a clean and readable Git history, this project uses the **Conventional Commits** specification for our Pull Request titles. Each PR title should be structured as follows:

```
<type>: <description>
```

**Common Types:**

*   **`feat`**: A new feature.
*   **`fix`**: A bug fix.
*   **`docs`**: Changes to documentation only.
*   **`style`**: Changes that do not affect the meaning of the code (white-space, formatting, etc).
*   **`refactor`**: A code change that neither fixes a bug nor adds a feature.
*   **`test`**: Adding missing tests or correcting existing tests.
*   **`chore`**: Changes to the build process or auxiliary tools and libraries.

**Example:**
```
feat: Add support for structured output generation
```
```
docs: Update CONTRIBUTING.md with commit message standard
```

A GitHub Action will automatically check your Pull Request title to ensure it meets this standard.

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

## Contributing to the Community Catalog

The `p8s` community catalog is a curated collection of modules that provide a "paved road" for developers to provision cloud infrastructure and application-level components. We welcome contributions of new modules to the catalog.

### The "One Module, One File" Model

To ensure the catalog is scalable, maintainable, and community-friendly, we use a "one module, one file" model, inspired by the design of Homebrew.

*   **One Module, One File:** Each module in the community catalog is defined in its own `module-name.yaml` file.
*   **Category-Based Directory Structure:** The catalogs are organized by category, providing an intuitive and discoverable experience for developers.

### How to Add a New Module

1.  **Find the Right Category:** Choose the appropriate category for your module from the existing directories in `pace-cli/catalogs`. If a suitable category does not exist, please open an issue to discuss creating a new one.
2.  **Create a New Module File:** Create a new `module-name.yaml` file in the appropriate category directory. The file name should be the same as the name of the resource in the `app.yaml` (e.g., `s3-bucket.yaml`).
3.  **Define the Module:** Define the module using our simple, declarative YAML schema.

    ```yaml
    # A short, descriptive title for the module
    #
    # A longer description of the module and what it does.
    #
    source: "<terraform-module-source>"
    version: "<terraform-module-version>"
    inputs:
      - from: "<app-yaml-input-name>"
        to: "<terraform-module-variable-name>"
        required: true
      - from: "<app-yaml-input-name>"
        to: "<terraform-module-variable-name>"
        required: false
        default: "<default-value>"
    ```
4.  **Open a Pull Request:** Open a pull request with your new module. Please include a clear and descriptive title and a detailed description of the module and why it should be added to the catalog.

Thank you again for your contribution!
