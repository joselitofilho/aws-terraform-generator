<link rel="stylesheet" href="markdown-styles-list.css">

# Contributing

Contributions to this project are [released](#contributions-under-repository-license) to the public under the [project's open source license](LICENSE).

This project adheres to a [Code of Conduct][code-of-conduct]. By participating, you are expected to honor this code.

[code-of-conduct]: CODE_OF_CONDUCT.md

The purpose of this tool is to generate code from diagrams, create diagrams from a code folder, and compare diagrams. For detailed instructions on how it operates, refer to the [README](README.md).

If you find any issues or have suggestions for improvements, feel free to create an [issue][issues] or submit a pull request. Your contribution is much appreciated.

## Contributions Under Repository License

Whenever you add Content to a repository containing notice of a license, you license that Content under the same terms, and you agree that you have the right to license that Content under those terms. If you have a separate agreement to license that Content under different terms, such as a contributor license agreement, that agreement will supersede.

Isn't this just how it works already? Yep. This is widely accepted as the norm in the open-source community; it's commonly referred to by the shorthand "inbound=outbound". We're just making it explicit.

## Submitting a Pull Request

1. Fork it.
2. Create a branch (`git checkout -b adds-support-to-<name>`)
3. Commit your changes (`git commit -am "Adds support to the <name> AWS resource"`)
4. Push to the branch (`git push origin adds-support-to-<name>`)
5. Open a [Pull Request][pull_request]
6. Enjoy a refreshing Drink and wait

## Testing

To run the tests:

```bash
$ task -d scripts tests
```

> [!IMPORTANT]
> There are some tests that verify the equality of a map. If they fail, run them again several times until they pass. This will soon be corrected by the testing library.

If nothing complains, congratulations!

## Releasing a new version

If you are the current maintainer of this gem:

1. Update documentations if necessary: [README](README.md), [configuration](CONFIGURATION.md), and [template](TEMPLATE.md)
2. Retrieve the latest tag created for the most recent version generated `git describe --tags --abbrev=0`
3. Generate a new tag by incrementing the latest one `git tag v0.0.0` (v + major + .minor + .patch)
4. Push the new tag created `git push --tags`
5. Create a new release on GitHub for the new tag created [here][release]

[issues]: https://github.com/joselitofilho/aws-terraform-generator/issues
[pull_request]: https://github.com/joselitofilho/aws-terraform-generator/pulls
[release]: https://github.com/joselitofilho/aws-terraform-generator/releases/new