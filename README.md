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
