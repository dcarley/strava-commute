# strava-commute

A utility that runs on [AWS Lambda][] to rename and tag commute rides on [Strava][].

[AWS Lambda]: https://aws.amazon.com/lambda/
[Strava]: https://www.strava.com/

It is inspired by [Alex Muller's][] [lambda-strava-commute-namer][]. The
differences are that I wanted to try out:

- [Lambda's native support for Go][lambda-go]
- [Strava's webhook push subscriptions][strava-webhook]

[Alex Muller's]: http://alex.mullr.net/blog/2017/09/using-lambda-to-do-bits-and-pieces/
[lambda-strava-commute-namer]: https://github.com/alexmuller/lambda-strava-commute-namer
[lambda-go]: https://aws.amazon.com/blogs/compute/announcing-go-support-for-aws-lambda/
[strava-webhook]: http://strava.github.io/api/partner/v3/events/

## Usage

### Prerequisites

1. Install [Terraform](https://www.terraform.io/).
1. Export your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.
1. Create an S3 bucket to store state and initialise Terraform:

        make init BUCKET_SUFFIX=<your_suffix>

1. Create an API application from [your settings on Strava][] with a
   callback domain of `localhost:8080`.

[your settings on Strava]: https://www.strava.com/settings/api

1. Create an API token that has `write` and `view_private` scopes for your
   account using [dcarley/oauth2-cli][]:

        oauth2-cli \
          -scope write,view_private \
          -id <your_client_id> \
          -secret <your_client_secret> \
          -auth https://www.strava.com/oauth/authorize \
          -token https://www.strava.com/oauth/token

[dcarley/oauth2-cli]: https://github.com/dcarley/oauth2-cli

1. Email [developers@strava.com](mailto:developers@strava.com) to ask them
   to enable webhook push subscriptions for your application. Include your
   application client ID and a brief description of what your application
   does.

### Config

Write a config file to `dist/config.json` in the following format:

    {
      "gear_id": "12345",
      "locations": {
        "London": {
          "min": [-0.510375, 51.286758],
          "max": [0.334015, 51.691875]
        },
        "Sheffield": {
          "min": [-1.801472, 53.304512],
          "max": [-1.324669, 53.503128]
        }
      }
    }

Notes about the format:

- The name of each location will be append to activities that have a
    matching start and/or end location. You can use a [bounding box
    utility][] to generate the longitude and latitude co-ordinates.
- The `gear_id` is optional. If present all activities with matching
    locations will be tagged with it. You can get the ID from the URL when
    looking at your bike/shoes on Strava.

[bounding box utility]: http://boundingbox.klokantech.com/

### Deployment

1. Export your Strava API token as `STRAVA_API_TOKEN`.
1. Build, package, and deploy:

        make

1. Create a push subscription if you're deploying for the first time or the
   URL has changed:

        make register \
          STRAVA_CLIENT_ID=<your_client_id> \
          STRAVA_CLIENT_SECRET=<your_client_secret>
