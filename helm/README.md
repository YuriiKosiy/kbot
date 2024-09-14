# KBOT Helm Chart Release Guide

This README is available in multiple languages:
- [English](README.md) - curent page
- [Українська](README.ua.md)

This guide will walk you through the steps to package and release a new version of the KBOT Helm chart.

## Prerequisites

- Ensure that the version in the `values.yaml` file and chart templates is updated to the desired release version (in this guide, `1.0.12` is used as an example).
- Ensure that you have [Helm](https://helm.sh/docs/intro/install/) installed.
- Ensure that [GitHub CLI (gh)](https://cli.github.com/manual/installation) is installed.
- **GitHub Token:** You will need a GitHub personal access token to authenticate via GitHub CLI for creating a release.

---

## Steps for Creating a Helm Chart Release

### 1. Package the New Version of the Helm Chart

First, you need to package the updated version of the Helm chart. Replace `1.0.12` with the version number that you have set in your `values.yaml` file. For example, to release version `1.0.12`, run:

```bash
helm package ./helm --version 1.0.12
```

This will create a packaged Helm chart named `helm-1.0.12.tgz`.

> **Note:** Make sure to replace `1.0.12` with the actual version you're releasing, which should match the version in your `values.yaml` file.

---

### 2. Install GitHub CLI (if necessary)

If you haven't installed GitHub CLI yet, follow the official [installation instructions](https://cli.github.com/manual/installation).

---

### 3. Authenticate GitHub CLI

Before creating the release, you need to authenticate with GitHub. You can do this via:

#### Option 1: Logging in via GitHub CLI

Run the following command:

```bash
gh auth login
```

Follow the instructions to log in via a web browser or token.

#### Option 2: Using a GitHub Access Token

You can set up authentication via GitHub token using a personal access token, especially useful if you want to automate the process. To do this, export your GitHub token as an environment variable in your session:

1. Obtain a GitHub personal access token from [GitHub Settings -> Developer Settings -> Personal Access Tokens](https://github.com/settings/tokens). Ensure that you grant it the `repo` scope to access repositories and create releases.

2. Export the token in your shell session before running further commands:

```bash
export GITHUB_TOKEN="your_personal_access_token"
```

Now, when `gh` commands are executed, GitHub CLI will use this token for authentication.

---

### 4. Create a New GitHub Release

Now, create a new release for the KBOT Helm chart, making sure the version number (`1.0.12`) matches what you set in your `values.yaml` file:

```bash
gh release create v1.0.12 helm-1.0.12.tgz -t "kbot v1.0.12" -n "Updated release of kbot v1.0.12 with new improvements"
```

> **Note:** Replace `v1.0.12` and `helm-1.0.12.tgz` with the correct version you are releasing.

Here's a breakdown of the command:

- `v1.0.12` — The version you are releasing. **Make sure this is consistent with your `values.yaml`.**
- `helm-1.0.12.tgz` — The file path to the Helm chart package you created earlier.
- `-t "kbot v1.0.12"` — The title of the release.
- `-n "Updated release of kbot v1.0.12 with new improvements"` — The description of the release (you can modify this as needed).

This will create a new release on GitHub and upload the Helm chart package as an asset.

---

### 5. Verify the Release

Finally, confirm that the release was successfully created by visiting your repository's releases page on GitHub. Ensure the chart file is uploaded and everything looks correct.

---

### Additional Notes

- Replace all instances of `1.0.12` with the version number you're releasing, which should be reflected in your `values.yaml`.
- Always check that you've updated version-related information in the Helm chart (e.g., in `Chart.yaml` and `values.yaml`).
- If you automate the process, you can set `GITHUB_TOKEN` environment variable beforehand for authentication and continuous integration pipelines.
- You can follow this guide whenever releasing a new KBOT Helm chart.
