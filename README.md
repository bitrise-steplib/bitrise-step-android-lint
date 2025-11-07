# Android Lint

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/bitrise-step-android-lint?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/bitrise-step-android-lint/releases)

Runs Lint on your Android project source files and detects potential syntax errors to keep your code error free.

<details>
<summary>Description</summary>

It highlights the code line where the error is found, explains the type of error and suggests corrections. The Step does not make the build fail if it spots any structural errors in your code. If you have inserted the **Deploy to Bitrise.io** Step in the Workflow, the tes results will be available in a viewable or downloadable Lint Report HTML or XML file which you can access on the Build's APPS & ARTIFACTS page.

### Configuring the Step

1. Set the **Project Location** input which, by default, points to the root directory of your Android project.
2. Set the module and variant you wish to lint in the **Module** and **Variant** fields.

Optionally, you can modify these inputs:
1. You can specify where the Lint reports should be found once the Step has run if you overwrite the **Report location pattern** input.
2. You can set if the Step should cache build outputs and dependencies, only the dependencies or nothing at all in the **Set level of cache** input.
3. You can set any gradle argument to the gradle task in the **Additional Gradle Arguments** input.

### Troubleshooting
Make sure you insert the Step before a build Step.
Make sure you type the correct module and variant names in the respective fields of the Step. If you are unsure about the exact names, you can check them in the **Project Structure** dialog of your project in Android Studio.

### Useful links
- [Improve your code with lint checks](https://developer.android.com/studio/write/lint)

### Related Steps
- [Android Build](https://www.bitrise.io/integrations/steps/android-build)
- [Android Unit Test](https://www.bitrise.io/integrations/steps/android-unit-test)
- [[BETA] Virtual Device Testing for Android](https://www.bitrise.io/integrations/steps/virtual-device-testing-for-android)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `project_location` | The root directory of your android project, for example, where your root build gradle file exists (also gradlew, settings.gradle, etc...) | required | `$BITRISE_SOURCE_DIR` |
| `module` | Set the module that you want to lint. To see your available modules please open your project in Android Studio and go in [Project Structure] and see the list on the left.  |  |  |
| `variant` | Set the variant that you want to lint. To see your available variants please open your project in Android Studio and go in [Project Structure] -> variants section.  |  |  |
| `report_path_pattern` | Will find the report file with the given pattern. If you need the xml file then you can use: "*/build/reports/lint-results*.xml"  | required | `*/build/reports/lint-results*.html` |
| `cache_level` | `all` - will cache build cache and dependencies `only_deps` - will cache dependencies only `none` - will not cache anything | required | `only_deps` |
| `arguments` | Extra arguments passed to the gradle task |  |  |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-android-lint/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-android-lint/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
