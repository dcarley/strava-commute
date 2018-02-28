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

### Config

Write a config file to `dist/config.json` in the following format:

    {
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

[bounding box utility]: http://boundingbox.klokantech.com/

### Deployment

1. Export your Strava API token as `STRAVA_API_TOKEN`.
1. Build, package, and deploy:

        make

1. Take the `url` in the output.
