# ghr

Github recruitment scraper

## Install

    $ go get github.com/mmcloughlin/ghr

## Usage

Initialize a datastore with the command

    $ ghr init --store=recruitment.db

This will create a SQLite3 database to save searches and prospects.  Then you
can start a repository search with

    $ ghr search --store=recruitment.db --query='stars:>=50 language:Go'

Note that this search will run slowly, because `ghr` is written to adhere to
GitHub API rate limits, and without a token the allowed request rates are very
slow. Please [generate a personal access
token](https://help.github.com/articles/creating-an-access-token-for-command-line-use/)
and provide it to `ghr` with the `--token` option.

    $ ghr search --store=recruitment.db --query='stars:>=50 language:Go' --token='cafe...beef'

GitHub also [requests that you set the `User-Agent`
header](https://developer.github.com/v3/#user-agent-required) appropriately
for API requests. Please use the `--useragent` option to `ghr` to set the user
agent header to your GitHub username, or something else appropriate.
