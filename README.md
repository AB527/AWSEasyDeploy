# AWS Easy Deploy

![Go](https://img.shields.io/badge/go-1.21+-00ACD7?style=flat&logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![AWS](https://img.shields.io/badge/AWS-Elastic%20Beanstalk-FF9900?style=flat&logo=amazonaws&logoColor=white)
![CI](https://img.shields.io/badge/CI-GitHub%20Actions%20%7C%20GitLab%20CI-2DA44E?style=flat)

<br>
<img width="1294" height="170" alt="image" src="https://github.com/user-attachments/assets/9cbc9208-71f1-4c9b-9fe3-5c0ac4885fed" />
<br>
<br>

**AWS Elastic Beanstalk automation CLI. AWS App Runner simplicity, without the cost.**

AWS Easy Deploy is a Go-based CLI tool that gives AWS Elastic Beanstalk the power and simplicity of AWS App Runner, automating environment initialization, CI/CD pipeline generation, S3 packaging, and environment config injection with just a few commands.

Stop spending 45 to 60 minutes on manual deployment setup. AWS Easy Deploy gets your project live in **~10 minutes**, while cutting runtime costs by **40 to 60%** compared to equivalent App Runner deployments for persistent workloads.

## Why AWS Easy Deploy?

- **Eliminate managed runtime costs** of App Runner while keeping its automation experience
- **Self-generating CI/CD.** Automatically creates GitHub Actions and GitLab CI workflows for you
- **Full pipeline automation.** Zip, S3 upload, Elastic Beanstalk version creation, and deployment
- **Reduced setup time** from 45 to 60 mins down to ~10 mins
- **40 to 60% cost savings** vs. equivalent AWS App Runner persistent workloads
- **Push `.env` files directly** to Elastic Beanstalk environment config in one command

## Prerequisites

Before using AWS Easy Deploy, ensure you have the **AWS CLI installed and configured** with a profile that has the following IAM policies attached:

| IAM Policy | Purpose |
|---|---|
| `AdministratorAccess-AWSElasticBeanstalk` | Create and manage Elastic Beanstalk applications and environments |
| `AmazonEC2FullAccess` | Provision EC2 instances for Elastic Beanstalk environments |
| `AmazonS3FullAccess` | Create buckets and upload deployment artifacts |
| `IAMFullAccess` | Create and manage IAM roles required by Elastic Beanstalk |

Configure your AWS profile if you haven't already:

```bash
aws configure
```

## Installation

### Option 1. Install Script (Recommended)

```bash
curl -sL https://raw.githubusercontent.com/AB527/AWSEasyDeploy/main/install.sh | bash
```

This automatically downloads the correct binary for your OS and places it in your PATH.

### Option 2. Go Install

If you have Go installed:

```bash
go install github.com/AB527/AWSEasyDeploy@latest
```

### Verify Installation

```bash
easy-deploy --help
```

## Commands

### `init` - Initialize a New Deployment

The `init` command sets up your project for Elastic Beanstalk deployments. It prompts you for a few inputs and generates a `.easy-deploy` config file along with a ready-to-use CI/CD pipeline.

```bash
easy-deploy init
```

You will be prompted for the following:

| Input | Description |
|---|---|
| **Project Name** | Used as the Elastic Beanstalk application and environment name |
| **Framework** | Runtime for your app (e.g. `nodejs`, `python`, `go`) |
| **Branch Name** | Git branch that triggers deployments (e.g. `main`) |
| **Launch in VPC?** | Whether to deploy inside a VPC (`yes` / `no`) |

Once complete, a `.easy-deploy` config file and CI/CD pipeline are generated in your project. Commit and push these to your repository. From that point on, every push to the configured branch will automatically zip your application, upload it to S3, create a new Elastic Beanstalk application version, and deploy it to your environment.

#### Reinitializing an Existing Project

If a `.easy-deploy` file already exists (e.g. a teammate clones the repo), `init` detects it and prompts you:

```
$ easy-deploy init

? Configuration exists – reinitialize from it?:
  Yes
  No
```

- **Reinitialize from existing config.** Sets up the local environment using the committed config. Recommended for teammates cloning an existing project.
- **Overwrite with a fresh initialization.** Starts from scratch and replaces the existing config.

### `push-env` - Push Environment Variables to Elastic Beanstalk

The `push-env` command reads a local env file and automatically pushes all key-value pairs as configuration options to your Elastic Beanstalk environment, with no manual Console clicks required.

```bash
easy-deploy push-env <filename>

# Examples:
easy-deploy push-env .env
easy-deploy push-env .env.production
```

What it does:

- Parses the specified env file
- Injects all variables into the Elastic Beanstalk environment configuration
- Triggers an environment reload to apply the changes
- Skips blank lines and comments (lines starting with `#`)

## Quick Start

```bash
# 1. Install
curl -sL https://raw.githubusercontent.com/AB527/AWSEasyDeploy/main/install.sh | bash

# 2. Navigate to your project directory
cd my-project

# 3. Initialize, following the prompts for project name, framework, branch, and VPC
easy-deploy init

# 4. Commit the generated config and pipeline files
git add .easy-deploy
git commit -m "Add easy-deploy config"

# 5. Push to your branch. Your code will be deployed automatically once pushed to the platform
git push origin main

# 6. Optionally push your environment variables
easy-deploy push-env .env
```

## Contributing

We welcome contributions from the community! Here's how you can get started:

1. Fork the repository and create your branch from `main`.
2. Follow the Getting Started steps to set up your local environment.
3. Make your changes and commit them with a clear message.
4. Open a Pull Request and describe the changes you've made.

Have an idea for a new tool or an improvement? [Open an issue](https://github.com/AB527/AWSEasyDeploy/issues) to discuss it first.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/AB527/AWSEasyDeploy/blob/main/LICENSE) file for details. 
