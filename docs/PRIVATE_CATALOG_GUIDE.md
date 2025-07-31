# Guide: Creating a Private `p8s` Catalog

This guide is for platform engineers (the "Producers") who want to create a private, company-specific catalog for `p8s`. A private catalog allows you to provide your application developers with a curated set of hardened, pre-approved infrastructure modules.

## The "Why": A Paved Road for Infrastructure

By creating a private catalog, you are building a "paved road" for your developers. This provides numerous benefits:
*   **Governance & Security:** Enforce security best practices and compliance standards in every piece of infrastructure.
*   **Velocity:** Free your application developers from the complexity of infrastructure management, allowing them to move faster.
*   **Consistency:** Ensure that all infrastructure is provisioned in a consistent, predictable, and supportable way.

## Step 1: Create Your Catalog Repository

The first step is to create a new, dedicated Git repository to house your catalog. This repository will contain your `catalog.yaml` and, optionally, any custom Terraform modules you maintain.

This repository should be accessible to the developers who will be using the `pace` CLI.

## Step 2: Build Your `catalog.yaml`

The `catalog.yaml` is the heart of your private catalog. It defines the resources available to your developers and maps their simple inputs to the underlying Terraform modules.

Here is a detailed breakdown of the schema:

```yaml
# catalog.yaml
resources:
  # This is a resource type, e.g., "postgres:v1"
  s3-bucket:
    # The source of the Terraform module
    source: "terraform-aws-modules/s3-bucket/aws"
    version: "3.15.1"

    # 'inputs' defines the API contract for this resource type.
    # It maps the developer's input from app.yaml to the module's variables.
    inputs:
      - from: "id"                   # Field name in app.yaml
        to: "bucket"                # Variable name in the Terraform module
        required: true
      - from: "versioning"
        to: "versioning.enabled"    # Dot notation for nested variables
        required: false
        default: false
```

*   **`source`**: The path to the Terraform module. This can be a reference to the public Terraform Registry, a GitHub repository, or a local path.
*   **`version`**: The specific version of the module to use. This is critical for ensuring stability.
*   **`inputs`**: This array defines the API for your resource.
    *   **`from`**: The key that a developer will use in their `app.yaml`.
    *   **`to`**: The variable name in the underlying Terraform module. You can use dot notation (e.g., `versioning.enabled`) to set nested variables.
    *   **`required`**: A boolean indicating whether the developer must provide this input.
    *   **`default`**: A default value to use if the input is not provided.

## Step 3: Using `pace init`

Once your `catalog.yaml` is committed to your repository, you can configure the `pace` CLI to use it with the `pace init` command.

This command tells `pace` where to find your private catalog. It will download the catalog and pre-fetch all the Terraform modules it references into a local cache (`~/.pace/cache`).

```bash
# Configure pace to use your private catalog
pace init --catalog-url git@github.com:my-company/my-pace-catalog.git
```

After running this command, any subsequent `pace generate` commands will use your private catalog instead of the default public one.

## Step 4: Advanced Features

`p8s` provides several advanced features for creating a robust and secure catalog.

### A. Value Mapping

For certain inputs, you may want to provide a simple, abstract option for developers that maps to a more complex value. For example, you can offer a `size` input that maps to specific instance classes.

This logic is handled within the `pace` tool itself and can be extended.

### B. Hardcoded Secure Defaults

You can enforce security by hardcoding critical variables directly in your templates or within the `pace` tool's logic. For example, you can ensure that all databases are created with `publicly_accessible = false`, regardless of user input. This is a powerful way to enforce security at the generation level.

By following this guide, you can create a powerful, private catalog that will accelerate your development teams while maintaining the highest standards of security and governance.
