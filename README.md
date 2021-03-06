<!--

Copyright 2017 Sonatype

Licensed under the Apache License, Version 2.0 (the "License"); 
you may not use this file except in compliance with the License. 
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software 
distributed under the License is distributed on an "AS IS" BASIS, 
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. 
See the License for the specific language governing permissions and 
limitations under the License.

-->
# Nexus Webhook Publish

[![Join the chat at https://gitter.im/sonatype/nexus-developers](https://badges.gitter.im/sonatype/nexus-developers.svg)](https://gitter.im/sonatype/nexus-developers?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## What is this?

This a lil Golang app that gets a Webhook payload from Nexus Repository Manager 3.1+ and publish to currently:

- npmjs.org

The primary use case for this app is if you are publishing npm components to Nexus Repository but you also want to publish them to npmjs (and eventually other public repos). You COULD do this as a part of your CI build, but why not use Golang to do it!

## Developing

To get up and running you'll need Golang and `go dep` installed. Once you've done that:

- `dep ensure` from the root

From there you can play around with it, build it, etc...

## Using

Right now this is pretty WIP, but you can test it out by:

- Configuring a [Repository Webhook Capability](https://help.sonatype.com/display/NXRM3/Webhooks) in Nexus Repository Manager, specifically the component event
- Set the secret key in the capability to something good! I have it set to duckduckgoose so you'll be so embarrassed you need to change it (this is a const in `main.go` for those looking)
- Ensure the `NEXUS_REPO_BASE_URL` is set to whatever your Nexus Repo instance is, I have it defaulted to how Nexus Repo runs by default (also in `main.go`)
- Set the address to send the webhook to `http://localhost:8000/publish` or wherever this go app will end up running
- Modify `webhook/webhook.go` and set your secret key in there
- `go run main.go`
- Install npm (likely version 4, as there is a bug in 5 currently that prevents publishing tarballs)
- Login on npm using `npm login` as the user that you will be publishing to npm as

Now Nexus Repository Manager should send a webhook over to this service. If the event is for an npm component (what we refer to new packages as), it will shell out to npm and publish it. Assumedly if you wanted to limit this just to a specific repo, some code can be added to do that by checking c.RepositoryName in one of the functions where component is being used.

## The Fine Print

It is worth noting that this is **NOT SUPPORTED** by Sonatype, and is a contribution of ours
to the open source community (read: you!)

Remember:

* Use this contribution at the risk tolerance that you have
* Do NOT file Sonatype support tickets related to this Golang project
* DO file issues here on GitHub, so that the community can pitch in

Phew, that was easier than I thought. Last but not least of all:

Have fun creating and using this plugin and the Nexus platform, we are glad to have you here!

## Getting help

Looking to contribute to our code but need some help? There's a few ways to get information:

* Chat with us on [Gitter](https://gitter.im/sonatype/nexus-developers)
* Check out the [Nexus3](http://stackoverflow.com/questions/tagged/nexus3) tag on Stack Overflow
* Check out the [Nexus Repository User List](https://groups.google.com/a/glists.sonatype.com/forum/?hl=en#!forum/nexus-users)
* Connect with [@sonatypeDev](https://twitter.com/sonatypedev) on Twitter

## Acknowledgements 

We stand on the shoulders of giants. Some of the code in this project was graciously borrowed (and attributed to):

- https://github.com/rjz/githubhook
