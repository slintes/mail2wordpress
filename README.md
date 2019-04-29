# mail2wordpress

This is a little app which helps to publish playlists.

## Context

I have to publish a playlist of a radio station on a Wordpress based website every 2 weeks. I get the playlist via
email as an CSV formatted attachment. Currently I'm using a local converter only (CSV to HTML), and create the new 
post on Wordpress manually. The goal is to have this fully automated.

## Idea

- forward email with CSV attachment to [IFTTT email trigger](https://ifttt.com/email)
- IFTTT serves the attachment via HTTP
- IFTTT triggers the mail2wordpress service (this project) with the [webhooks action](https://ifttt.com/services/maker_webhooks)
- mail2wordpress extracts the attachment URL from the HTTP request body and downloads the attachment
- mail2wordpress converts the CSV content to HTML
- mail2wordpress creates a draft post on the Wordpress website with the converted HTML
- last manual step: approve the draft post

## TODO

- add tests...

# License

Copyright 2019 Marc Sluiter

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
