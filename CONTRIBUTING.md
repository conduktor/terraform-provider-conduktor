<!-- omit in toc -->
# Contributing to terraform-provider-conduktor

First off, thanks for taking the time to contribute! ❤️

All types of contributions are encouraged and valued. See the [Table of Contents](#table-of-contents) for different ways to help and details about how this project handles them. Please make sure to read the relevant section before making your contribution. It will make it a lot easier for us maintainers and smooth out the experience for all involved. The community looks forward to your contributions. 🎉

> And if you like the project, but just don't have time to contribute, that's fine. There are other easy ways to support the project and show your appreciation, which we would also be very happy about:
> - Star the project
> - Tweet about it
> - Refer this project in your project's readme
> - Mention the project at local meetups and tell your friends/colleagues

<!-- omit in toc -->
## Table of Contents

- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Your First Code Contribution](#your-first-code-contribution)
- [Improving The Documentation](#improving-the-documentation)
- [Styleguides](#styleguides)
- [Commit Messages](#commit-messages)
- [Join The Project Team](#join-the-project-team)



## I Have a Question

> If you want to ask a question, we assume that you have read the available [Documentation](https://docs.conduktor.io/).

Before you ask a question, it is best to search for existing [Issues](https://github.com/conduktor/terraform-provider-conduktor/issues) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue. It is also advisable to search the internet for answers first.

If you then still feel the need to ask a question and need clarification, we recommend the following:

- Open an [Issue](https://github.com/conduktor/terraform-provider-conduktor/issues/new).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions (go, docker, etc), depending on what seems relevant.

We will then take care of the issue as soon as possible.

## I Want To Contribute

### Reporting Bugs

<!-- omit in toc -->
#### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you up for more information. Therefore, we ask you to investigate carefully, collect information and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

- Make sure that you are using the latest version.
- Determine if your bug is really a bug and not an error on your side e.g. using incompatible environment components/versions (Make sure that you have read the [documentation](https://docs.conduktor.io/). If you are looking for support, you might want to check [this section](#i-have-a-question)).
- To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in the [bug tracker](https://github.com/conduktor/terraform-provider-conduktor/issues?q=label%3Abug).
- Collect information about the bug:
  - Version of the terraform-provider-conduktor you are using
  - Version of Conduktor Console or Gateway you are targeting
  - Stack trace (Traceback)
  - OS, Platform and Version (Windows, Linux, macOS, x86, ARM)
  - Version of the interpreter, compiler, SDK, runtime environment, package manager, depending on what seems relevant.
  - Possibly your input and the output
  - Can you reliably reproduce the issue? And can you also reproduce it with older versions?

<!-- omit in toc -->
#### How Do I Submit a Good Bug Report?

> You must never report security related issues, vulnerabilities or bugs including sensitive information to the issue tracker, or elsewhere in public. Instead sensitive bugs must be sent by email to [security@conduktor.io]().

We use GitHub issues to track bugs and errors. If you run into an issue with the project:

- Open an [Issue](https://github.com/conduktor/terraform-provider-conduktor/issues/new). (Since we can't be sure at this point whether it is a bug or not, we ask you not to talk about a bug yet and not to label the issue.)
- Explain the behavior you would expect and the actual behavior.
- Please provide as much context as possible and describe the *reproduction steps* that someone else can follow to recreate the issue on their own. This usually includes your code. For good bug reports you should isolate the problem and create a reduced test case.
- Provide the information you collected in the previous section.

You can also contact us at [support@conduktor.io]() if you need help with the issue.


### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for Conduktor Terraform Provider, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

<!-- omit in toc -->
#### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Read the [documentation](https://docs.conduktor.io/) carefully and find out if the functionality is already covered, maybe by an individual configuration.
- Perform a [search](https://github.com/conduktor/terraform-provider-conduktor/issues) to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.
- Search for similar suggestions in the [Conduktor Roadmap](https://product.conduktor.help/) and Submit an idea if it is not already there.
- Find out whether your idea fits with the scope and aims of the project. It's up to you to make a strong case to convince the project's developers of the merits of this feature. Keep in mind that we want features that will be useful to the majority of our users and not just a small subset.

<!-- omit in toc -->
#### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/conduktor/terraform-provider-conduktor/issues) or in [Conduktor Roadmap](https://product.conduktor.help/).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Provide a **step-by-step description of the suggested enhancement** in as many details as possible.
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why. At this point you can also tell which alternatives do not work for you.
- **Explain why this enhancement would be useful** to most Conduktor Terraform Provider users. You may also want to point out the other projects that solved it better and which could serve as inspiration.

### Your First Code Contribution

To run the project locally, you will need Go 1.22+ installed and configured on your machine.
Here some useful commands to get you started:

Install pre-commit git hooks:
```shell
make install-githooks
```

Run code formatting and linting:
```shell
make go-fmt
make go-lint
```

Build the provider locally:
```shell
make build
```

Deploy the provider locally on `~/.terraform.d/plugins/` directory:
```shell
make deploy-locally
```

Run acceptance tests:
This will start a local Conduktor Console instance and run the acceptance tests against it.
```shell
make testacc
```

Don't forget to run the tests to make sure everything is working as expected before submitting a pull request.

## Styleguides
### Commit Messages
Use explicit commit message that follows the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification. This convention makes it easier to understand the changes in a project and to automate the versioning process.


<!-- omit in toc -->
